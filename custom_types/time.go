package custom_types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	_time time.Time
}

func NewTime(t time.Time) (*Time, error) {
	onlyTime, err := time.Parse(time.TimeOnly, t.Format(time.TimeOnly))
	if err != nil {
		return nil, err
	}
	return &Time{_time: onlyTime}, nil
}

func (t Time) UnderlyingTime() time.Time {
	return t._time
}

func (t Time) Equal(arg Time) bool {
	return t._time.Format(time.TimeOnly) == arg._time.Format(time.TimeOnly)
}

func (t Time) Sub(arg Time) time.Duration {
	return t._time.Sub(arg._time)
}

func (t *Time) Scan(src any) error {
	if src == nil {
		return nil
	}

	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot scan type %T into custom_types.Time", src)
	}

	parsedTime, err := time.Parse(time.TimeOnly, str)
	if err != nil {
		return err
	}

	t._time = parsedTime
	return nil
}

func (t *Time) Value() (driver.Value, error) {
	return t._time.Format(time.TimeOnly), nil
}

func ParseTime(str string) (*Time, error) {
	t, err := time.Parse(time.TimeOnly, str)
	if err != nil {
		return nil, err
	}
	return &Time{_time: t}, nil
}

func TimeDiffAbs(a, b Time) time.Duration {
	return a.Sub(b).Abs()
}
