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

func (rs *ReserveState) FromString(val string) {
	switch val {
	case "reserved":
		*rs = Reserved
	case "cancelled":
		*rs = Cancelled
	case "accepted":
		*rs = Accepted
	}
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

func (tt *TransferType) FromString(val string) {
	switch val {
	case "deposit":
		*tt = Deposit
	case "withdraw":
		*tt = Withdraw
	}
}
