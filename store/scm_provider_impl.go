package store

import (
	"fmt"
	"github.com/dhruv1397/pr-monitor/types"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

var _ SCMProvider = (*scmProviderImpl)(nil)

const configFileName = ".prm_config"

type scmProviderImpl struct {
}

func NewSCMProviderImpl() SCMProvider {
	return &scmProviderImpl{}
}

func (s *scmProviderImpl) Create(provider types.SCMProvider) error {
	existingProviders, err := s.readYAML()
	if err != nil {
		return fmt.Errorf("error listing existing SCM providers before creating new: %w", err)
	}

	if existingProviders[provider.Name] != nil {
		return fmt.Errorf("SCM provider %s already exists", provider.Name)
	}

	provider.Updated = time.Now().UnixMilli()
	provider.Created = time.Now().UnixMilli()

	existingProviders[provider.Name] = &provider

	err = s.writeYAML(existingProviders)
	if err != nil {
		return fmt.Errorf("error writing SCM provider config: %w", err)
	}
	return nil
}

func (s *scmProviderImpl) UpdateBulk(providers []types.SCMProvider) error {
	existingProviders, err := s.readYAML()
	if err != nil {
		return fmt.Errorf("error listing existing SCM providers before updating: %w", err)
	}

	for _, provider := range providers {
		if existingProviders[provider.Name] == nil {
			return fmt.Errorf("SCM provider %s does not exist", provider.Name)
		}
		provider.Updated = time.Now().UnixMilli()
		existingProviders[provider.Name] = &provider
	}

	err = s.writeYAML(existingProviders)
	if err != nil {
		return fmt.Errorf("error writing SCM provider config: %w", err)
	}
	return nil
}

func (s *scmProviderImpl) List(providerType string, providerName string) ([]*types.SCMProvider, error) {
	providerMap, err := s.readYAML()
	if err != nil {
		return nil, fmt.Errorf("error listing SCM providers: %w", err)
	}
	if providerMap == nil {
		return []*types.SCMProvider{}, nil
	}
	providers := make([]*types.SCMProvider, 0)
	for name, provider := range providerMap {
		if name != "" && provider != nil &&
			(providerName == "" || providerName == provider.Name) &&
			(providerType == "" || providerType == provider.Type) {
			providers = append(providers, provider)
		}
	}

	return providers, nil
}

func (s *scmProviderImpl) Delete(name string) error {
	existingProviders, err := s.readYAML()
	if err != nil {
		return fmt.Errorf("error listing existing SCM providers before deleting: %w", err)
	}

	if existingProviders[name] == nil {
		fmt.Printf("SCM provider %s does not exist\n", name)
		return nil
	}

	existingProviders[name] = nil

	err = s.writeYAML(existingProviders)
	if err != nil {
		return fmt.Errorf("error writing SCM provider config: %w", err)
	}
	return nil
}

func (s *scmProviderImpl) getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}
	configFilePath := filepath.Join(homeDir, configFileName)
	return configFilePath, nil
}

func (s *scmProviderImpl) readYAML() (map[string]*types.SCMProvider, error) {
	configFilePath, err := s.getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("error getting SCM provider config file path before reading: %w", err)
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, fmt.Errorf("could not create SCM provider config file %s: %w", configFilePath, err)
		}
		defer file.Close()
	}

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading SCM provider config file %s: %w", configFilePath, err)
	}

	var config = &types.SCMConfig{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("error deserialising SCM provider config: %w", err)
	}
	providerMap := map[string]*types.SCMProvider{}
	for _, provider := range config.Providers {
		providerMap[provider.Name] = provider
	}

	return providerMap, nil
}

func (s *scmProviderImpl) writeYAML(providerMap map[string]*types.SCMProvider) error {
	configFilePath, err := s.getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting SCM provider config file path before writing: %w", err)
	}

	providers := make([]*types.SCMProvider, 0)
	for name, provider := range providerMap {
		if name != "" && provider != nil {
			providers = append(providers, provider)
		}
	}
	config := &types.SCMConfig{Providers: providers}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error serialising SCM provider config: %w", err)
	}

	err = os.WriteFile(configFilePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing SCM provider config: %w", err)
	}

	return nil
}
