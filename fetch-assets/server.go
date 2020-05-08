package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type CreateReactAppAssets struct {
	scripts []string
	links   []string
}

func main() {
	fmt.Println("hello world")
	cra := parse()
	fmt.Println(cra)
}

func parse() CreateReactAppAssets {
	var body string
	var err error
	switch os.Getenv("APP_ENV") {
	case "production":
		body, err = fetchFromBuild()
	default:
		// case "development":
		body, err = fetchFromReactServer()
		// 	panic(fmt.Sprintf("Unsupported environment: %s", os.Getenv("APP_ENV")))
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(body)

	scriptRegex := regexp.MustCompile("<script.*</script>")
	linksRegex := regexp.MustCompile("<link\\s+href=\"[^\"]*\"\\s+rel=\"stylesheet\">")
	scripts := scriptRegex.FindAllString(body, -1)
	links := linksRegex.FindAllString(body, -1)

	return CreateReactAppAssets{scripts, links}
}

func fetchFromReactServer() (string, error) {
	url := "http://localhost:3002/"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Cannot fetch %s", url))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}

func fetchFromBuild() (string, error) {
	env := "REACT_DIR"
	react_dir, exists := os.LookupEnv(env)
	if !exists {
		return "", errors.New(fmt.Sprintf("Cannot fetch environment variable: %s", env))
	}

	filename := react_dir + "/build/index.html"
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return "", errors.New(fmt.Sprintf("Path does not exist: %s", filename))
	}

	if info.IsDir() {
		return "", errors.New(fmt.Sprintf("Path is a directory: %s", filename))
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
