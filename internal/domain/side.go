package domain

import "strconv"

// InvalidSideError indicates that the provided string does not
// match any known Side value.
type InvalidSideError struct {
	Have string
}

func (e InvalidSideError) Error() string {
	return "invalid side: " + e.Have
}

// Side describes the direction of an order: buy or sell.
type Side byte

const (
	sideInvalid Side = 0

	// SideBid represents a buy order.
	SideBid Side = 1

	// SideAsk represents a sell order.
	SideAsk Side = 2
)

// SideFromString creates a Side from its string representation
// ("bid" or "ask"). Returns InvalidSideError if the string does not
// match any known value.
func SideFromString(str string) Side {
	switch str {
	case "bid":
		return SideBid
	case "ask":
		return SideAsk
	default:
		return sideInvalid
	}
}

// String returns the string representation of Side.
// For invalid values it returns "invalid(<number>)".
func (s Side) String() string {
	switch s {
	case SideBid:
		return "bid"
	case SideAsk:
		return "ask"
	default:
		return "invalid(" + strconv.Itoa(int(s)) + ")"
	}
}

// IsValid reports whether the value is one of the known
// valid Side values (SideBid, SideAsk).
func (s Side) IsValid() bool {
	switch s {
	case SideBid, SideAsk:
		return true
	default:
		return false
	}
}
