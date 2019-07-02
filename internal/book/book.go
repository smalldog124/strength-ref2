package book

import "time"

func (s *Seating) State(now time.Time) SeatingState {
	if s.Booked {
		return Booked
	}

	if s.ExpireTimestamp > 0 && (GetTimestamp(now)-s.ExpireTimestamp) < TimeLimitMS {
		return Reserved
	} else {
		return Free
	}
}

//// Funcs ////
func GetTimestamp(d time.Time) int64 {
	return d.UnixNano() / int64(time.Millisecond)
}
