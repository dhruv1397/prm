package clientbuilder

import (
	"github.com/dhruv1397/pr-monitor/prclient"
	"github.com/dhruv1397/pr-monitor/scmclient"
	"github.com/dhruv1397/pr-monitor/types"
)

func GetHarnessSCMClient(host string, pat string) (*scmclient.HarnessSCMClient, error) {
	return scmclient.NewHarnessSCMClient(host, pat)
}

func GetHarnessPRClient(host string, user *types.User, repos []*types.Repo) (prclient.PRClient, error) {
	return prclient.NewHarnessPRClient(host, user, repos)
}
