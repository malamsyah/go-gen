# Go API Generator
This Go project dynamically generates API models, controllers, and routes based on a specified YAML configuration file. It streamlines the setup of a CRUD API by automating the boilerplate code creation process.

## Features
Model Generation: Automatically generates model files using GORM based on the structure defined in the YAML file.
Controller Generation: Creates controller files with predefined CRUD operations.
Route Generation: Sets up routing using the Gin framework according to the definitions in the YAML file.
Automatic File Creation: Based on the YAML configuration, the necessary Go files for models, controllers, and routes are generated on the fly.

## Prerequisites
Before running this project, ensure you have the following installed:

- Go (version 1.13 or higher)
- Gin-Gonic Framework
- GORM
- YAML v3 package for Go

## You can install the necessary Go packages using:

```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get gopkg.in/yaml.v3
```

## Configuration
Define your API models and routes in a YAML file. Here's a sample structure for the config.yaml:

```yaml
Copy code
models:
  User:
    ID: "int"
    Username: "string"
    Email: "string"
routes:
  - path: "/users"
    method: "GET"
    controller: "GetUsers"
    model: "User"
  - path: "/users/{id}"
    method: "GET"
    controller: "GetUser"
    model: "User"
  - path: "/users"
    method: "POST"
    controller: "CreateUser"
    model: "User"
```
## Usage
Prepare your YAML configuration file and name it as desired, e.g., config.yaml.
Run the program by passing the YAML file as an argument:

```bash
go run main.go config.yaml
```

The program reads the configuration file and generates Go source files for models, controllers, and the router based on the definitions provided.

## Generated Files
The program generates the following types of files:

- Model files (models/*.go): These files contain the GORM model definitions.
- Controller files (controllers/*.go): These files implement the CRUD operations as defined in the routes.
- Route file (routes/routes.go): This file sets up the routing using Gin.

## Contributing
Contributions are welcome! Please feel free to submit pull requests, create issues for bugs, or suggest improvements.

