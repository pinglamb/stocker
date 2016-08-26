package main

import (
  "fmt"
  "os"
  "os/exec"
  "bytes"
  "strings"
  "io/ioutil"
  "encoding/json"
)

type DockerInspect struct {
  Config DockerImageConfig `json:"Config"`
}

type DockerImageConfig struct {
  ExposedPorts map[string]*json.RawMessage `json:"ExposedPorts"`
}

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
    fmt.Println("Usage: stocker add SERVICE_NAME [-f] (e.g. stocker add postgres:9.5)")
    success = false
    return
  }

  success = true
  return
}

func commandAdd(args []string) (err error) {
  ensureDockerComposeYamlExists()

  service := args[0]
  version := "latest"
  if strings.Contains(service, ":") {
    serviceAndVersion := strings.Split(service, ":")
    service = serviceAndVersion[0]
    version = serviceAndVersion[1]
  }
  image := fmt.Sprintf("%s:%s", service, version)

  force := false
  if len(args) > 1 && args[1] == "-f" {
    force = true
  }

  b, err := ioutil.ReadFile("docker-compose.yml")
  check(err)

  content := string(b)

  if !force {
    if strings.Contains(content, fmt.Sprintf("%s:", service)) {
      fmt.Println("Found", service, "service in your docker-compose.yml, -f to add it anyway")
      return nil
    }
  }

  pullCmd := exec.Command("docker", "pull", image)
  pullCmd.Stdout = os.Stdout
  pullCmd.Stderr = os.Stderr
  err = pullCmd.Run()
  check(err)

  inspectCmd := exec.Command("docker", "inspect", image)
  var output bytes.Buffer
  inspectCmd.Stdout = &output
  err = inspectCmd.Run()
  check(err)
  var j []DockerInspect
  err = json.Unmarshal(output.Bytes(), &j)
  check(err)
  var ports bytes.Buffer
  if len(j[0].Config.ExposedPorts) > 0 {
    ports.WriteString("    ports:\n")
    for k := range j[0].Config.ExposedPorts {
      port := strings.Split(k, "/")[0]
      ports.WriteString(fmt.Sprintf("      - %s:%s\n", port, port))
    }
  }

  serviceConfig := fmt.Sprintf("  %s:\n    image: \"%s\"\n%s", service, image, ports.String())
  newContent := strings.Replace(content, "services:\n", fmt.Sprintf("services:\n%s", serviceConfig), 1)

  err = ioutil.WriteFile("docker-compose.yml", []byte(newContent), 0644)
  check(err)

  fmt.Println("Added", service, "service to your docker-compose.yml")

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
