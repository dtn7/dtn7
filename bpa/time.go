package bpa

import "time"

// DTNTime is an integer indicating the time like the Unix time, just starting
// from the year 2000 instead of 1970. This is specified in section 4.1.6.
type DTNTime uint64

const seconds1970To2k = 946684800

// Unix returns the Unix timestamp for this DTNTime.
func (t *DTNTime) Unix() int64 {
	return (int64)(*t) + seconds1970To2k
}

// Time returns a UTC-based time.Time for this DTNTime.
func (t *DTNTime) Time() time.Time {
	return time.Unix(t.Unix(), 0).UTC()
}

// DTNTimeFromTime returns the DTNTime for the time.Time.
func DTNTimeFromTime(t time.Time) DTNTime {
	return (DTNTime)(t.Unix() - seconds1970To2k)
}

// DTNTimeNow returns the current (UTC) time as DTNTime.
func DTNTimeNow() DTNTime {
	return DTNTimeFromTime(time.Now())
}

// CreationTimestamp is a tuple of a DTNTime and a sequence number to differ
// bundles with the same DTNTime (seconds) from the same endpoint. It is
// specified in section 4.1.7.
type CreationTimestamp struct {
	time           DTNTime
	sequenceNumber uint
}
