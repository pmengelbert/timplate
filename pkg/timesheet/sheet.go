package timesheet

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type (
	Time struct {
		Hour, Minute int8
	}

	BulletList []string

	Record struct {
		Date        string     `json:"date"`
		Hours       string     `json:"hours"`
		Description BulletList `json:"description"`
		Times       BulletList `json:"times"`
		TimeSum     float32
	}

	Sheet struct {
		Name      string   `json:"name"`
		Rate      int      `json:"rate"`
		TaxRate   float64  `json:"taxRate"`
		StartDate string   `json:"startDate"`
		EndDate   string   `json:"endDate"`
		Records   []Record `json:"records"`
	}
)

func (s Sheet) TotalHours() float64 {
	var sum float64
	for _, r := range s.Records {
		sum += float64(r.TimeSum)
	}
	return sum
}

func (s Sheet) TotalPay() float64 {
	return s.TotalHours() * float64(s.Rate)
}

func (s Sheet) AfterTax() float64 {
	return s.TotalPay() * (1 - s.TaxRate)
}

func (s *Sheet) CapitalizeDescriptions() {
	for i := range s.Records {
		for j, v := range s.Records[i].Description {
			a := strings.Split(v, " ")
			a[0] = strings.Title(a[0])
			str := strings.Join(a, " ")
			s.Records[i].Description[j] = str
		}
	}
}

func Parse(t string) (Time, error) {
	yamlError := fmt.Errorf("bad time in yaml")
	match := regexp.MustCompile("([0-9]{1,2})(:?)([0-9]{2})\\s*((a|p)m?)?")
	if match.FindString(t) != t {
		return Time{}, yamlError
	}
	hourString := match.ReplaceAllString(t, "$1")
	minuteString := match.ReplaceAllString(t, "$3")
	suffix := match.ReplaceAllString(t, "$5")

	var hours, minutes int8
	switch suffix {
	case "p":
		hours += 12
	case "a":
		if hours > 12 {
			return Time{}, yamlError
		}
	}

	i, err := strconv.Atoi(hourString)
	if err != nil {
		return Time{}, err
	}

	j, err := strconv.Atoi(minuteString)
	if err != nil {
		return Time{}, err
	}

	hours += int8(i)
	minutes += int8(j)

	if hours > 23 || minutes > 59 {
		return Time{}, yamlError
	}

	return Time{hours, minutes}, nil
}

func (t1 Time) DifferenceInHours(t2 Time) float32 {
	minDiff := t1.Minute - t2.Minute
	if t1.Hour < t2.Hour {
		t1.Hour += 24
	}

	if minDiff < 0 {
		minDiff += 60
		t1.Hour--
	}

	hrDiff := t1.Hour - t2.Hour
	minFraction := float32(minDiff) / 60.0
	return float32(hrDiff) + minFraction
}
