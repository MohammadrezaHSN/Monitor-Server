package handler

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func waitForKey() {
	fmt.Println("\nPress any key to exit...")
	_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Route() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Program crashed with error:", r)
			waitForKey()
		}
	}()

	router := gin.Default()
	fmt.Println("Agent")
	router.GET("/api/agents", agentGetHandler)
	router.POST("/api/agents/register", agentRegisterHandler)
	router.POST("/api/agents/update", agentUpdateHandler)
	router.POST("/api/agents/delete", agentDeleteHandler)
	router.POST("/api/agents/get-agent-with-id", agentGetWithIdHandler)

	fmt.Println("Request")
	router.GET("/api/requests", requestGetHandler)
	router.POST("/api/request/create", requestCreateHandler)
	router.POST("/api/request/update", requestUpdateHandler)
	router.POST("/api/request/delete", requestDeleteHandler)
	router.POST("/api/request/get-request-with-id", requestGetWithIdHandler)

	fmt.Println("Response")
	router.GET("/api/response", responseGetHandler)
	router.POST("/api/response/create", responseCreateHandler)
	router.POST("/api/response/update", responseUpdateHandler)
	router.POST("/api/response/delete", responseDeleteHandler)
	router.POST("/api/response/get-response-with-id", responseGetWithIdHandler)

	fmt.Println("File")
	router.GET("/api/files", fileGetHandler)
	router.POST("/api/files/store", fileStoreHandler)
	router.POST("/api/files/update", fileUpdateHandler)
	router.POST("/api/files/delete", fileDeleteHandler)
	router.POST("/api/files/get-file-with-id", fileGetWithIdHandler)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
		waitForKey()
		return
	}

}
