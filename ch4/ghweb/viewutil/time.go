package viewutil

import (
	"strings"
	"time"
)

type ViewTime time.Time

func NewTime() ViewTime {
	return ViewTime(time.Now())
}

func (vt ViewTime) String() string {
	return time.Time(vt).Format(time.RFC1123)
}

func (vt *ViewTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02T15:04:05Z", s)
	*vt = ViewTime(t)
	return
}

func (vt *ViewTime) MarshalJSON() ([]byte, error) {
	return []byte(vt.String()), nil
}

func (vt *ViewTime) Time() time.Time {
	return time.Time(*vt)
}
