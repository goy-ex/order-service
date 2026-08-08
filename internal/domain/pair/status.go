package pair

import "strconv"

type PairStatus byte

const (
	pairStatusInvalid  PairStatus = 0
	PairStatusInactive PairStatus = 1
	PairStatusActive   PairStatus = 2
)

func PairStatusFromString(s string) PairStatus {
	switch s {
	case "inactive":
		return PairStatusInactive
	case "active":
		return PairStatusActive
	default:
		return pairStatusInvalid
	}
}

func (s PairStatus) IsValid() bool {
	switch s {
	case PairStatusInactive, PairStatusActive:
		return true
	default:
		return false
	}
}

func (s PairStatus) String() string {
	switch s {
	case PairStatusInactive:
		return "inactive"
	case PairStatusActive:
		return "active"
	case pairStatusInvalid:
		return "invalid"
	default:
		return "invalid(" + strconv.Itoa(int(s)) + ")"
	}
}
