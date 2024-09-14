package scmclient

import (
	"context"
	"github.com/dhruv1397/pr-monitor/types"
	"net/http"
)

type HarnessSCMClient struct {
	httpClient        *http.Client
	pat               string
	accountIdentifier string
	host              string
}

func NewHarnessSCMClient(pat string) (*HarnessSCMClient, error) {
	return &HarnessSCMClient{
		httpClient: http.DefaultClient,
		pat:        pat,
	}, nil
}

func (h *HarnessSCMClient) GetUser(ctx context.Context) (*types.HarnessUser, error) {

}

func (h *HarnessSCMClient) GetRepos(ctx context.Context, user *types.HarnessUser) ([]*types.HarnessRepo, error) {

}
