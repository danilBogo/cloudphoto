package services

import (
	"cloudphoto/internal/constants"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"os/user"
	"path/filepath"
)

type ConfigManager struct {
	user *user.User
}

type IniConfig struct {
	Bucket      string
	AccessKey   string
	SecretKey   string
	Region      string
	EndpointURL string
}

type keyValue struct {
	key   string
	value string
}

func NewConfigManager() (*ConfigManager, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &ConfigManager{user: currentUser}, nil
}

func (cm ConfigManager) GenerateIni(config *IniConfig) error {
	path := cm.getConfigFilePath()

	cfg := ini.Empty()

	section, err := cfg.NewSection(constants.DefaultSectionName)
	if err != nil {
		return err
	}

	var keyValues = []keyValue{
		{constants.Bucket, config.Bucket},
		{constants.AccessKey, config.AccessKey},
		{constants.SecretKey, config.SecretKey},
		{constants.Region, config.Region},
		{constants.EndpointURL, config.EndpointURL},
	}

	for _, kv := range keyValues {
		_, err = section.NewKey(kv.key, kv.value)
		if err != nil {
			return err
		}
	}

	if !fileExists(path) {
		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}

		file, err := os.Create(path)
		_ = file.Close()
		if err != nil {
			return err
		}
	}

	err = cfg.SaveTo(path)
	if err != nil {
		return err
	}

	fmt.Printf("Ini file successfully created at %v\n", path)

	return nil
}

func (cm ConfigManager) TryGetConfig() (*IniConfig, error) {
	path := cm.getConfigFilePath()

	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	section := cfg.Section(constants.DefaultSectionName)
	config := &IniConfig{
		Bucket:      section.Key(constants.Bucket).String(),
		AccessKey:   section.Key(constants.AccessKey).String(),
		SecretKey:   section.Key(constants.SecretKey).String(),
		Region:      section.Key(constants.Region).String(),
		EndpointURL: section.Key(constants.EndpointURL).String(),
	}

	if !cm.isValidConfig(*config) {
		return nil, errors.New("ini config file is not valid")
	}

	return config, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (cm ConfigManager) getConfigFilePath() string {
	return filepath.Join(cm.user.HomeDir, ".config", "cloudphoto", "cloudphotorc", "cloudphoto.ini")
}

func (cm ConfigManager) isValidConfig(config IniConfig) bool {
	return len(config.Bucket) > 0 &&
		len(config.AccessKey) > 0 &&
		len(config.SecretKey) > 0 &&
		len(config.Region) > 0 &&
		len(config.EndpointURL) > 0
}
