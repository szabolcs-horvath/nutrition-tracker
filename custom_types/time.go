package custom_types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	_time time.Time
}

func (t Time) UnderlyingTime() time.Time {
	return t._time
}

func ParseTime(str string) (*Time, error) {
	t, err := time.Parse(time.TimeOnly, str)
	if err != nil {
		return nil, err
	}
	return &Time{_time: t}, nil
}

func (t Time) Equal(arg Time) bool {
	return t._time.Format(time.TimeOnly) == arg._time.Format(time.TimeOnly)
}

func (t Time) Scan(src any) error {
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

func (t Time) Value() (driver.Value, error) {
	return t._time.Format(time.TimeOnly), nil
}
