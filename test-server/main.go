package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // テスト用なので全て許可
	},
}

type Message struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"clientInfo"`
}

type InitializeResponse struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"serverInfo"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		// メッセージを受信
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Unmarshal error: %v", err)
			continue
		}

		log.Printf("Received message: %s", msg.Method)

		// メッセージタイプに応じて処理
		switch msg.Method {
		case "initialize":
			handleInitialize(conn, &msg)
		case "tools/list":
			handleToolsList(conn, &msg)
		case "tools/call":
			handleToolsCall(conn, &msg)
		case "ping":
			handlePing(conn, &msg)
		default:
			log.Printf("Unknown method: %s", msg.Method)
		}
	}
}

func handleInitialize(conn *websocket.Conn, msg *Message) {
	var req InitializeRequest
	if err := json.Unmarshal(msg.Params, &req); err != nil {
		log.Printf("Failed to unmarshal initialize request: %v", err)
		return
	}

	log.Printf("Initializing with client: %s v%s", req.ClientInfo.Name, req.ClientInfo.Version)

	response := InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		ServerInfo: struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    "test-mcp-server",
			Version: "1.0.0",
		},
	}

	result, _ := json.Marshal(response)
	responseMsg := Message{
		ID:     msg.ID,
		Method: "initialize",
		Result: result,
	}

	if err := conn.WriteJSON(responseMsg); err != nil {
		log.Printf("Failed to send initialize response: %v", err)
	}
}

func handleToolsList(conn *websocket.Conn, msg *Message) {
	tools := []map[string]interface{}{
		{
			"name":        "echo",
			"description": "Echo back the input",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{"message"},
			},
		},
	}

	result, _ := json.Marshal(map[string]interface{}{
		"tools": tools,
	})

	responseMsg := Message{
		ID:     msg.ID,
		Method: "tools/list",
		Result: result,
	}

	if err := conn.WriteJSON(responseMsg); err != nil {
		log.Printf("Failed to send tools/list response: %v", err)
	}
}

func handleToolsCall(conn *websocket.Conn, msg *Message) {
	var toolCall struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(msg.Params, &toolCall); err != nil {
		log.Printf("Failed to unmarshal tool call: %v", err)
		return
	}

	log.Printf("Tool call: %s with args: %v", toolCall.Name, toolCall.Arguments)

	// 簡単なエコーツールの実装
	var result interface{}
	if toolCall.Name == "echo" {
		if message, ok := toolCall.Arguments["message"].(string); ok {
			result = map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf("Echo: %s", message),
					},
				},
				"isError": false,
			}
		}
	}

	resultData, _ := json.Marshal(result)
	responseMsg := Message{
		ID:     msg.ID,
		Method: "tools/call",
		Result: resultData,
	}

	if err := conn.WriteJSON(responseMsg); err != nil {
		log.Printf("Failed to send tools/call response: %v", err)
	}
}

func handlePing(conn *websocket.Conn, msg *Message) {
	responseMsg := Message{
		ID:     msg.ID,
		Method: "pong",
	}

	if err := conn.WriteJSON(responseMsg); err != nil {
		log.Printf("Failed to send pong: %v", err)
	}
}

func main() {
	http.HandleFunc("/", handleWebSocket)

	port := ":3000"
	log.Printf("Starting MCP test server on port %s", port)
	log.Printf("Connect your client to: ws://localhost%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
