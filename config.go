package main

import (
	"errors"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

// constants
const configFileName = "config.yaml"

// finals (lists)
var AvailableTemplates = []string{"logging"}

// errs
var ErrRouteAlreadyExists = errors.New("route already exists")
var ErrTemplateDoesntExist = errors.New("template doesnt exists")
var ErrTemplateAlreadyActive = errors.New("template already active")

type Config struct {
	Routes      []string `yaml:"routes"`
	Templates   []string `yaml:"templates"`
	Middlewares []string `yaml:"middlewares"`
}

func LoadConfig() (Config, error) {
	contents, err := os.ReadFile(configFileName)

	if err != nil {
		return Config{}, err
	}

	var values Config
	err = yaml.Unmarshal(contents, &values)

	if err != nil {
		return Config{}, err
	}

	return values, nil
}

func AddToRoutes(value string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	if slices.Contains(config.Routes, value) {
		return ErrRouteAlreadyExists
	}

	config.Routes = append(config.Routes, value)

	err = remarshalConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func AddTemplate(name string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	if !slices.Contains(AvailableTemplates, name) {
		return ErrTemplateDoesntExist
	}

	if slices.Contains(config.Templates, name) {
		return ErrTemplateAlreadyActive
	}

	config.Templates = append(config.Templates, name)

	err = remarshalConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func remarshalConfig(newConfig Config) error {
	// marshal it back
	updatedContents, err := yaml.Marshal(newConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFileName, updatedContents, 0644)
	if err != nil {
		return err
	}

	return nil
}
