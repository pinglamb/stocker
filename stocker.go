package main

import "fmt"
import "os"

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
  return nil
}

func commandUp() (err error) {
  return nil
}

func ensureDockerComposeYamlExists() {
  if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
    fmt.Println("Creating docker-compose.yml")
    f, err := os.Create("docker-compose.yml")
    check(err)
    _, err = f.WriteString("version: \"2\"\n\nservices:\n")
    check(err)

    f.Sync()
  } else {
    fmt.Println("Found docker-compose.yml")
  }
}

func check(e error) {
  if e != nil {
    panic(e)
  }
}
