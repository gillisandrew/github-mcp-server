package cmd

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	gogithub "github.com/google/go-github/v69/github"
)

func getGitHubClient(host, token string) (*gogithub.Client, error) {
	httpClient, err := api.NewHTTPClient(api.ClientOptions{
		Host:      host,
		AuthToken: token,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := gogithub.NewClient(httpClient)

	return client, nil

}
