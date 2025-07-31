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

type StructuredOutputSchema struct {
	Type        Type                              `json:"type,omitempty"`
	Description string                            `json:"description,omitempty"`
	Properties  map[string]StructuredOutputSchema `json:"properties,omitempty"`
	Items       *StructuredOutputSchema           `json:"items,omitempty"`
	MaxItems    *int64                            `json:"maxItems,omitempty"`
	MinItems    *int64                            `json:"minItems,omitempty"`
	Required    []string                          `json:"required,omitempty"`
}
