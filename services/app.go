package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/chocnut/sentry-api/domain"
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

	fmt.Printf("ID    | User Count| Title | Permalink    | \n")
	for _, issue := range data {
		fmt.Printf("|%s| %-10v| %30v| %-10v|\n", issue.ID, issue.UserCount, issue.Title, issue.Permalink)
	}

	fmt.Println("done!")
}
