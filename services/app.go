package services

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/chocnut/sentry-api/domain"
	humanize "github.com/dustin/go-humanize"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
)

var limit int
var sortBy string
var period string
var env string

func initFlags() {
	flag.IntVar(&limit, "limit", 25, "Result limit")
	flag.IntVar(&limit, "l", 25, "Result limit (shorthand)")

	flag.StringVar(&sortBy, "sort", "freq", "Sort by")
	flag.StringVar(&sortBy, "s", "freq", "Sort by (shorthand)")

	flag.StringVar(&period, "period", "14d", "Period")
	flag.StringVar(&period, "p", "14d", "Period (shorthand)")

	flag.StringVar(&env, "env", "production", "Environment")
	flag.StringVar(&env, "e", "production", "Environment (shorthand)")

	flag.Parse()
}

/*
Run ...
Bootstrap
*/
func Run() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initFlags()

	sentryURL := os.Getenv("SENTRY_URL")
	token := os.Getenv("SENTRY_TOKEN")

	sentryURL += fmt.Sprintf("/issues/?query=environment:%s+is:unresolved&sort=%s&statsPeriod=%s&limit=%d", env, sortBy, period, limit)
	req, error := http.NewRequest("GET", sentryURL, nil)

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: time.Second * 40,
	}
	resp, error := client.Do(req)

	if error != nil {
		fmt.Println("oops")
	}

	decoder := json.NewDecoder(resp.Body)
	var data []domain.Issue

	error = decoder.Decode(&data)
	if error != nil {
		fmt.Println(error)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[j].UserCount < data[i].UserCount
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Event Count", "User Count", "Last Seen", "Title", "Permalink"})

	for _, issue := range data {
		i64, _ := strconv.ParseInt(issue.Count, 10, 32)
		lastSeen, _ := time.Parse(time.RFC3339, issue.LastSeen)
		row := []string{issue.ID, humanize.Comma(i64), humanize.Comma(issue.UserCount), humanize.Time(lastSeen), issue.Title, issue.Permalink}

		if issue.UserCount >= 5000 {
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}})
		} else if issue.UserCount >= 1000 && issue.UserCount <= 5000 {
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}})
		} else {
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}})
		}

	}
	table.Render()
}
