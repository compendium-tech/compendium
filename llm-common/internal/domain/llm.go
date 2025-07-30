package domain

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
	TypeString    Type = "STRING"
	TypeNumber    Type = "NUMBER"
	TypeInteger   Type = "INTEGER"
	TypeBoolean   Type = "BOOLEAN"
	TypeArray     Type = "ARRAY"
	TypeObject    Type = "OBJECT"
	TypeNULL      Type = "NULL"
)

type Role string
type Type string

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
	Type        Type
	Name        string
	Description string
	IsRequired  bool
	Enum        []string
}

type Response struct {
	Message Message
}
