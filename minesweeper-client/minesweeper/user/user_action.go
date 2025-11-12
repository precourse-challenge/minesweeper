package user

type Action string

const (
	Open    Action = "open"
	Flag    Action = "flag"
	Quit    Action = "quit"
	Unknown Action = "unknown"
)

func From(input string) Action {
	switch input {
	case "open":
		return Open
	case "flag":
		return Flag
	case "quit":
		return Quit
	default:
		return Unknown
	}
}
