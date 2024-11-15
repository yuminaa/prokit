package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	version = "1.0.0"

	reset     = "\033[0m"
	bold      = "\033[1m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	gray      = "\033[90m"
	dim       = "\033[2m"
	italic    = "\033[3m"
	underline = "\033[4m"
)

type ProjectConfig struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Author       string   `json:"author"`
	Files        []string `json:"files"`
	Dependencies []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"dependencies"`
	Scripts map[string]string `json:"scripts"`
}

type CLIFlags struct {
	language    string
	projectType string
	name        string
	output      string
	version     bool
}

func main() {
	flags := parseFlags()

	if flags.version {
		printVersion()
		return
	}

	if !isValidLanguage(flags.language) {
		fmt.Printf("%s%s✘ Error: Unsupported language: %s%s\n", bold, red, flags.language, reset)
		fmt.Printf("\n%sSupported languages:%s", cyan, reset)
		fmt.Printf("  Go, TS, C, CXX Python\n")
		os.Exit(1)
	}

	fmt.Printf("%s%s⚡ Loading configuration for %s project...%s\n", bold, blue, flags.language, reset)
	config, err := loadLanguageConfig(flags.language)
	if err != nil {
		fmt.Printf("%s%s✘ Error: %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	if err := createProject(flags, config); err != nil {
		fmt.Printf("%s%s✘ Error: %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	fmt.Printf("\n%s%s✓ Successfully created %s project: %s%s\n", bold, green, flags.language, flags.name, reset)
	fmt.Printf("%s  Location: %s%s\n", gray, filepath.Join(flags.output, flags.name), reset)
}

func printVersion() {
	fmt.Printf("%s%sproKit%s %sv%s%s\n", bold, magenta, reset, italic, version, reset)
	fmt.Printf("%sA minimal project scaffolding tool%s\n", gray, reset)
}

func parseFlags() CLIFlags {
	flags := CLIFlags{}

	flag.BoolVar(&flags.version, "version", false, "Show version information")
	flag.StringVar(&flags.language, "lang", "",
		"Programming language (go, ts, c, cxx, python)")
	flag.StringVar(&flags.projectType, "type", "app",
		"Project type (app or library)")
	flag.StringVar(&flags.name, "name", "",
		"Project name")
	flag.StringVar(&flags.output, "output", ".",
		"Output directory (default: current directory)")

	flag.Usage = func() {
		fmt.Printf("\n%s%sproKit%s - Project Generator%s\n", bold, magenta, reset, dim)
		fmt.Printf("%sVersion %s%s\n\n", dim, version, reset)

		fmt.Printf("%s%sUsage:%s\n", bold, yellow, reset)
		fmt.Printf("  prokit -lang=<language> -name=<project-name> [options]\n\n")

		fmt.Printf("%s%sExamples:%s\n", bold, yellow, reset)
		fmt.Printf("  %sprokit -lang=go -name=myproject%s\n", cyan, reset)
		fmt.Printf("  %sprokit -lang=ts -name=mylibrary -type=library%s\n", cyan, reset)
		fmt.Printf("  %sprokit -lang=python -name=myapp -output=./projects%s\n\n", cyan, reset)

		fmt.Printf("%s%sFlags:%s\n", bold, yellow, reset)
		fmt.Printf("  %s-lang%s string\n", green, reset)
		fmt.Printf("        Programming language (go, ts, c, cxx, python)\n\n")
		fmt.Printf("  %s-name%s string\n", green, reset)
		fmt.Printf("        Project name\n\n")
		fmt.Printf("  %s-type%s string\n", green, reset)
		fmt.Printf("        Project type (app or library) (default: \"app\")\n\n")
		fmt.Printf("  %s-output%s string\n", green, reset)
		fmt.Printf("        Output directory (default: current directory)\n\n")
		fmt.Printf("  %s-version%s\n", green, reset)
		fmt.Printf("        Show version information\n\n")

		fmt.Printf("%s%sLanguages:%s %sGo, TS, C, CXX, C#, Python%s\n\n",
			bold, yellow, reset, cyan, reset)
	}

	flag.Parse()

	if flags.version {
		return flags
	}

	if flags.language == "" || flags.name == "" {
		fmt.Printf("%s%s✘ Error: language (-lang) and name (-name) are required%s\n\n", bold, red, reset)
		flag.Usage()
		os.Exit(1)
	}

	flags.name = cleanProjectName(flags.name)
	if flags.name == "" {
		fmt.Printf("%s%s✘ Error: invalid project name%s\n", bold, red, reset)
		os.Exit(1)
	}

	flags.language = strings.ToLower(flags.language)
	if err := os.MkdirAll(flags.output, 0755); err != nil {
		fmt.Printf("%s%s✘ Error creating output directory: %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	return flags
}

func cleanProjectName(name string) string {
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' {
			return r
		}
		return -1
	}, name)

	name = strings.ToLower(name)
	name = strings.Trim(name, "-_")
	return name
}

func isValidLanguage(lang string) bool {
	validLanguages := map[string]bool{
		"go":     true,
		"ts":     true,
		"c":      true,
		"cxx":    true,
		"csharp": true,
		"python": true,
	}
	return validLanguages[lang]
}

func loadLanguageConfig(language string) (*ProjectConfig, error) {
	configPath := filepath.Join("internal", "config", fmt.Sprintf("%s.json", language))
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config ProjectConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func createProject(flags CLIFlags, config *ProjectConfig) error {
	projectDir := filepath.Join(flags.output, flags.name)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	fmt.Printf("\n%sCreating project structure...%s\n", blue, reset)
	for _, file := range config.Files {
		filePath := filepath.Join(projectDir, file)

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %v", file, err)
		}

		if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %v", file, err)
		}

		fmt.Printf("  %s%s✓%s %s\n", dim, green, reset, file)
	}

	return nil
}
