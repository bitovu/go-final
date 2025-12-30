package api

import (
	"time"

	"final/pkg/db"
)

const format = "20060102"

func checkDate(task *db.Task, now time.Time) error {
	if task.Date == "" {
		task.Date = now.Format(format)
		return nil
	}

	t, err := time.Parse(format, task.Date)
	if err != nil {
		return err
	}

	today := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)

	if t.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format(format)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}
