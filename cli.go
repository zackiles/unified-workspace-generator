package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var cliDir string
var workspaceDir string
var globalConfigDir string
var workspaceFile string

func main() {
	// Determine the CLI's directory (the directory where this executable is located)
	cliDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	// Set default paths for config and workspace inside the CLI directory
	globalConfigDir = filepath.Join(cliDir, "global-config")
	workspaceFile = filepath.Join(cliDir, "projects.code-workspace")

	app := &cli.App{
		Name:  "create-project",
		Usage: "Create a new Go or JavaScript project in a specified workspace directory",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Project name",
			},
			&cli.StringFlag{
				Name:    "language",
				Aliases: []string{"l"},
				Usage:   "Project language (js/go)",
			},
			&cli.StringFlag{
				Name:    "workspace",
				Aliases: []string{"w"},
				Usage:   "Root workspace directory for projects",
			},
		},
		Action: func(c *cli.Context) error {
			projectName := c.String("name")
			lang := c.String("language")
			workspaceDir = c.String("workspace")

			// Prompt for the workspace directory if not provided
			if workspaceDir == "" {
				fmt.Print("Enter the root workspace directory (where all project folders will be created): ")
				fmt.Scanln(&workspaceDir)
			}

			// Update globalConfigDir and workspaceFile paths based on workspaceDir
			globalConfigDir = filepath.Join(workspaceDir, "global-config")
			workspaceFile = filepath.Join(workspaceDir, "projects.code-workspace")

			// Prompt for project name if not provided
			if projectName == "" {
				fmt.Print("Enter project name: ")
				fmt.Scanln(&projectName)
			}

			// Prompt for project language if not provided
			if lang == "" {
				fmt.Print("Enter language (js/go): ")
				fmt.Scanln(&lang)
			}

			// Copy the CLI's configuration files to the workspace
			copyConfigToWorkspace()

			// Create the project
			createProject(projectName, lang)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// copyConfigToWorkspace copies the CLI's config files to the user's workspace
func copyConfigToWorkspace() {
	// Ensure the global-config directory exists in the workspace
	err := os.MkdirAll(globalConfigDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating global config directory:", err)
		return
	}

	// Copy the files from the CLI's global-config folder to the workspace
	err = copyFile(filepath.Join(cliDir, "global-config", ".eslintrc.json"), filepath.Join(globalConfigDir, ".eslintrc.json"))
	if err != nil {
		fmt.Println("Error copying ESLint config:", err)
	}
	err = copyFile(filepath.Join(cliDir, "global-config", ".prettierrc"), filepath.Join(globalConfigDir, ".prettierrc"))
	if err != nil {
		fmt.Println("Error copying Prettier config:", err)
	}
	err = copyFile(filepath.Join(cliDir, "global-config", ".golangci.yml"), filepath.Join(globalConfigDir, ".golangci.yml"))
	if err != nil {
		fmt.Println("Error copying GolangCI config:", err)
	}
}

// createProject sets up the project with the specified language (js or go)
func createProject(name string, lang string) {
	projectPath := filepath.Join(workspaceDir, name)

	// Create the project directory
	os.MkdirAll(projectPath, os.ModePerm)
	fmt.Printf("Created project folder: %s\n", projectPath)

	// Initialize Git repository
	gitInit(projectPath)

	// Create language-specific setup
	if lang == "js" {
		setupJavaScriptProject(projectPath)
	} else if lang == "go" {
		setupGolangProject(projectPath)
	} else {
		fmt.Println("Unsupported language. Please choose 'js' or 'go'.")
		return
	}

	// Add the project to the workspace
	updateWorkspace(name)
}

// setupJavaScriptProject handles JS-specific setup
func setupJavaScriptProject(path string) {
	// Run npm init and install eslint + prettier
	runCommand("npm", []string{"init", "-y"}, path)
	runCommand("npm", []string{"install", "--save-dev", "eslint", "prettier"}, path)

	// Copy ESLint and Prettier config files
	copyFile(filepath.Join(globalConfigDir, ".eslintrc.json"), filepath.Join(path, ".eslintrc.json"))
	copyFile(filepath.Join(globalConfigDir, ".prettierrc"), filepath.Join(path, ".prettierrc"))

	fmt.Println("JavaScript project setup complete.")
}

// setupGolangProject handles Go-specific setup
func setupGolangProject(path string) {
	// Initialize go mod and get golangci-lint
	runCommand("go", []string{"mod", "init", "github.com/yourusername/" + filepath.Base(path)}, path)
	runCommand("go", []string{"get", "-u", "github.com/golangci/golangci-lint/cmd/golangci-lint"}, path)

	// Copy GolangCI config
	copyFile(filepath.Join(globalConfigDir, ".golangci.yml"), filepath.Join(path, ".golangci.yml"))

	fmt.Println("Golang project setup complete.")
}

// gitInit initializes a Git repository in the project folder
func gitInit(path string) {
	runCommand("git", []string{"init"}, path)
	fmt.Println("Initialized Git repository.")
}

// updateWorkspace updates the workspace file with the new project path
func updateWorkspace(projectName string) {
	projectPath := filepath.Join(workspaceDir, projectName)

	if _, err := os.Stat(workspaceFile); os.IsNotExist(err) {
		// Create a new workspace file if it doesn't exist
		file, _ := os.Create(workspaceFile)
		file.WriteString("{\"folders\": [], \"settings\": {}}\n")
		file.Close()
	}

	// Read workspace content
	content, _ := os.ReadFile(workspaceFile)
	if !strings.Contains(string(content), projectPath) {
		// Add the project path to the workspace
		newContent := strings.Replace(string(content), "\"folders\": [", "\"folders\": [{ \"path\": \""+projectPath+"\" },", 1)
		os.WriteFile(workspaceFile, []byte(newContent), os.ModePerm)
		fmt.Printf("Added %s to the workspace.\n", projectName)
	}
}

// runCommand runs shell commands in the specified directory
func runCommand(command string, args []string, dir string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	}
}

// copyFile copies files from source to destination and returns an error if it fails
func copyFile(src, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error reading source file: %w", err)
	}
	err = os.WriteFile(dest, input, 0644)
	if err != nil {
		return fmt.Errorf("error writing destination file: %w", err)
	}
	fmt.Printf("Copied %s to %s\n", filepath.Base(src), dest)
	return nil
}
