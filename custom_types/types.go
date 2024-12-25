package custom_types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	Time time.Time
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

	t.Time = parsedTime
	return nil
}

func (t Time) Value() (driver.Value, error) {
	return t.Time.Format(time.TimeOnly), nil
}

type Date struct {
	Time time.Time
}

func (d Date) Scan(src any) error {
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

	d.Time = parsedDate
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return d.Time.Format(time.DateOnly), nil
}
