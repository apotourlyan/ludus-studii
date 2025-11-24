package types

type Command string

const (
	CommandInit   Command = "init"
	CommandCreate Command = "create"
	CommandUp     Command = "up"
	CommandDown   Command = "down"
)
