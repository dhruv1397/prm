package clientbuilder

import (
	"github.com/dhruv1397/prm/prclient"
	"github.com/dhruv1397/prm/scmclient"
	"github.com/dhruv1397/prm/types"
)

func GetHarnessSCMClient(host string, pat string) (*scmclient.HarnessSCMClient, error) {
	return scmclient.NewHarnessSCMClient(host, pat)
}

func GetHarnessPRClient(host string, user *types.User, repos []*types.Repo) (prclient.PRClient, error) {
	return prclient.NewHarnessPRClient(host, user, repos)
}
