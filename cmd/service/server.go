package service

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"html/template"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ipLogs       []string
	commandsList []string
	mu           sync.Mutex
	templates    *template.Template
)

const uploadBaseDir = `C:\Users\Gopher\Desktop\job\go\goServer\server\files`

type PageData struct {
	Title string
	Data  any
}

func loadTemplates() {
	templates = template.Must(template.ParseGlob("web/templates/*.html"))

}

func renderTemplate(w http.ResponseWriter, name string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
	}
}

// ========== HTTP handler ==========
func messageHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	mu.Lock()
	ipLogs = append(ipLogs, ip)
	mu.Unlock()

	renderTemplate(w, "index.html", PageData{
		Title: "Go Server Dashboard",
		Data: map[string]any{
			"IP": ip,
		},
	})
}

func showIps(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("format") == "json" {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"ips": ipLogs,
		})
		return
	}

	renderTemplate(w, "ips.html", PageData{
		Title: "Collected IPs",
	})
}

func commands(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Query().Get("format") == "json" {
			mu.Lock()
			defer mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"commands": commandsList,
			})
			return
		}

		renderTemplate(w, "commands.html", PageData{
			Title: "Command Center",
		})
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		cmd := strings.TrimSpace(string(body))
		if cmd == "" {
			http.Error(w, "Empty command not allowed", http.StatusBadRequest)
			return
		}

		mu.Lock()
		commandsList = append(commandsList, cmd)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Command accepted!",
		})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func sanitizeForWindowsFolder(ip string) string {
	return strings.ReplaceAll(ip, ":", "-")
}

// ========== HTTP/HTTPS upload ==========
func handleHTTPUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "upload.html", PageData{
			Title: "Upload Center",
		})
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	safeIP := sanitizeForWindowsFolder(ip)

	fileName := r.Header.Get("X-File-Name")
	if fileName == "" {
		fileName = "uploaded_file.bin"
	}
	fileName = filepath.Base(fileName)

	targetDir := filepath.Join(uploadBaseDir, "http", safeIP)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		http.Error(w, "Server error creating directory", http.StatusInternalServerError)
		return
	}

	targetPath := filepath.Join(targetDir, fileName)
	dst, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, "Server error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, r.Body); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("File %s uploaded successfully via HTTP", fileName),
	})
}

// ========== TCP server ==========
func startTCPFileServer(ip, port string) {
	listener, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		fmt.Println("TCP listener error:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()

			var nameLen uint16
			if err := binary.Read(c, binary.BigEndian, &nameLen); err != nil {
				return
			}

			nameBytes := make([]byte, nameLen)
			if _, err := io.ReadFull(c, nameBytes); err != nil {
				return
			}

			fileName := filepath.Base(string(nameBytes))

			ip, _, err := net.SplitHostPort(c.RemoteAddr().String())
			if err != nil {
				ip = c.RemoteAddr().String()
			}

			safeIP := sanitizeForWindowsFolder(ip)
			targetDir := filepath.Join(uploadBaseDir, "tcp", safeIP)
			os.MkdirAll(targetDir, 0755)

			dst, err := os.Create(filepath.Join(targetDir, fileName))
			if err != nil {
				return
			}
			defer dst.Close()

			io.Copy(dst, c)
		}(conn)
	}
}

// ========== UDP server ==========
type udpClientState struct {
	file        *os.File
	lastWrite   time.Time
	expectedSeq uint32
	buffer      map[uint32][]byte
	complete    bool
}

const maxUDPWindow = 128

func startUDPFileServer(ip, port string) {
	addr, err := net.ResolveUDPAddr("udp4", ip+":"+port)
	if err != nil {
		fmt.Println("UDP resolve error:", err)
		return
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("UDP listen error:", err)
		return
	}
	defer conn.Close()

	clients := make(map[string]*udpClientState)
	var udpMu sync.Mutex

	go func() {
		for {
			time.Sleep(10 * time.Second)
			udpMu.Lock()
			for clientKey, state := range clients {
				if time.Since(state.lastWrite) > 30*time.Second {
					if state.file != nil {
						state.file.Close()
					}
					delete(clients, clientKey)
					fmt.Printf("UDP: Cleaned up inactive session %s\n", clientKey)
				}
			}
			udpMu.Unlock()
		}
	}()

	buf := make([]byte, 2000)
	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil || n < 4 {
			continue
		}

		clientKey := sanitizeForWindowsFolder(remote.String())
		seq := binary.BigEndian.Uint32(buf[:4])

		udpMu.Lock()

		if seq == 0 {
			if n < 6 {
				udpMu.Unlock()
				continue
			}
			nameLen := int(binary.BigEndian.Uint16(buf[4:6]))
			if 6+nameLen > n {
				udpMu.Unlock()
				continue
			}

			state, exists := clients[clientKey]
			if !exists {
				fileName := filepath.Base(string(buf[6 : 6+nameLen]))
				targetDir := filepath.Join(uploadBaseDir, "udp", clientKey)
				if err := os.MkdirAll(targetDir, 0755); err != nil {
					fmt.Println("UDP: Failed to create directory:", err)
					udpMu.Unlock()
					continue
				}

				f, err := os.Create(filepath.Join(targetDir, fileName))
				if err != nil {
					fmt.Println("UDP: Failed to create file:", err)
					udpMu.Unlock()
					continue
				}

				state = &udpClientState{
					file:        f,
					lastWrite:   time.Now(),
					expectedSeq: 1,
					buffer:      make(map[uint32][]byte),
					complete:    false,
				}
				clients[clientKey] = state
			}

			ack := make([]byte, 4)
			binary.BigEndian.PutUint32(ack, 0)
			conn.WriteToUDP(ack, remote)
			udpMu.Unlock()
			continue
		}

		state, exists := clients[clientKey]
		if !exists || state.complete {
			udpMu.Unlock()
			continue
		}

		if n < 5 {
			udpMu.Unlock()
			continue
		}

		flag := buf[4]
		payload := buf[5:n]

		if seq < state.expectedSeq {
			ack := make([]byte, 4)
			binary.BigEndian.PutUint32(ack, state.expectedSeq)
			conn.WriteToUDP(ack, remote)
			udpMu.Unlock()
			continue
		}

		if seq >= state.expectedSeq+maxUDPWindow {
			udpMu.Unlock()
			continue
		}

		entry := make([]byte, 1+len(payload))
		entry[0] = flag
		copy(entry[1:], payload)
		state.buffer[seq] = entry
		state.lastWrite = time.Now()

		for {
			data, ok := state.buffer[state.expectedSeq]
			if !ok {
				break
			}
			delete(state.buffer, state.expectedSeq)
			f := data[0]
			p := data[1:]

			if len(p) > 0 {
				state.file.Write(p)
			}
			if f == 0x01 {
				state.complete = true
				state.file.Close()
				state.file = nil
				go func(key string) {
					time.Sleep(2 * time.Second)
					udpMu.Lock()
					delete(clients, key)
					udpMu.Unlock()
				}(clientKey)
			}
			state.expectedSeq++
		}

		ack := make([]byte, 4)
		binary.BigEndian.PutUint32(ack, state.expectedSeq)
		conn.WriteToUDP(ack, remote)
		udpMu.Unlock()
	}
}

// ========== HTTPS server ==========
func generateSelfSignedCert() (tls.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}
	httpsTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"GoServer Test"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &httpsTemplate, &httpsTemplate, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return tls.X509KeyPair(certPEM, keyPEM)
}

func startHTTPSFileServer(ip, port string) {
	cert, err := generateSelfSignedCert()
	if err != nil {
		fmt.Println("TLS cert error:", err)
		return
	}
	server := &http.Server{
		Addr:      ip + ":" + port,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
	}
	fmt.Printf("HTTPS server on %s:%s\n", ip, port)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Println("HTTPS error:", err)
	}
}

func promptInput(promptMsg, defaultValue string) string {
	fmt.Println(promptMsg)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return defaultValue
	}
	return input
}

func fetchSharedFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Query().Get("ip") == "" {
		renderTemplate(w, "fetch_share.html", PageData{
			Title: "Fetch Shared File",
		})
		return
	}

	targetIP := r.URL.Query().Get("ip")
	targetPort := r.URL.Query().Get("port")
	filename := r.URL.Query().Get("file")

	if targetIP == "" || targetPort == "" || filename == "" {
		http.Error(w, "Missing parameters: ip, port, file", http.StatusBadRequest)
		return
	}

	clientURL := fmt.Sprintf("http://%s:%s/%s", targetIP, targetPort, filename)
	resp, err := http.Get(clientURL)
	if err != nil {
		http.Error(w, "Failed to connect to client share: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, resp.Body)
}

func RunServer() {
	err := os.MkdirAll(uploadBaseDir, 0755)
	if err != nil {
		fmt.Println(err)
	}
	loadTemplates()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", messageHandler)
	http.HandleFunc("/ips", showIps)
	http.HandleFunc("/commands", commands)
	http.HandleFunc("/upload", handleHTTPUpload)
	http.HandleFunc("/fetch-share", fetchSharedFileHandler)

	ip := promptInput("Insert IP ADDRESS (Default 0.0.0.0)", "0.0.0.0")
	httpPort := promptInput("Insert HTTP port (Default 8080): ", "8080")
	tcpPort := promptInput("Insert TCP port (Default 9090): ", "9090")
	udpPort := promptInput("Insert UDP port (Default 9091): ", "9091")
	httpsPort := promptInput("Insert HTTPS port (Default 8443): ", "8443")

	go startTCPFileServer(ip, tcpPort)

	go startUDPFileServer(ip, udpPort)
	go startHTTPSFileServer(ip, httpsPort)

	fmt.Printf("HTTP server on %s:%s\n", ip, httpPort)
	if err := http.ListenAndServe(ip+":"+httpPort, nil); err != nil {
		panic(err)
	}
}
