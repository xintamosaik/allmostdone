package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TodoInput struct {
	Short       string
	Description string
	DueDate     *time.Time
	CostOfDelay int16
	Effort      string
}

type TodoFormData struct {
	Input          TodoInput
	Error          string
	DueDateRaw     string
	CostOfDelayRaw string
}

func (d TodoFormData) DueDateValue() string {
	if d.DueDateRaw != "" {
		return d.DueDateRaw
	}
	return todoInputDueDateValue(d.Input)
}

func (d TodoFormData) CostOfDelayValue() string {
	if d.CostOfDelayRaw != "" {
		return d.CostOfDelayRaw
	}
	return todoInputCostOfDelayValue(d.Input)
}

func todoInputDueDateValue(in TodoInput) string {
	if in.DueDate == nil {
		return ""
	}
	return in.DueDate.Format("2006-01-02")
}

func todoInputCostOfDelayValue(in TodoInput) string {
	return strconv.FormatInt(int64(in.CostOfDelay), 10)
}

// helper to parse common form fields from a request
func parseTodoForm(r *http.Request) (in TodoInput, err error) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err = r.ParseMultipartForm(1 << 20)
	} else {
		err = r.ParseForm()
	}
	if err != nil {
		return
	}

	in.Short = strings.TrimSpace(r.FormValue("short"))
	if in.Short == "" {
		err = fmt.Errorf("short is required")
		return
	}

	in.Description = strings.TrimSpace(r.FormValue("description"))

	dateStr := strings.TrimSpace(r.FormValue("due_date"))
	if dateStr != "" {
		var dt time.Time
		dt, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			err = fmt.Errorf("due_date must be in YYYY-MM-DD format")
			return
		}
		in.DueDate = &dt
	}

	if codStr := strings.TrimSpace(r.FormValue("cost_of_delay")); codStr != "" {
		var tmp int
		tmp, err = strconv.Atoi(codStr)
		if err != nil {
			err = fmt.Errorf("cost_of_delay must be a number")
			return
		}
		if tmp < -2 || tmp > 2 {
			err = fmt.Errorf("cost_of_delay must be between -2 and 2")
			return
		}
		in.CostOfDelay = int16(tmp)
	}

	in.Effort = strings.TrimSpace(r.FormValue("effort"))
	switch in.Effort {
	case "mins", "hours", "days", "weeks", "months":
	case "":
		in.Effort = "hours"
	default:
		err = fmt.Errorf("invalid effort")
		return
	}

	return
}
