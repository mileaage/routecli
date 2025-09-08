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
var AvailableTemplates = []string{"home", "data-dashboard", "video-list", "greeting", "profile", "logging"}

// errs
var ErrRouteAlreadyExists = errors.New("route already exists")
var ErrTemplateDoesntExist = errors.New("template doesnt exists")
var ErrTemplateAlreadyActive = errors.New("template already active")

type Route struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

type Config struct {
	Routes      []Route  `yaml:"routes"`
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

func AddToRoutes(path, template string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// check if already in routes
	for _, route := range config.Routes {
		if route.Path == path {
			return ErrRouteAlreadyExists
		}
	}

	// Check if template exists in available templates
	if !slices.Contains(AvailableTemplates, template) {
		return ErrTemplateDoesntExist
	}

	newRoute := Route{
		Path:     path,
		Template: template,
	}

	config.Routes = append(config.Routes, newRoute)

	// Auto-add template to active templates if not already there
	if !slices.Contains(config.Templates, template) {
		config.Templates = append(config.Templates, template)
	}

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

// Helper function to get template for a given route path
func GetTemplateForRoute(path string) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	for _, route := range config.Routes {
		if route.Path == path {
			return route.Template, nil
		}
	}

	return "", errors.New("route not found")
}

// Helper function to list all routes with their templates
func ListRoutes() ([]Route, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.Routes, nil
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
