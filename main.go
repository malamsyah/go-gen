package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Define the structure of the YAML file
type APIConfig struct {
	Models map[string]map[string]string `yaml:"models"`
	Routes []Route                      `yaml:"routes"`
}

type Field struct {
	Name string `yaml:"Name"`
	Type string `yaml:"Type"`
}

type Route struct {
	Path       string `yaml:"path"`
	Method     string `yaml:"method"`
	Controller string `yaml:"controller"`
	Model      string `yaml:"model"`
}

func main() {
	// Check if the file path argument is provided
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <file_path>", os.Args[0])
	}

	// Get the file path from the first argument
	filePath := os.Args[1]

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	var config APIConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Generate files
	generateModelFiles(config.Models)
	generateControllerFiles(config.Models, config.Routes)
	generateRouteFile(config.Routes)
	generateMainFile()
}

func generateModelFiles(models map[string]map[string]string) {
	fmt.Println()
	fmt.Println("+---------------------------+")
	fmt.Println("| Generating model files... |")
	fmt.Println("+---------------------------+")
	fmt.Println()
	for modelName, fieldMap := range models {
		fmt.Println()
		f := os.Stdout

		// fmt.Println("Model Name:", modelName)
		// fmt.Println(fieldMap)

		fields := make([]Field, 0)
		for fieldName, fieldType := range fieldMap {
			fields = append(fields, Field{
				Name: fieldName,
				Type: fieldType,
			})
		}

		tmpl := template.Must(template.New("model").Funcs(template.FuncMap{
			"ToUpper": strings.ToUpper,
			"ToLower": strings.ToLower,
		}).Parse(modelTemplate))
		tmpl.Execute(f, struct {
			ModelName string
			Fields    []Field
		}{
			ModelName: modelName,
			Fields:    fields,
		})
	}
}

func generateControllerFiles(models map[string]map[string]string, routes []Route) {
	fmt.Println()
	fmt.Println("+--------------------------------+")
	fmt.Println("| Generating controller files... |")
	fmt.Println("+--------------------------------+")

	for modelName := range models {
		fmt.Println()
		f := os.Stdout

		// Filter routes for the current model
		var modelRoutes []Route
		for _, route := range routes {
			if route.Model == modelName {
				modelRoutes = append(modelRoutes, route)
			}
		}

		tmpl := template.Must(template.New("controller").Funcs(template.FuncMap{
			"ToLower": strings.ToLower,
		}).Parse(controllerTemplate))
		tmpl.Execute(f, struct {
			ModelName string
			Routes    []Route
		}{
			ModelName: modelName,
			Routes:    modelRoutes,
		})
	}
}

func generateRouteFile(routes []Route) {
	fmt.Println()
	fmt.Println("+--------------------------+")
	fmt.Println("| Generating route file... |")
	fmt.Println("+--------------------------+")
	fmt.Println()
	f := os.Stdout

	tmpl := template.Must(template.New("routes").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).Parse(routesTemplate))
	tmpl.Execute(f, struct {
		Routes []Route
	}{
		Routes: routes,
	})
}

func generateMainFile() {
	fmt.Println()
	fmt.Println("+-------------------------+")
	fmt.Println("| Generating main file... |")
	fmt.Println("+-------------------------+")
	fmt.Println()
	f := os.Stdout

	tmpl := template.Must(template.New("main").Parse(mainTemplate))
	tmpl.Execute(f, nil)
}

// Templates
var modelTemplate = `package models

import "gorm.io/gorm"

type {{ .ModelName }} struct {
	gorm.Model
{{- range .Fields }}
	{{ .Name }} {{ .Type }} ` + "`json:\"{{ .Name | ToLower }}\"`" + `
{{- end }}
}
`

var controllerTemplate = `package controllers

import (
	"net/http"
	"my-crud-api/config"
	"my-crud-api/models"
	"github.com/gin-gonic/gin"
)

{{- range .Routes }}
func {{ .Controller }}(c *gin.Context) {
	var {{ .Model | ToLower }} models.{{ .Model }}
	{{ if eq .Method "POST" }}if err := c.ShouldBindJSON(&{{ .Model | ToLower }}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&{{ .Model | ToLower }})
	c.JSON(http.StatusOK, &{{ .Model | ToLower }}){{ end }}
	{{ if eq .Method "GET" }}{{ if eq (index .Path 7) '{' }}
	if err := config.DB.First(&{{ .Model | ToLower }}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, &{{ .Model | ToLower }}){{ else }}
	var {{ .Model | ToLower }}s []models.{{ .Model }}
	config.DB.Find(&{{ .Model | ToLower }}s)
	c.JSON(http.StatusOK, &{{ .Model | ToLower }}s){{ end }}{{ end }}
	{{ if eq .Method "PUT" }}if err := config.DB.First(&{{ .Model | ToLower }}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	if err := c.ShouldBindJSON(&{{ .Model | ToLower }}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&{{ .Model | ToLower }})
	c.JSON(http.StatusOK, &{{ .Model | ToLower }}){{ end }}
	{{ if eq .Method "DELETE" }}if err := config.DB.Delete(&{{ .Model | ToLower }}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "{{ .Model }} deleted"}){{ end }}
}
{{- end }}
`

var routesTemplate = `package routes

import (
	"my-crud-api/controllers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
{{- range .Routes }}
	router.{{ .Method }}("{{ .Path }}", controllers.{{ .Controller }})
{{- end }}
}
`

var mainTemplate = `package main

import (
	"my-crud-api/config"
	"my-crud-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	router := gin.Default()
	routes.InitializeRoutes(router)

	router.Run(":8080")
}
`
