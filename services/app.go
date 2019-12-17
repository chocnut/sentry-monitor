package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/chocnut/sentry-api/domain"
	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
)

var sentryURL = "https://sentry.infostreamgroup.com/api/0/projects/reflex/sa_web/issues/?query=environment:production+is:unresolved&sort=freq&statsPeriod=14d&limit=25"

// var sentryURL = "https://sentry.infostreamgroup.com/api/0/projects/reflex/wypv3_webapp/issues/?query=environment:production+is:unresolved&sort=freq&statsPeriod=14d&limit=25"

var bearer = "Bearer 522a48dcfa15484ab5f540864b71ccfd46bfa41794fd440fa3a4d947b9c8717b"

// var bearer = "Bearer 28ae8bfebc444d61996fa694c347bbd3d42fd504983f493bb1b0e6db175ab5b1"

/*
Run ...
bootstrap
*/
func Run() {

	req, error := http.NewRequest("GET", sentryURL, nil)

	req.Header.Add("Authorization", bearer)

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
	table.SetHeader([]string{"ID", "Event Count", "User Count", "Title", "Permalink"})

	for _, issue := range data {
		i64, _ := strconv.ParseInt(issue.Count, 10, 32)
		row := []string{issue.ID, humanize.Comma(i64), humanize.Comma(issue.UserCount), issue.Title, issue.Permalink}

		if issue.UserCount >= 5000 {
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}})
		} else if issue.UserCount >= 1000 && issue.UserCount <= 5000 {
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}})
		} else {
			table.Rich(row, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlueColor}})
		}

	}
	table.Render()
}
