package timesheet

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	BulletList []string

	Record struct {
		Date        string     `json:"date"`
		Hours       string     `json:"hours"`
		Description BulletList `json:"description"`
		Times       BulletList `json:"times"`
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
	for _, v := range s.Records {
		f, err := strconv.ParseFloat(v.Hours, 64)
		if err != nil {
			fmt.Println("Hours provided aren't a number")
			os.Exit(1)
		}
		sum += f
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
