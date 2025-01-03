package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Define flags
	layerName := flag.String("name", "", "Name of the Lambda layer")
	runtimes := flag.String("runtimes", "python3.10", "Comma-separated list of compatible runtimes")
	requirementsFile := flag.String("requirements", "requirements.txt", "Path to the requirements.txt file")
	deploy := flag.Bool("deploy", false, "Set this flag to deploy the Lambda layer")
	flag.Parse()

	// Validate flags
	if *layerName == "" {
		log.Fatalf("--name must be specified")
	}

	// Build the Lambda layer
	err := buildLambdaLayer(*layerName, *runtimes, *requirementsFile)
	if err != nil {
		log.Fatalf("Error building Lambda layer: %v", err)
	}

	// Deploy the Lambda layer if requested
	if *deploy {
		err = deployLayer(*layerName, *runtimes)
		if err != nil {
			log.Fatalf("Error deploying Lambda layer: %v", err)
		}
	}
}

func buildLambdaLayer(layerName string, compatibleRuntimes string, requirementsFile string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	absRequirementsFile := filepath.Join(currentDir, requirementsFile)

	fmt.Println("Building Lambda layer with the following details:")
	fmt.Printf("Layer Name: %s\n", layerName)
	fmt.Printf("Compatible Runtimes: %s\n", compatibleRuntimes)
	fmt.Printf("Requirements File: %s\n", absRequirementsFile)
	fmt.Println(strings.Repeat("#", 80))

	if !confirmPrompt("Should I proceed building this Layer...is Docker running? (yes/no): ") {
		return errors.New("operation canceled by user")
	}

	dockerImage := "public.ecr.aws/lambda/python:3.10"
	containerName := "amzn"

	executeCommand("docker pull " + dockerImage)

	existingContainers := executeCommand("docker ps -a --format '{{.Names}}'")
	if strings.Contains(existingContainers, containerName) {
		fmt.Printf("Container %s already exists. Removing it...\n", containerName)
		executeCommand("docker rm -f " + containerName)
	}

	executeCommand(fmt.Sprintf("docker run --name %s -d -t --rm %s /bin/bash", containerName, dockerImage))

	if _, err := os.Stat(absRequirementsFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("requirements file '%s' not found", absRequirementsFile)
	}

	fmt.Printf("Copying requirements file from %s to container...\n", absRequirementsFile)
	executeCommand(fmt.Sprintf("docker cp %s %s:/root/requirements.txt", absRequirementsFile, containerName))

	commands := []string{
		"yum update -y",
		"yum install -y zip",
		"python3.10 -m ensurepip",
		"python3.10 -m pip install --upgrade pip",
		"python3.10 -m pip install -r /root/requirements.txt -t /root/package/python",
		"find /root/package/ -name '*.pyc' -delete",
		"rm -rf /root/package/urllib3* /root/package/six* /root/package/botocore* /root/package/idna* /root/package/tomli* /root/package/jmespath* /root/package/charset_normalizer*",
		"find /root/package/ -name '*.dist-info' -exec rm -rf {} +",
		"find /root/package/ -type d -name '__pycache__' -exec rm -r {} +",
		fmt.Sprintf("cd /root/package && zip -r9qv ../%s.zip *", layerName),
	}

	for _, cmd := range commands {
		fmt.Printf("Executing inside container: %s\n", cmd)
		executeCommand(fmt.Sprintf("docker exec %s sh -c \"%s\"", containerName, cmd))
	}

	outputDir := currentDir
	fmt.Printf("Copying the resulting ZIP file to %s/%s.zip...\n", outputDir, layerName)
	executeCommand(fmt.Sprintf("docker cp %s:/root/%s.zip %s/%s.zip", containerName, layerName, outputDir, layerName))

	fmt.Println("Cleaning up Docker container and temporary files...")
	executeCommand("docker rm -f " + containerName)

	return nil
}

func executeCommand(command string) string {
	fmt.Printf("Running: %s\n", command)
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		exitWithError(fmt.Sprintf("Command failed: %s", command))
	}
	return out.String()
}

func confirmPrompt(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return response == "yes"
}

func exitWithError(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func deployLayer(layerName string, compatibleRuntimes string) error {
	layerDescription := fmt.Sprintf("%s for Lambda functions", strings.Title(layerName))
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	outputDir := currentDir
	zipFilePath := filepath.Join(outputDir, fmt.Sprintf("%s.zip", layerName))

	if !confirmPrompt("Should I Publish the layer to AWS Lambda? (yes/no): ") {
		fmt.Println("Skipping publishing the layer to AWS Lambda.")
		return nil
	}

	fmt.Printf("Publishing the layer to AWS Lambda with name '%s'...\n", layerName)
	executeCommand(fmt.Sprintf(
		"aws lambda publish-layer-version --layer-name %s --description '%s' --zip-file fileb://%s --compatible-runtimes %s",
		layerName, layerDescription, zipFilePath, compatibleRuntimes,
	))

	return nil
}
