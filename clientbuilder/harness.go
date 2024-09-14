package clientbuilder

import (
	"github.com/dhruv1397/pr-monitor/prclient"
	"github.com/dhruv1397/pr-monitor/scmclient"
	"github.com/dhruv1397/pr-monitor/types"
)

func GetHarnessSCMClient(pat string) (*scmclient.HarnessSCMClient, error) {
	return scmclient.NewHarnessSCMClient(pat)
}

func GetHarnessPRClient(user *types.HarnessUser, repo *types.HarnessRepo) (prclient.PRClient, error) {
	return prclient.NewHarnessPRClient(user, repo)
}
