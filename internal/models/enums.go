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
	Crediting TransferType = iota
	Debiting
)

func (tt TransferType) String() string {
	switch tt {
	case Crediting:
		return "crediting"
	case Debiting:
		return "debiting"
	}
	return "unknown"
}
