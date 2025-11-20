package user

type Action string

const (
	Open          Action = "open"
	Flag          Action = "flag"
	Exit          Action = "exit"
	Retry         Action = "retry"
	UnknownAction Action = "unknown"
)

func ActionFrom(input string) Action {
	switch input {
	case "open":
		return Open
	case "flag":
		return Flag
	case "exit":
		return Exit
	case "retry":
		return Retry
	default:
		return UnknownAction
	}
}
