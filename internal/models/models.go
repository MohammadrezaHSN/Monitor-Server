package models

import (
	"time"
)

type Agent struct {
	ID            int64      `json:"id"`
	SystemID      string     `json:"system_id"`
	Interval      uint16     `json:"interval"`
	IsAlive       bool       `json:"is_alive"`
	Description   string     `json:"description"`
	LastActive    *time.Time `json:"last_active,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	NextRequestAt *time.Time `json:"next_request_at,omitempty"`
}

type Request struct {
	ID                 int64      `json:"id"`
	AgentID            int64      `json:"agent_id"`
	Context            string     `json:"context"`
	IsDone             bool       `json:"is_done"`
	CreatedAt          time.Time  `json:"created_at"`
	ReceivedToAgentAt  *time.Time `json:"received_to_agent_at,omitempty"`
	EstimateResponseAt *time.Time `json:"estimate_response_at,omitempty"`
}

type Response struct {
	ID                 int64      `json:"id"`
	AgentID            int64      `json:"agent_id"`
	RequestID          int64      `json:"request_id"`
	Result             string     `json:"result"`
	IsSuccessful       bool       `json:"is_successful"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	ReceivedToServerAt time.Time  `json:"received_to_server_at"`
}

type File struct {
	ID                       int64      `json:"id"`
	AgentID                  int64      `json:"agent_id"`
	RequestID                int64      `json:"request_id"`
	FilePathFromAgentSystem  string     `json:"file_path_from_agent_system"`
	FilePathFromServerSystem string     `json:"file_path_from_server_system"`
	IsCompleted              bool       `json:"is_completed"`
	UploadAt                 *time.Time `json:"upload_at,omitempty"`
	ReceivedToServerAt       time.Time  `json:"received_to_server_at"`
}
