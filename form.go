package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// helper to parse common form fields from a request
func parseTodoForm(r *http.Request) (short string, description string, dueDate *time.Time, costOfDelay int16, effort string, err error) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err = r.ParseMultipartForm(1 << 20)
	} else {
		err = r.ParseForm()
	}
	if err != nil {
		return
	}

	short = r.FormValue("short")
	description = r.FormValue("description")

	dateStr := r.FormValue("due_date")
	if dateStr != "" {
		var dt time.Time
		dt, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			err = fmt.Errorf("due_date must be in YYYY-MM-DD format")
			return
		}
		dueDate = &dt
	}

	if codStr := r.FormValue("cost_of_delay"); codStr != "" {
		var tmp int
		tmp, err = strconv.Atoi(codStr)
		if err != nil {
			return
		}
		if tmp < -2 || tmp > 2 {
			err = fmt.Errorf("cost_of_delay must be between -2 and 2")
			return
		}
		costOfDelay = int16(tmp)
	}

	effort = r.FormValue("effort")
	switch effort {
	case "mins", "hours", "days", "weeks", "months":
	case "":
		effort = "hours"
	default:
		err = fmt.Errorf("invalid effort")
		return
	}

	return
}
