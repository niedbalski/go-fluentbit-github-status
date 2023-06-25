package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/calyptia/plugin"
)

const (
	GithubStatusEndpointJSON = "https://www.githubstatus.com/api/v2/status.json"
)

type GithubStatusResponse struct {
	Page struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		URL       string    `json:"url"`
		TimeZone  string    `json:"time_zone"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"page"`
	Status struct {
		Indicator   string `json:"indicator"`
		Description string `json:"description"`
	} `json:"status"`
}

// Plugin needs to be registered as an input type plugin in the initialisation phase
func init() {
	plugin.RegisterInput("go-fluentbit-github-status", "Golang input plugin for checking the status of github", &GithubStatusPlugin{})
}

type GithubStatusPlugin struct{}

func (plug *GithubStatusPlugin) Init(ctx context.Context, fbit *plugin.Fluentbit) error {
	return nil
}

func (plug *GithubStatusPlugin) Collect(ctx context.Context, ch chan<- plugin.Message) error {
	tick := time.NewTicker(time.Second * 5)

	fmt.Println("initialized")
	for {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}

			return nil
		case <-tick.C:
			resp, err := http.Get(GithubStatusEndpointJSON)
			if err != nil {
				log.Fatal(err)
			}

			responseData, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			var githubStatusResponse GithubStatusResponse
			err = json.Unmarshal(responseData, &githubStatusResponse)
			if err != nil {
				log.Fatal(err)
			}

			ch <- plugin.Message{
				Time: time.Now(),
				Record: map[string]any{
					"data": githubStatusResponse,
				},
			}
		}
	}
}

func main() {}
