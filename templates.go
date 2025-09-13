package main

import (
	"html/template"
	"path/filepath"
)

func GetTemplateFile(templateName string) (*template.Template, error) {
	htmlTemplate, err := template.ParseFiles(filepath.Join("web/template", templateName+".html"))
	if err != nil {
		return &template.Template{}, err
	}

	return htmlTemplate, nil
}
