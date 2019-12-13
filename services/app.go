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
	"github.com/olekukonko/tablewriter"
)

var sentryURL = "https://sentry.infostreamgroup.com/api/0/projects/reflex/sa_web/issues/?query=environment:production+is:unresolved&sort=freq&statsPeriod=14d&limit=25"

var bearer = "Bearer 522a48dcfa15484ab5f540864b71ccfd46bfa41794fd440fa3a4d947b9c8717b"

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
	table.SetHeader([]string{"ID", "User Count", "Title", "Permalink"})

	for _, issue := range data {
		table.Append([]string{issue.ID, strconv.Itoa(issue.UserCount), issue.Title, issue.Permalink})
	}
	table.Render()
}
