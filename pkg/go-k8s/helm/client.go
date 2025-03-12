package helm

import (
	"helm.sh/helm/v3/pkg/cli"
)

type Client struct {
	Settings *cli.EnvSettings
}

func NewClient() *Client {
	return &Client{
		Settings: cli.New(),
	}
}
