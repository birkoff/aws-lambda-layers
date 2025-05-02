package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runtimeToDockerImage(runtime string) string {
	runtime = strings.TrimPrefix(runtime, "python")
	return fmt.Sprintf("public.ecr.aws/lambda/python:%s", runtime)
}

func main() {
	layerName := flag.String("name", "", "Name of the Lambda layer")
	runtime := flag.String("runtime", "python3.10", "Python runtime version")
	requirementsFile := flag.String("requirements", "requirements.txt", "Path to requirements.txt")
	deploy := flag.Bool("deploy", false, "Deploy the Lambda layer")
	flag.Parse()

	if *layerName == "" {
		log.Fatal("--name must be specified")
	}

	if err := buildLambdaLayer(*layerName, *runtime, *requirementsFile); err != nil {
		log.Fatalf("Error building layer: %v", err)
	}

	if *deploy {
		if err := deployLayer(*layerName, *runtime); err != nil {
			log.Fatalf("Error deploying layer: %v", err)
		}
	}
}

func buildLambdaLayer(layerName, runtime, requirementsFile string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	absRequirementsFile := filepath.Join(currentDir, requirementsFile)
	if _, err := os.Stat(absRequirementsFile); err != nil {
		return fmt.Errorf("requirements file not found: %w", err)
	}

	// if !confirmPrompt("Proceed with building the layer? (yes/no): ") {
	// 	return errors.New("operation canceled")
	// }

	dockerImage := runtimeToDockerImage(runtime)
	containerName := "lambda-builder"

	// Cleanup any existing container
	if output := executeCommand("docker ps -a --format '{{.Names}}'"); strings.Contains(output, containerName) {
		executeCommand(fmt.Sprintf("docker rm -f %s", containerName))
	}

	// Start container with proper platform and bind mounts
	executeCommand(fmt.Sprintf(
		"docker run --platform linux/amd64 --name %s -d -t --rm -v %s:/host %s /bin/bash",
		containerName, currentDir, dockerImage,
	))

	// Install system dependencies and Python packages
	commands := []string{
		// Base system updates
		"yum update -y",
		"yum install -y gcc python3-devel",

		// PostgreSQL dependencies
		"yum install -y gcc python3-devel postgresql-devel postgresql-libs",
		// SQL Server dependencies
		"ACCEPT_EULA=Y yum install -y msodbcsql17 freetds freetds-devel unixodbc unixODBC-devel",
		
		// Python toolchain
		fmt.Sprintf("%s -m pip install --upgrade pip", runtime),

		// Package installation
		fmt.Sprintf("%s -m pip install -r /host/%s -t /root/package/python", runtime, requirementsFile),


		// Library bundling
		"cp /usr/lib64/libpq.so* /root/package/python/",
		"cp /usr/lib64/libodbc.so* /root/package/python/",
		
		// Cleanup
		"find /root/package/ -type d -name '__pycache__' -exec rm -rf {} +",
		fmt.Sprintf("cd /root/package && zip -r9qv ../%s.zip *", layerName),

	}

	for _, cmd := range commands {
		executeCommand(fmt.Sprintf("docker exec %s sh -c %q", containerName, cmd))
	}


	outputDir := currentDir
	fmt.Printf("Copying the resulting ZIP file to %s/%s.zip...\n", outputDir, layerName)
	executeCommand(fmt.Sprintf("docker cp %s:/root/%s.zip %s/%s.zip", containerName, layerName, outputDir, layerName))

	fmt.Println("Cleaning up Docker container and temporary files...")
	executeCommand(fmt.Sprintf("docker rm -f %s", containerName))


	return nil
}

func executeCommand(command string) string {
	fmt.Printf("Running: %s\n", command)
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Command failed: %v\n", err)
		os.Exit(1)
	}
	return out.String()
}

func confirmPrompt(prompt string) bool {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ToLower(scanner.Text()) == "yes"
}

func deployLayer(layerName, runtime string) error {
	if !confirmPrompt("Publish layer to AWS Lambda? (yes/no): ") {
		return nil
	}

	description := fmt.Sprintf("%s layer for %s", layerName, runtime)
	executeCommand(fmt.Sprintf(
		"aws lambda publish-layer-version --layer-name %s --description %q --zip-file fileb://%s.zip --compatible-runtimes %s",
		layerName, description, layerName, runtime,
	))

	return nil
}