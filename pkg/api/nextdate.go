package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(format, dstart)
	if err != nil {
		return "", errors.New("date format error")
	}

	if repeat == "" {
		return "", errors.New("Repeat time not valid")
	}
	if strings.HasPrefix(repeat, "d ") {
		N := strings.Split(repeat, " ")
		if len(N) != 2 {
			return "", errors.New("Repeat format error")
		}
		n, err := strconv.Atoi(N[1])
		if err != nil {
			return "", err
		}
		if n < 1 || n > 400 {
			return "", errors.New("Repeat format error")
		}
		for {
			date = date.AddDate(0, 0, n)
			if afterNow(date, now) {
				break
			}
		}
	} else if repeat == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	} else {
		return "", errors.New("invalid repeat")
	}
	return date.Format("20060102"), nil
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	if y1 != y2 {
		return y1 > y2
	}
	if m1 != m2 {
		return m1 > m2
	}
	return d1 > d2
}
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(format, nowStr)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": "now format error",
			})
			return
		}
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "next date error",
		})
		return
	}

	w.Write([]byte(result))
}
