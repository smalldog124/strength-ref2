package book

type Seating struct {
	ExpireTimestamp int64 `json:"expireTimestamp"`
	Booked          bool  `json:"booked"`
}

type SeatingState int

const (
	Free     SeatingState = 0
	Reserved SeatingState = 1
	Booked   SeatingState = 2
)

const TimeLimitMS = 10 * 1000
