package enum

type PaymentStatus int

const (
	Unpaid PaymentStatus = iota
	Paid
	Failed
)

func (s PaymentStatus) String() string {
	switch s {
	case Unpaid:
		return "Unpaid"
	case Paid:
		return "Paid"
	case Failed:
		return "Failed"
	default:
		return "Unpaid"
	}
}
