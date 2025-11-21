package enum

type Status int

const (
	Pending Status = iota
	Canceled
	Completed
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Canceled:
		return "Canceled"
	case Completed:
		return "Completed"
	default:
		return "Pending"
	}
}
