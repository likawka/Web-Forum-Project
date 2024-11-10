package handlers

import (
	"fmt"
	"github.com/0-LY/Forum-test/pkg/api"

	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

func renderTemplate(w http.ResponseWriter, pageName string, data api.ParserConfig) {
	page, ok := api.Pages[pageName]
	if !ok {
		fmt.Println("Page not found:", pageName)
		return
	}

	path := "web/templates/"
	var componentPaths []string

	for _, c := range page.Components {
		componentPaths = append(componentPaths, path+"components/"+c+".html")
	}
	componentPaths = append(componentPaths, path+"base.layout.html")

	templates := append([]string{path + page.Template}, componentPaths...)
	parsedTemplate, err := template.ParseFiles(templates...)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	err = parsedTemplate.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}

func writeStructToFile(data interface{}, name string) error {
	file, err := os.Create("test/" + name + ".json")
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
