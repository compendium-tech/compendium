package domain

const (
	RoleSystem    LLMRole = "system"
	RoleUser      LLMRole = "user"
	RoleAssistant LLMRole = "assistant"

	TypeString  LLMType = "STRING"
	TypeNumber  LLMType = "NUMBER"
	TypeInteger LLMType = "INTEGER"
	TypeBoolean LLMType = "BOOLEAN"
	TypeArray   LLMType = "ARRAY"
	TypeObject  LLMType = "OBJECT"
)

type LLMRole string
type LLMType string

type LLMMessage struct {
	Role      LLMRole
	Text      string
	ToolCalls []LLMToolCall
}

type LLMToolCall struct {
	ID         string
	Name       string
	Parameters map[string]any
}

type LLMToolDefinition struct {
	Name             string
	Description      string
	ParametersSchema *LLMSchema
}

type LLMSchema struct {
	Type        LLMType
	Description string
	Properties  map[string]LLMSchema
	Items       *LLMSchema
	MaxItems    *int64
	MinItems    *int64
	Required    []string
}
