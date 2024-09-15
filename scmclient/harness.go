package scmclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhruv1397/pr-monitor/types"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type HarnessSCMClient struct {
	httpClient        *http.Client
	pat               string
	accountIdentifier string
	host              string
}

func NewHarnessSCMClient(host string, pat string) (*HarnessSCMClient, error) {
	parts := strings.Split(pat, ".")
	return &HarnessSCMClient{
		accountIdentifier: parts[1],
		httpClient:        http.DefaultClient,
		pat:               pat,
		host:              host,
	}, nil
}

func (h *HarnessSCMClient) GetUser(ctx context.Context) (*types.User, error) {
	email, err := h.getEmail(ctx)
	if err != nil {
		return nil, err
	}
	id, err := h.getPrincipalID(ctx, email)
	if err != nil {
		return nil, err
	}

	return &types.User{
		PrincipalID: id,
		Email:       email,
		PAT:         h.pat,
	}, nil
}

func (h *HarnessSCMClient) GetRepos(ctx context.Context) ([]*types.Repo, error) {
	orgs, err := h.getOrgs(ctx)
	if err != nil {
		return nil, err
	}

	var allRepos []*types.Repo
	var repoMutex sync.Mutex

	var wg sync.WaitGroup
	errChan := make(chan error, len(orgs)*200)
	repoChan := make(chan *types.Repo, len(orgs)*200)

	for _, org := range orgs {
		wg.Add(1)
		go func(org string) {
			defer wg.Done()

			projects, err := h.getProjects(ctx, org)
			if err != nil {
				errChan <- err
				return
			}

			var projectWg sync.WaitGroup
			for _, project := range projects {
				projectWg.Add(1)
				go func(project string) {
					defer projectWg.Done()

					repos, err := h.getRepos(ctx, org, project)
					if err != nil {
						errChan <- err
						return
					}

					for _, repo := range repos {
						currentRepo := &types.Repo{
							AccountIdentifier: h.accountIdentifier,
							OrgIdentifier:     org,
							ProjectIdentifier: project,
							RepoIdentifier:    repo,
						}
						repoChan <- currentRepo
					}
				}(project)
			}

			projectWg.Wait()
		}(org)
	}

	go func() {
		wg.Wait()
		close(repoChan)
		close(errChan)
	}()

	var errs []error

	for repoChan != nil || errChan != nil {
		select {
		case resp, ok := <-repoChan:
			if !ok {
				repoChan = nil
			} else {
				repoMutex.Lock()
				allRepos = append(allRepos, resp)
				repoMutex.Unlock()
			}
		case errValue, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				errs = append(errs, errValue)
			}
		}
	}

	if len(errs) > 0 {
		return allRepos, fmt.Errorf("errors encountered: %v", errs)
	}

	return allRepos, nil
}

func (h *HarnessSCMClient) getEmail(ctx context.Context) (string, error) {
	apiURL := fmt.Sprintf("%s%s%s", h.host, "/ng/api/user/currentUser?accountIdentifier=", h.accountIdentifier)
	responseObj := &types.EmailResponse{}
	err := h.get(ctx, apiURL, responseObj)
	if err != nil {
		return "", err
	}
	return responseObj.EmailData.Email, nil
}

func (h *HarnessSCMClient) getPrincipalID(ctx context.Context, email string) (int64, error) {
	apiURL := fmt.Sprintf("%s%s%s%s%s", h.host, "/gateway/code/api/v1/principals?query=", email,
		"&type=user&accountIdentifier=", h.accountIdentifier)
	responseObj := make([]*types.PrincipalData, 0)
	err := h.get(ctx, apiURL, &responseObj)
	if err != nil {
		return 0, err
	}

	return responseObj[0].ID, nil
}

func (h *HarnessSCMClient) getOrgs(ctx context.Context) ([]string, error) {
	var orgs = make([]string, 0)
	apiURL := fmt.Sprintf("%s%s", h.host, "/v1/orgs?page=0&limit=200&sort=name&order=ASC")
	responseObj := make([]*types.OrgResponse, 0)
	err := h.get(ctx, apiURL, &responseObj)
	if err != nil {
		return orgs, err
	}
	for _, org := range responseObj {
		orgs = append(orgs, org.OrgData.Identifier)
	}
	return orgs, nil
}

func (h *HarnessSCMClient) getProjects(ctx context.Context, org string) ([]string, error) {
	var projects = make([]string, 0)
	apiURL := fmt.Sprintf("%s%s%s%s", h.host, "/v1/orgs/", org, "/projects?has_module=true&module_type=CODE&page=0&limit=200&sort=name&order=ASC")
	responseObj := make([]*types.ProjectResponse, 0)
	err := h.get(ctx, apiURL, &responseObj)
	if err != nil {
		return projects, err
	}
	for _, project := range responseObj {
		projects = append(projects, project.ProjectData.Identifier)
	}
	return projects, nil
}

func (h *HarnessSCMClient) getRepos(ctx context.Context, org string, project string) ([]string, error) {
	var repos = make([]string, 0)
	apiURL := fmt.Sprintf("%s%s%s%s%s%s%s%s", h.host, "/code/api/v1/repos?accountIdentifier=", h.accountIdentifier, "&orgIdentifier=", org, "&projectIdentifier=", project, "&page=1&limit=200")
	responseObj := make([]*types.RepoData, 0)
	err := h.get(ctx, apiURL, &responseObj)
	if err != nil {
		return repos, err
	}
	for _, repo := range responseObj {
		repos = append(repos, repo.Identifier)
	}
	return repos, nil
}

func (h *HarnessSCMClient) get(ctx context.Context, url string, v any) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("x-api-key", h.pat)

	response, err := h.httpClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}
