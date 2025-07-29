package domain

const (
	ROLE_SYSTEM    Role = "system"
	ROLE_USER      Role = "user"
	ROLE_ASSISTANT Role = "assistant"
	ROLE_TOOL      Role = "tool"
)

type Role string

type Message struct {
	Role        Role
	TextContent string
	ToolCalls   []ToolCall
}

type ToolCall struct {
	ID         string
	Name       string
	Parameters map[string]any
}

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []ToolParameter
}

type ToolParameter struct {
	Type        string
	Name        string
	Description string
	IsRequired  bool
	Enum        []string
}

type Response struct {
	Message        Message
	GoogleSearches []string
}
