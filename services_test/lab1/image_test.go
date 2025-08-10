package services_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestImageExists(t *testing.T) {
	user := os.Getenv("REPO_USER")
	if user == "" {
		t.Fatal("REPO_USER environment variable not set")
	}

	tag := os.Getenv("TAG")
	if tag == "" {
		// set the tag to lab1 by default
		tag = "lab1"
	}

	imageName := user + "/restaurant_microservice:" + tag

	// Check if the Docker image exists
	cmd := exec.Command("docker", "images", "-q", imageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to check image existence: %v\nOutput: %s", err, output)
	}

	if len(output) == 0 {
		t.Errorf("Docker image '%s' does not exist", imageName)
	} else {
		fmt.Printf("Docker image '%s' exists\n", imageName)
	}
}
