package mcduck_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/manzanit0/mcduck/internal/mcduck"
)

func TestCalculateTotalsPerCategory(t *testing.T) {
	testCases := []struct {
		expenses []mcduck.Expense
		result   map[string]map[string]uint64
	}{
		{
			expenses: []mcduck.Expense{
				{Date: time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC), Category: "a", Amount: 1},
				{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "a", Amount: 1},
				{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "b", Amount: 1},
				{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "c", Amount: 1},
				{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "c", Amount: 2},
			},
			result: map[string]map[string]uint64{
				"2006-01": {"a": 1},
				"2006-02": {"a": 1, "b": 1, "c": 3},
			},
		},
		{
			expenses: []mcduck.Expense{
				{Date: time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC), Category: "a", Amount: 13},
				{Date: time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC), Category: "a", Amount: 12},
			},
			result: map[string]map[string]uint64{
				"2006-01": {"a": 13},
				"2008-01": {"a": 12},
			},
		},
	}
	for x, tC := range testCases {
		t.Run(fmt.Sprintf("case %d", x), func(t *testing.T) {
			totals := mcduck.CalculateTotalsPerCategory(tC.expenses)
			if len(totals) != len(tC.result) {
				t.Fatalf("expected %d results, got %d", len(tC.result), len(totals))
			}

			for month, monthTotals := range totals {
				for category, total := range monthTotals {
					if tC.result[month][category] != total {
						t.Errorf("expected %d for %s-%s, got %d", tC.result[month][category], month, category, total)
					}
				}
			}
		})
	}
}

func TestCalculateMonthOverMonthTotals(t *testing.T) {
	input := []mcduck.Expense{
		{Date: time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC), Category: "a", Amount: 1},
		{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "a", Amount: 1},
		{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "b", Amount: 2},
		{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "c", Amount: 3},
		{Date: time.Date(2006, 2, 1, 0, 0, 0, 0, time.UTC), Category: "c", Amount: 3},
		{Date: time.Date(2006, 3, 1, 0, 0, 0, 0, time.UTC), Category: "d", Amount: 4},
	}

	want := map[string]map[string]uint64{
		"a": {"2006-01": 1, "2006-02": 1, "2006-03": 0},
		"b": {"2006-01": 0, "2006-02": 2, "2006-03": 0},
		"c": {"2006-01": 0, "2006-02": 6, "2006-03": 0},
		"d": {"2006-01": 0, "2006-02": 0, "2006-03": 4},
	}

	got := mcduck.CalculateMonthOverMonthTotals(input)

	if len(want) != len(got) {
		t.Fatalf("expected %d results, got %d", len(want), len(got))
	}

	for category, amountsByMonth := range got {
		for month, amount := range amountsByMonth {
			if want[category][month] != amount {
				t.Errorf("wanted %d for %s in %s, got %d", want[category][month], category, month, amount)
			}
		}
	}
}

func TestGetTop3ExpenseCategories(t *testing.T) {
	testCases := []struct {
		desc      string
		input     []mcduck.Expense
		monthYear string
		output    []mcduck.CategoryAggregate
	}{
		{
			desc:      "when less than three categories are provided, then they're all returned",
			monthYear: mcduck.NewMonthYear(time.Date(2008, time.February, 2, 0, 0, 0, 0, time.UTC)),
			input: []mcduck.Expense{
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "foo", Amount: 11},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "foo", Amount: 11},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "bar", Amount: 33},
			},
			output: []mcduck.CategoryAggregate{
				{Category: "bar", MonthYear: "2008-02", TotalAmount: 33},
				{Category: "foo", MonthYear: "2008-02", TotalAmount: 22},
			},
		},
		{
			desc:      "when more than three categories are provided, then only the top three are returned",
			monthYear: mcduck.NewMonthYear(time.Date(2008, time.February, 2, 0, 0, 0, 0, time.UTC)),
			input: []mcduck.Expense{
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "foo", Amount: 11},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "foo", Amount: 11},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "bar", Amount: 33},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "baz", Amount: 102},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "baz", Amount: 44},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "baz", Amount: 55},
				{Date: time.Date(2008, time.February, 11, 0, 0, 0, 0, time.UTC), Subcategory: "nope", Amount: 5},
			},
			output: []mcduck.CategoryAggregate{
				{Category: "baz", MonthYear: "2008-02", TotalAmount: 201},
				{Category: "bar", MonthYear: "2008-02", TotalAmount: 33},
				{Category: "foo", MonthYear: "2008-02", TotalAmount: 22},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			aggregate := mcduck.GetTop3ExpenseCategories(tC.input, tC.monthYear)
			for i := range aggregate {
				if aggregate[i].Category != tC.output[i].Category {
					t.Error("unexpected category", aggregate[i].Category, "expected", tC.output[i].Category)
				}

				if aggregate[i].TotalAmount != tC.output[i].TotalAmount {
					t.Error("unexpected amount", aggregate[i].TotalAmount, "expected", tC.output[i].TotalAmount)
				}

				if aggregate[i].MonthYear != tC.output[i].MonthYear {
					t.Error("unexpected date ", aggregate[i].MonthYear, "expected", tC.output[i].MonthYear)
				}
			}
		})
	}
}

func TestFromCSV(t *testing.T) {
	t.Run("when the file is empty, an error is returned", func(t *testing.T) {
		expenses, err := mcduck.FromCSV(bytes.NewBufferString(""))

		if err == nil {
			t.Fatalf("expected an error, got nil")
		}

		if len(expenses) != 0 {
			t.Fatalf("expected zero expenses, got %v", len(expenses))
		}
	})

	t.Run("when the column separator is a semi-colon, the expenses are parsed successfully", func(t *testing.T) {
		expenses, err := mcduck.FromCSV(bytes.NewBufferString(`
date;amount;category;subcategory
2022-04-02;2.82;food;meat
2022-04-02;8.22;transport;gasoline
`))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(expenses) != 2 {
			t.Fatalf("expected two expenses, got %v", len(expenses))
		}

		e := expenses[0]
		if e.Amount != 282 {
			t.Errorf("expected amount to be 2.82, got %v", e.Amount)
		}

		if e.Date.Format("2006-01-02") != "2022-04-02" {
			t.Errorf("expected date to be 2022-04-02, got %v", e.Date)
		}

		if e.Category != "food" {
			t.Errorf("expected category to be food, got %v", e.Category)
		}

		if e.Subcategory != "meat" {
			t.Errorf("expected subcategory to be meat, got %v", e.Subcategory)
		}

		e = expenses[1]
		if e.Amount != 822 {
			t.Errorf("expected amount to be 8.22, got %v", e.Amount)
		}

		if e.Date.Format("2006-01-02") != "2022-04-02" {
			t.Errorf("expected date to be 2022-04-02, got %v", e.Date)
		}

		if e.Category != "transport" {
			t.Errorf("expected category to be transport, got %v", e.Category)
		}

		if e.Subcategory != "gasoline" {
			t.Errorf("expected subcategory to be gasoline, got %v", e.Subcategory)
		}
	})

	t.Run("when the column separator is a comma, the expenses are parsed successfully", func(t *testing.T) {
		expenses, err := mcduck.FromCSV(bytes.NewBufferString(`
date,amount,category,subcategory
2022-04-02,2.82,food,meat
2022-04-02,8.22,transport,gasoline
`))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(expenses) != 2 {
			t.Fatalf("expected two expenses, got %v", len(expenses))
		}

		e := expenses[0]
		if e.Amount != 282 {
			t.Errorf("expected amount to be 2.82, got %v", e.Amount)
		}

		if e.Date.Format("2006-01-02") != "2022-04-02" {
			t.Errorf("expected date to be 2022-04-02, got %v", e.Date)
		}

		if e.Category != "food" {
			t.Errorf("expected category to be food, got %v", e.Category)
		}

		if e.Subcategory != "meat" {
			t.Errorf("expected subcategory to be meat, got %v", e.Subcategory)
		}

		e = expenses[1]
		if e.Amount != 822 {
			t.Errorf("expected amount to be 8.22, got %v", e.Amount)
		}

		if e.Date.Format("2006-01-02") != "2022-04-02" {
			t.Errorf("expected date to be 2022-04-02, got %v", e.Date)
		}

		if e.Category != "transport" {
			t.Errorf("expected category to be transport, got %v", e.Category)
		}

		if e.Subcategory != "gasoline" {
			t.Errorf("expected subcategory to be gasoline, got %v", e.Subcategory)
		}
	})

	t.Run("when the column separator is neither a comma nor a semi-colon, an error is returned", func(t *testing.T) {
		expenses, err := mcduck.FromCSV(bytes.NewBufferString(`
date?amount?category?subcategory
2022-04-02?2.82?food?meat
2022-04-02?8.22?transport?gasoline
`))

		if err == nil {
			t.Fatalf("expected an error, got nil")
		}

		if len(expenses) != 0 {
			t.Fatalf("expected zero expenses, got %v", len(expenses))
		}
	})

	t.Run("when the amounts floating point separator is a comma, the expenses are parsed succesfully", func(t *testing.T) {
		expenses, err := mcduck.FromCSV(bytes.NewBufferString(`
date;amount;category;subcategory
2022-04-02;2,82;food;meat
2022-04-02;8,22;transport;gasoline
`))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(expenses) != 2 {
			t.Fatalf("expected two expenses, got %v", len(expenses))
		}

		e := expenses[0]
		if e.Amount != 282 {
			t.Errorf("expected amount to be 2.82, got %v", e.Amount)
		}

		if e.Date.Format("2006-01-02") != "2022-04-02" {
			t.Errorf("expected date to be 2022-04-02, got %v", e.Date)
		}

		if e.Category != "food" {
			t.Errorf("expected category to be food, got %v", e.Category)
		}

		if e.Subcategory != "meat" {
			t.Errorf("expected subcategory to be meat, got %v", e.Subcategory)
		}

		e = expenses[1]
		if e.Amount != 822 {
			t.Errorf("expected amount to be 8.22, got %v", e.Amount)
		}

		if e.Date.Format("2006-01-02") != "2022-04-02" {
			t.Errorf("expected date to be 2022-04-02, got %v", e.Date)
		}

		if e.Category != "transport" {
			t.Errorf("expected category to be transport, got %v", e.Category)
		}

		if e.Subcategory != "gasoline" {
			t.Errorf("expected subcategory to be gasoline, got %v", e.Subcategory)
		}
	})
}
