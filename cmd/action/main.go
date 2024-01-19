package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// all of the env params and inputs are defined here
var (
	//required
	path    = flag.String("path", ".", "Path to main Go main package.")
	name    = flag.String("name", "", "Override action name, the default name is the package name.")
	desc    = flag.String("desc", "", "Override action description, the default description is the package synopsis.")
	image   = flag.String("image", "golang:1.14.2-alpine3.11", "Override Docker image to run the action with (See https://hub.docker.com/_/golang?tab=tags).")
	install = flag.String("install", "", "Comma separated list of requirements to 'apk add'.")
	icon    = flag.String("icon", "", "Set branding icon. (See options at https://feathericons.com).")
	color   = flag.String("color", "", "Set branding color. (white, yellow, blue, green, orange, red, purple or gray-dark).")
	domain  = flag.String("domain", "", "Company domain associated with registered corporate accound in SwaggerHub (See account at https://swagger.io)")

	//description Email for commit message.
	//default posener@gmail.com
	email = os.Getenv("email")
	//description Github token for PR comments. Optional.
	githubToken      = os.Getenv("GITHUB_TOKEN")
	gitBranch        = os.Getenv("GIT_BRANCH")
	swaggerhubApiKey = os.Getenv("SWAGGERHUB_API_KEY")
	swaggerFile      = os.Getenv("SWAGGER_FILE")
)

func prepareVersion() string {
	version := strings.ReplaceAll(gitBranch, "/", "_")
	version = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '.', r == '-':
			return r
		}
		return -1
	}, version)
	return version
}

func publish(openApi map[string]interface{}, key, domain, apiName, apiVersion string) error {
	// oas 3.0 and private only for now
	// just to make it work and test
	url := fmt.Sprintf("https://api.swaggerhub.com/apis/%s/%s?version=%s&oas=3.0&isPrivate=true&force=true", domain, apiName, apiVersion)

	payload, err := json.Marshal(openApi)
	if err != nil {
		return fmt.Errorf("failed to parse swagger.json with OpenAPI standard: %s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for given swagger.json file: %s", err)
	}
	req.Header.Set("Authorization", key)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to submit request over SwaggerHub: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("failed to publish swagger.json to SwaggerHub %s", resBody)
	}

	fmt.Printf("Publish %s of %s is successfull", apiName, apiVersion)
	return nil
}

func main() {
	flag.Parse()

	if swaggerhubApiKey == "" {
		fmt.Println("SWAGGERHUB_API_KEY is not defined")
		return
	}

	version := prepareVersion()
	if version == "" {
		fmt.Println("Empty version, skipping upload")
		return
	}

	fileBytes, err := ioutil.ReadFile(swaggerFile)
	if err != nil {
		fmt.Printf("Error reading swagger.json: %s\n", err)
		return
	}

	var openAPI map[string]interface{}
	err = json.Unmarshal(fileBytes, &openAPI)
	if err != nil {
		fmt.Printf("Desearilization error for given swagger.json %s\n", err)
		return
	}

	err = publish(openAPI, swaggerhubApiKey, *domain, *name, version)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}
