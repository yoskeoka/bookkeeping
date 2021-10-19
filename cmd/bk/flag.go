package main

import (
	"fmt"
	"time"
)

type dateFlag struct {
	date *time.Time
}

func (d *dateFlag) Set(v string) error {
	t, err := time.Parse("20060102", v)
	if err == nil {
		d.date = &t
		return nil
	}

	return fmt.Errorf("-date format doesn't match any of supported format, where the supported format is 'yyyymmdd'")
}

func (d *dateFlag) String() string {
	return d.date.Format("20060102")
}
