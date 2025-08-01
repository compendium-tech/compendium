package domain

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"

	TypeString  Type = "STRING"
	TypeNumber  Type = "NUMBER"
	TypeInteger Type = "INTEGER"
	TypeBoolean Type = "BOOLEAN"
	TypeArray   Type = "ARRAY"
	TypeObject  Type = "OBJECT"
)

type Role string
type Type string

type Message struct {
	Role      Role
	Text      string
	ToolCalls []ToolCall
}

type ToolCall struct {
	ID         string
	Name       string
	Parameters map[string]any
}

type ToolDefinition struct {
	Name             string
	Description      string
	ParametersSchema *Schema
}

type Schema struct {
	Type        Type
	Description string
	Properties  map[string]Schema
	Items       *Schema
	MaxItems    *int64
	MinItems    *int64
	Required    []string
}
