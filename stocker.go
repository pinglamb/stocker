package main

import (
  "fmt"
  "os"
  "strings"
  "io/ioutil"
)

func main() {
  command, args, success := extractArgs()
  if !success { return }

  if command == "add" {
    commandAdd(args)
  } else if command == "up" {
    commandUp()
  }
}

func extractArgs() (command string, args []string, success bool) {
  if len(os.Args) < 2 {
    fmt.Println("Usage: stocker [add/up]")
    success = false
    return
  }

  command = os.Args[1]
  args = os.Args[2:]

  if command != "add" && command != "up" {
    fmt.Println("Usage: stocker [add/up]")
    success = false
    return
  }

  if command == "add" && len(args) == 0 {
    fmt.Println("Usage: stocker add [service name] (e.g. postgres)")
    success = false
    return
  }

  success = true
  return
}

func commandAdd(args []string) (err error) {
  ensureDockerComposeYamlExists()

  service := args[0]
  serviceConfig := fmt.Sprintf("  %s:\n    image: %s:latest\n", service, service)

  b, err := ioutil.ReadFile("docker-compose.yml")
  check(err)

  content := string(b)
  newContent := strings.Replace(content, "services:\n", fmt.Sprintf("services:\n%s", serviceConfig), 1)

  err = ioutil.WriteFile("docker-compose.yml", []byte(newContent), 0644)
  check(err)

  fmt.Println("Added %s service to your docker-compose.yml", service)

  return nil
}

func commandUp() (err error) {
  return nil
}

func ensureDockerComposeYamlExists() {
  if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
    err := ioutil.WriteFile("docker-compose.yml", []byte("version: \"2\"\n\nservices:\n"), 0644)
    check(err)
  } else {
    fmt.Println("Found docker-compose.yml")
  }
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}
