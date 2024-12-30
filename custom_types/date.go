package custom_types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Date struct {
	_time time.Time
}

func NewDate(t time.Time) (*Date, error) {
	onlyDate, err := time.Parse(time.DateOnly, t.Format(time.DateOnly))
	if err != nil {
		return nil, err
	}
	return &Date{_time: onlyDate}, nil
}

func (d *Date) UnderlyingTime() time.Time {
	return d._time
}

func (d *Date) Sub(arg Date) time.Duration {
	return d._time.Sub(arg._time)
}

func (d *Date) Equal(arg Date) bool {
	return d._time.Format(time.DateOnly) == arg._time.Format(time.DateOnly)
}

func (d *Date) String() string {
	return d._time.Format(time.DateOnly)
}

func (d *Date) Scan(src any) error {
	if src == nil {
		return nil
	}

	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot scan type %T into custom_types.Date", src)
	}

	parsedDate, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}

	d._time = parsedDate
	return nil
}

func (d *Date) Value() (driver.Value, error) {
	return d._time.Format(time.DateOnly), nil
}

func ParseDate(str string) (*Date, error) {
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return nil, err
	}
	return &Date{_time: t}, nil
}

func DateDiffAbs(a, b Date) time.Duration {
	return a.Sub(b).Abs()
}
