package main

import (
	"errors"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

// constants
const configFileName = "config.yaml"

// finals
var AvailableTemplates = []string{"home", "data-dashboard", "video-list", "greeting", "profile", "logging"}

// errs
var ErrRouteAlreadyExists = errors.New("route already exists")
var ErrTemplateDoesntExist = errors.New("template doesnt exists")
var ErrTemplateAlreadyActive = errors.New("template already active")

type Route struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

type TemplateData struct {
	Name string         `yaml:"name"`
	Data map[string]any `yaml:"data,omitempty"`
}

type Config struct {
	Routes       []Route        `yaml:"routes"`
	Templates    []string       `yaml:"templates"`
	TemplateData []TemplateData `yaml:"template_data,omitempty"`
	Middlewares  []string       `yaml:"middlewares"`
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

	// check if template exists
	if !slices.Contains(AvailableTemplates, template) {
		return ErrTemplateDoesntExist
	}

	newRoute := Route{
		Path:     path,
		Template: template,
	}

	config.Routes = append(config.Routes, newRoute)

	// add it not already there
	if !slices.Contains(config.Templates, template) {
		config.Templates = append(config.Templates, template)
	}

	// add it only if it isnt there
	if !hasTemplateData(config.TemplateData, template) {
		config.TemplateData = append(config.TemplateData, TemplateData{
			Name: template,
			Data: make(map[string]interface{}),
		})
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

	// add it
	if !hasTemplateData(config.TemplateData, name) {
		config.TemplateData = append(config.TemplateData, TemplateData{
			Name: name,
			Data: make(map[string]interface{}),
		})
	}

	err = remarshalConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// set data for existing template
func SetTemplateData(templateName string, data map[string]any) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Check if template exists in available templates
	if !slices.Contains(AvailableTemplates, templateName) {
		return ErrTemplateDoesntExist
	}

	// Find existing template data or create new one
	found := false
	for i, td := range config.TemplateData {
		if td.Name == templateName {
			config.TemplateData[i].Data = data
			found = true
			break
		}
	}

	if !found {
		config.TemplateData = append(config.TemplateData, TemplateData{
			Name: templateName,
			Data: data,
		})
	}

	// Ensure template is in active templates list
	if !slices.Contains(config.Templates, templateName) {
		config.Templates = append(config.Templates, templateName)
	}

	err = remarshalConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// retrieve data for an existing template
func GetTemplateData(templateName string) (map[string]any, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	for _, td := range config.TemplateData {
		if td.Name == templateName {

			if td.Name == "home" {
				links := []string{}
				for _, route := range config.Routes {
					if len(route.Path) > 1 {
						links = append(links, route.Path[1:]) // omit first character
					}
				}

				td.Data["AvailableLinks"] = links

				return td.Data, nil
			}

			return td.Data, nil

		}
	}

	return nil, errors.New("template data not found")
}

// updates specific fields in template data
func UpdateTemplateData(templateName string, updates map[string]any) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// check if exists
	if !slices.Contains(AvailableTemplates, templateName) {
		return ErrTemplateDoesntExist
	}

	// Find existing template data
	found := false
	for i, td := range config.TemplateData {
		if td.Name == templateName {
			if td.Data == nil {
				td.Data = make(map[string]any)
			}
			// Update existing data with new values
			for key, value := range updates {
				td.Data[key] = value
			}
			config.TemplateData[i] = td
			found = true
			break
		}
	}

	if !found {
		// Create new template data entry
		config.TemplateData = append(config.TemplateData, TemplateData{
			Name: templateName,
			Data: updates,
		})
	}

	err = remarshalConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// ListTemplatesWithData returns all templates and their associated data
func ListTemplatesWithData() ([]TemplateData, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.TemplateData, nil
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

// GetTemplateAndDataForRoute returns both template name and its data for a given route
func GetTemplateAndDataForRoute(path string) (string, map[string]any, error) {
	templateName, err := GetTemplateForRoute(path)
	if err != nil {
		return "", nil, err
	}

	templateData, err := GetTemplateData(templateName)
	if err != nil {
		// Return template name even if data doesn't exist
		return templateName, make(map[string]any), nil
	}

	return templateName, templateData, nil
}

// Helper function to list all routes with their templates
func ListRoutes() ([]Route, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.Routes, nil
}

// Helper function to check if template data exists
func hasTemplateData(templateData []TemplateData, templateName string) bool {
	for _, td := range templateData {
		if td.Name == templateName {
			return true
		}
	}
	return false
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
