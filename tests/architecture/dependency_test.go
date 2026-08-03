package architecture_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const module = "github.com/DrizzDev/platform"

type edge struct {
	owner      string
	dependency string
}

type location string

type layer string

const (
	domain         layer = "domain"
	application    layer = "application"
	adapter        layer = "adapter"
	infrastructure layer = "infrastructure"
	transport      layer = "transport"
	platform       layer = "platform"
)

func TestDependencies(test *testing.T) {
	test.Parallel()

	repository := repository{root: filepath.Join("..", ".."), test: test}
	var violations []string
	for owner, dependencies := range repository.dependencies() {
		for _, dependency := range dependencies {
			edge := edge{owner: owner, dependency: dependency}
			if violation := repository.dependency(edge); violation != "" {
				violations = append(violations, violation)
			}
		}
	}
	if len(violations) != 0 {
		test.Fatalf("architecture violations:\n%s", strings.Join(violations, "\n"))
	}
}

func (repository repository) dependencies() map[string][]string {
	repository.test.Helper()
	command := exec.CommandContext(context.Background(), "go", "list", "-json", "./internal/...")
	command.Dir = repository.root
	command.Env = append(os.Environ(), "GOCACHE="+filepath.Join(repository.test.TempDir(), "cache"))
	output, failure := command.Output()
	if failure != nil {
		var exit *exec.ExitError
		if errors.As(failure, &exit) {
			repository.test.Fatalf("go list: %v\n%s", failure, exit.Stderr)
		}
		repository.test.Fatal(failure)
	}

	result := make(map[string][]string)
	decoder := json.NewDecoder(bytes.NewReader(output))
	for decoder.More() {
		var value struct {
			Path    string   `json:"ImportPath"`
			Imports []string `json:"Imports"`
		}
		if failure := decoder.Decode(&value); failure != nil {
			repository.test.Fatal(failure)
		}
		result[value.Path] = value.Imports
	}
	return result
}

func (repository repository) dependency(edge edge) string {
	external := !strings.HasPrefix(edge.dependency, module+"/")
	if external {
		return repository.external(edge)
	}
	return repository.internal(edge)
}

func (repository repository) external(edge edge) string {
	if repository.core(edge.owner) && strings.Contains(edge.dependency, ".") {
		return edge.owner + " imports " + edge.dependency + ": core layers may use only the standard library"
	}
	return ""
}

func (repository repository) internal(edge edge) string {
	owner := location(edge.owner)
	switch {
	case owner.has(domain):
		return repository.domain(edge)
	case owner.has(application):
		return repository.application(edge)
	case owner.has(platform):
		return repository.platform(edge)
	case owner.has(transport):
		return repository.transport(edge)
	default:
		return ""
	}
}

func (repository repository) domain(edge edge) string {
	if !location(edge.dependency).has(domain) {
		return edge.owner + " imports " + edge.dependency + ": domain dependencies must point inward"
	}
	return ""
}

func (repository repository) application(edge edge) string {
	if repository.outer(edge.dependency) {
		return edge.owner + " imports " + edge.dependency + ": application dependencies must point inward"
	}
	return ""
}

func (repository repository) platform(edge edge) string {
	if repository.product(edge.dependency) {
		return edge.owner + " imports " + edge.dependency + ": platform must remain independent"
	}
	return ""
}

func (repository repository) transport(edge edge) string {
	dependency := location(edge.dependency)
	switch {
	case dependency.has(adapter):
		return edge.owner + " imports " + edge.dependency + ": transports may not compose adapters"
	case dependency.has(transport) && edge.owner != edge.dependency:
		return edge.owner + " imports " + edge.dependency + ": transports must communicate through owned contracts"
	default:
		return ""
	}
}

func (repository repository) core(path string) bool {
	location := location(path)
	return location.has(domain) || location.has(application)
}

func (repository repository) outer(path string) bool {
	location := location(path)
	return location.has(adapter) ||
		location.has(infrastructure) ||
		location.has(transport) ||
		location.has(platform)
}

func (repository repository) product(path string) bool {
	location := location(path)
	return location.has(domain) ||
		location.has(application) ||
		location.has(adapter) ||
		location.has(infrastructure) ||
		location.has(transport)
}

func (location location) has(layer layer) bool {
	path := string(location)
	name := string(layer)
	if layer == platform {
		return strings.Contains(path, "/internal/platform/") || strings.HasSuffix(path, "/internal/platform")
	}
	return strings.Contains(path, "/"+name+"/") || strings.HasSuffix(path, "/"+name)
}
