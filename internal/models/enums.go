package models

type ReserveState int

const (
	Reserved ReserveState = iota
	Cancelled
	Accepted
)

func (rs ReserveState) String() string {
	switch rs {
	case Reserved:
		return "reserved"
	case Cancelled:
		return "cancelled"
	case Accepted:
		return "accepted"
	}
	return "unknown"
}

type TransferType int

const (
	Deposit TransferType = iota
	Withdraw
)

func (tt TransferType) String() string {
	switch tt {
	case Deposit:
		return "deposit"
	case Withdraw:
		return "withdraw"
	}
	return "unknown"
}
