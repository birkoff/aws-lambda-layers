package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "lambda-layer-tool",
		Short: "A CLI tool to build and publish AWS Lambda layers",
	}

	var buildCmd = &cobra.Command{
		Use:   "build <layer-name> <compatible-runtimes>",
		Short: "Build a Lambda layer",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			layerName := args[0]
			compatibleRuntimes := args[1]
			if err := buildLambdaLayer(layerName, compatibleRuntimes); err != nil {
				exitWithError(err.Error())
			}
		},
	}

	rootCmd.AddCommand(buildCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildLambdaLayer(layerName, compatibleRuntimes string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	requirementsFile := filepath.Join(currentDir, fmt.Sprintf("%s-requirements.txt", layerName))
	packageName := layerName
	layerDescription := fmt.Sprintf("%s for Lambda functions", strings.Title(layerName))

	fmt.Println("Building Lambda layer with the following details:")
	fmt.Printf("Layer Name: %s\n", layerName)
	fmt.Printf("Compatible Runtimes: %s\n", compatibleRuntimes)
	fmt.Printf("Requirements File: %s\n", requirementsFile)
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

	if _, err := os.Stat(requirementsFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("requirements file '%s' not found", requirementsFile)
	}

	absPath, _ := filepath.Abs(requirementsFile)
	fmt.Printf("Copying requirements file from %s to container...\n", absPath)
	executeCommand(fmt.Sprintf("docker cp %s %s:/root/requirements.txt", absPath, containerName))

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
		fmt.Sprintf("cd /root/package && zip -r9qv ../%s.zip *", packageName),
	}

	for _, cmd := range commands {
		fmt.Printf("Executing inside container: %s\n", cmd)
		executeCommand(fmt.Sprintf("docker exec %s sh -c \"%s\"", containerName, cmd))
	}

	outputDir := currentDir
	fmt.Printf("Copying the resulting ZIP file to %s/%s.zip...\n", outputDir, packageName)
	executeCommand(fmt.Sprintf("docker cp %s:/root/%s.zip %s/%s.zip", containerName, packageName, outputDir, packageName))

	if confirmPrompt("Should I Publish the layer to AWS Lambda? (yes/no): ") {
		fmt.Printf("Publishing the layer to AWS Lambda with name '%s'...\n", layerName)
		executeCommand(fmt.Sprintf(
			"aws lambda publish-layer-version --layer-name %s --description '%s' --zip-file fileb://%s/%s.zip --compatible-runtimes %s",
			layerName, layerDescription, outputDir, packageName, compatibleRuntimes,
		))
	} else {
		fmt.Println("Skipping publishing the layer to AWS Lambda.")
	}

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
