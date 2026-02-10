package schema

const (
	tsNever          = "never"
	tsUnknown        = "unknown"
	tsRecordUnknown  = "Record<string, unknown>"
	tsEmptyObject    = "{}"
	schemaTypeString = "string"
	schemaTypeNull   = "null"
	enumValuePrefix  = "Value"
	enumNumberPrefix = "VALUE_"
)

const (
	modeDefault schemaMode = iota
	modeInput
	modeOutput
)

const (
	enumInvalid enumKind = iota
	enumString
	enumNumber
)
