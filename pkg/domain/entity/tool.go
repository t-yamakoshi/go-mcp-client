package entity

// Tool represents a tool definition
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolCall represents a tool call request
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents a tool call result
type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError"`
}

// Content represents content in a tool result
type Content struct {
	Type     string      `json:"type"`
	Text     string      `json:"text,omitempty"`
	ImageURL string      `json:"imageUrl,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
