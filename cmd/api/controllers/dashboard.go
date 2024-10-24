package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"

	"github.com/manzanit0/mcduck/internal/mcduck"
	"github.com/manzanit0/mcduck/pkg/auth"
	"github.com/manzanit0/mcduck/pkg/xtrace"
)

var chartColours = []string{
	"rgba(255, 99, 132)",
	"rgba(255, 159, 64)",
	"rgba(255, 205, 86)",
	"rgba(75, 192, 192)",
	"rgba(54, 162, 235)",
	"rgba(153, 102, 255)",
	"rgba(201, 203, 207)",
	// repeated
	"rgba(255, 99, 132)",
	"rgba(255, 159, 64)",
	"rgba(255, 205, 86)",
	"rgba(75, 192, 192)",
	"rgba(54, 162, 235)",
	"rgba(153, 102, 255)",
	"rgba(201, 203, 207)",
}

var chartBackgroundColours = []string{
	"rgba(255, 99, 132, 0.2)",
	"rgba(255, 159, 64, 0.2)",
	"rgba(255, 205, 86, 0.2)",
	"rgba(75, 192, 192, 0.2)",
	"rgba(54, 162, 235, 0.2)",
	"rgba(153, 102, 255, 0.2)",
	"rgba(201, 203, 207, 0.2)",
	// repeated
	"rgba(255, 99, 132, 0.2)",
	"rgba(255, 159, 64, 0.2)",
	"rgba(255, 205, 86, 0.2)",
	"rgba(75, 192, 192, 0.2)",
	"rgba(54, 162, 235, 0.2)",
	"rgba(153, 102, 255, 0.2)",
	"rgba(201, 203, 207, 0.2)",
}

type ChartData struct {
	Title    string
	Labels   []string
	Datasets []Dataset
}

type Dataset struct {
	Label            string
	BorderColour     string
	BackgroundColour string
	Hidden           bool
	Data             []string
}

type DashboardController struct {
	Expenses   *mcduck.ExpenseRepository
	SampleData []mcduck.Expense
}

func (d *DashboardController) LiveDemo(c *gin.Context) {
	expenses := d.SampleData

	mcduck.SortByDate(expenses)

	mostRecent := mcduck.FindMostRecentTime(expenses)
	mostRecentMonthYear := mcduck.NewMonthYear(mostRecent)

	categoryTotals := mcduck.CalculateTotalsPerCategory(expenses)
	categoryLabels := getSecondClassifier(categoryTotals)
	categoryChartData := buildChartData(categoryLabels, categoryTotals)

	// Since this is for public demoing, we might as well show-off the whole data
	// off the bat.
	for i := range categoryChartData.Datasets {
		categoryChartData.Datasets[i].Hidden = false
	}

	var subcategoryCharts []ChartData
	for cat, subcats := range GroupSubcategoriesByCategory(expenses) {
		filtered := FilterByCategory(expenses, cat)
		subcategoryTotals := mcduck.CalculateTotalsPerSubCategory(filtered)
		subcategoryChartData := buildChartData(subcats, subcategoryTotals)

		subcategoryChartData.Title = cat

		subcategoryCharts = append(subcategoryCharts, subcategoryChartData)
	}

	totalSpendsArr := TotalSpendLastThreeMonths(expenses)
	top3Categories := mcduck.GetTop3ExpenseCategories(expenses, mostRecentMonthYear)

	var formattedTop3Categories []FormattedCategoryAggregate
	for _, t := range top3Categories {
		formattedTop3Categories = append(formattedTop3Categories, FormattedCategoryAggregate{
			Category:    t.Category,
			MonthYear:   t.MonthYear,
			TotalAmount: fmt.Sprintf("%0.2f", mcduck.ConvertToDollar(t.TotalAmount)),
		})
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"PrettyMonthYear":        mostRecent.Format("January 2006"),
		"NoExpenses":             len(expenses) == 0,
		"Categories":             categoryLabels,
		"CategoriesChartData":    categoryChartData,
		"SubcategoriesChartData": subcategoryCharts,
		"TopCategories":          formattedTop3Categories,
		"TotalSpends":            totalSpendsArr,
	})
}

type FormattedCategoryAggregate struct {
	Category    string
	MonthYear   string
	TotalAmount string
}

func (d *DashboardController) Dashboard(c *gin.Context) {
	ctx, span := xtrace.GetSpan(c.Request.Context())
	user := auth.GetUserEmail(c)

	expenses, err := d.Expenses.ListExpenses(ctx, user)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed list expenses", "error", err.Error())
		c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
		return
	}

	if len(expenses) == 0 {
		expenses = []mcduck.Expense{}
	}

	mcduck.SortByDate(expenses)

	mostRecent := mcduck.FindMostRecentTime(expenses)
	mostRecentMonthYear := mcduck.NewMonthYear(mostRecent)

	categoryTotals := mcduck.CalculateTotalsPerCategory(expenses)
	categoryLabels := getSecondClassifier(categoryTotals)
	categoryChartData := buildChartData(categoryLabels, categoryTotals)

	var subcategoryCharts []ChartData
	for cat, subcats := range GroupSubcategoriesByCategory(expenses) {
		filtered := FilterByCategory(expenses, cat)
		subcategoryTotals := mcduck.CalculateTotalsPerSubCategory(filtered)
		subcategoryChartData := buildChartData(subcats, subcategoryTotals)

		subcategoryChartData.Title = cat

		subcategoryCharts = append(subcategoryCharts, subcategoryChartData)
	}

	totalSpendsArr := TotalSpendLastThreeMonths(expenses)

	top3Categories := mcduck.GetTop3ExpenseCategories(expenses, mostRecentMonthYear)

	var formattedTop3Categories []FormattedCategoryAggregate
	for _, t := range top3Categories {
		formattedTop3Categories = append(formattedTop3Categories, FormattedCategoryAggregate{
			Category:    t.Category,
			MonthYear:   t.MonthYear,
			TotalAmount: fmt.Sprintf("%0.2f", mcduck.ConvertToDollar(t.TotalAmount)),
		})
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"PrettyMonthYear":        mostRecent.Format("January 2006"),
		"NoExpenses":             len(expenses) == 0,
		"Categories":             categoryLabels,
		"CategoriesChartData":    categoryChartData,
		"SubcategoriesChartData": subcategoryCharts,
		"TopCategories":          formattedTop3Categories,
		"TotalSpends":            totalSpendsArr,
		"User":                   user,
	})
}

type MonthlySpend struct {
	date      time.Time
	amount    uint64
	MonthYear string
	Amount    string
}

func TotalSpendLastThreeMonths(expenses []mcduck.Expense) []*MonthlySpend {
	latest := mcduck.FindMostRecentTime(expenses)
	totalSpends := map[string]*MonthlySpend{}
	for i := range expenses {
		if isOlderThanLastThreeMonths(expenses[i].Date, latest) {
			continue
		}

		key := expenses[i].Date.Format("January 2006")
		val, ok := totalSpends[key]
		if !ok {
			totalSpends[key] = &MonthlySpend{
				date:      expenses[i].Date,
				MonthYear: key,
				amount:    expenses[i].Amount,
				Amount:    fmt.Sprintf("%.2f", mcduck.ConvertToDollar(expenses[i].Amount)),
			}
		} else {
			val.amount += expenses[i].Amount
			val.Amount = fmt.Sprintf("%.2f", mcduck.ConvertToDollar(val.amount))
		}
	}

	sortedTotalSpends := []*MonthlySpend{}
	for _, a := range totalSpends {
		sortedTotalSpends = append(sortedTotalSpends, a)
	}

	sort.Slice(sortedTotalSpends, func(i, j int) bool {
		return sortedTotalSpends[i].date.Before(sortedTotalSpends[j].date)
	})

	return sortedTotalSpends
}

func isOlderThanLastThreeMonths(t time.Time, latest time.Time) bool {
	// 15th of March 2022-> 15th of December 2022
	year, month, _ := latest.AddDate(0, -2, 0).Date()

	// 1st of December 2022
	beginningOf3MonthsAgo := time.Date(year, month, 1, 0, 0, 0, 0, time.Now().Location())

	return t.Before(beginningOf3MonthsAgo)
}

func FilterByCategory(list []mcduck.Expense, cat string) []mcduck.Expense {
	var filtered []mcduck.Expense
	for i := range list {
		if list[i].Category == cat {
			filtered = append(filtered, list[i])
		}
	}

	return filtered
}

func GroupSubcategoriesByCategory(list []mcduck.Expense) map[string][]string {
	m := map[string]map[string]bool{}
	for _, e := range list {
		if _, ok := m[e.Category]; !ok {
			m[e.Category] = map[string]bool{}
		}

		m[e.Category][e.Subcategory] = true
	}

	mm := map[string][]string{}
	for k, v := range m {
		if _, ok := mm[k]; !ok {
			mm[k] = []string{}
		}

		for s := range v {
			mm[k] = append(mm[k], s)
		}
	}

	return mm
}

func buildChartData(labels []string, totals map[string]map[string]uint64) ChartData {
	var datasets []Dataset
	for monthYear, amountsByCategory := range totals { // totalsByMonth[monthYear][expense.Category] += expense.Amount
		var data []string
		for _, label := range labels {
			if amount, ok := amountsByCategory[label]; ok {
				data = append(data, fmt.Sprintf("%.2f", mcduck.ConvertToDollar(amount)))
			} else {
				data = append(data, "0.00")
			}
		}

		datasets = append(datasets, Dataset{
			Label:  monthYear,
			Data:   data,
			Hidden: true,
		})
	}

	// FIXME: very naive sort. We would want to do a time comparison.
	sort.Slice(datasets, func(i, j int) bool {
		return datasets[i].Label < datasets[j].Label
	})

	// By default we only show the current month.
	if len(datasets) > 0 {
		datasets[len(datasets)-1].Hidden = false
	}

	for i := range datasets {
		datasets[i].BorderColour = chartColours[i]
		datasets[i].BackgroundColour = chartBackgroundColours[i]
	}

	return ChartData{
		Labels:   labels,
		Datasets: datasets,
	}
}

func (d *DashboardController) UploadExpenses(c *gin.Context) {
	ctx, span := xtrace.GetSpan(c.Request.Context())

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "get form error: %s", err.Error())
		return
	}

	if len(form.File["files"]) == 0 {
		c.String(http.StatusBadRequest, "no files uploaded")
		return
	}

	file := form.File["files"][0]
	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed upload file", "error", err.Error())
		c.String(http.StatusInternalServerError, "upload file error: %s", err.Error())
		return
	}

	expenses, err := readExpensesFromCSV(filename)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed read expenses from CSV", "error", err.Error())
		c.String(http.StatusInternalServerError, "file parsing error: %s", err.Error())
		return
	}

	// If the user is logged in, save those upload expenses
	user := auth.GetUserEmail(c)
	if user != "" {
		err = d.Expenses.CreateExpenses(c.Request.Context(), mcduck.ExpensesBatch{UserEmail: user, Records: expenses})
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			slog.ErrorContext(ctx, "failed create expenses", "error", err.Error())
			c.HTML(http.StatusOK, "error.html", gin.H{"error": err.Error()})
			return
		}
	}

	mcduck.SortByDate(expenses)

	mostRecent := mcduck.FindMostRecentTime(expenses)
	mostRecentMonthYear := mcduck.NewMonthYear(mostRecent)

	categoryTotals := mcduck.CalculateTotalsPerCategory(expenses)
	categoryLabels := getSecondClassifier(categoryTotals)
	categoryChartData := buildChartData(categoryLabels, categoryTotals)

	subcategoryTotals := mcduck.CalculateTotalsPerSubCategory(expenses)
	subcategoryLabels := getSecondClassifier(subcategoryTotals)
	subcategoryChartData := buildChartData(subcategoryLabels, subcategoryTotals)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"PrettyMonthYear":        mostRecent.Format("January 2006"),
		"NoExpenses":             len(expenses) == 0,
		"Categories":             categoryLabels,
		"CategoriesChartData":    categoryChartData,
		"SubCategories":          subcategoryLabels,
		"SubCategoriesChartData": subcategoryChartData,
		"TopCategories":          mcduck.GetTop3ExpenseCategories(expenses, mostRecentMonthYear),
		"User":                   user,
	})
}

func getSecondClassifier(calculations map[string]map[string]uint64) []string {
	classifierMap := map[string]bool{}
	classifierSlice := []string{}
	for _, amountByClassifier := range calculations {
		for secondClassifier := range amountByClassifier {
			if ok := classifierMap[secondClassifier]; !ok {
				classifierMap[secondClassifier] = true
				classifierSlice = append(classifierSlice, secondClassifier)
			}
		}
	}

	sort.Strings(classifierSlice)
	return classifierSlice
}

func readExpensesFromCSV(filename string) ([]mcduck.Expense, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return mcduck.FromCSV(f)
}
