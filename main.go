package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	major := flag.Bool("major", false, "Bump major version")
	minor := flag.Bool("minor", false, "Bump minor version")
	mainBranch := flag.String("main-branch", "main", "Name of the main branch")
	flag.Parse()

	var filePath string
	args := flag.Args()
	if len(args) > 0 {
		filePath = args[0]
	} else {
		filePath = "VERSION"
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("Error: VERSION file not found at %s\n", absPath)
		os.Exit(1)
	}

	currentBranch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		fmt.Printf("Error getting current branch: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(string(currentBranch)) == *mainBranch {
		fmt.Printf("Error: Cannot bump version on '%s' branch!\n", *mainBranch)
		os.Exit(1)
	}

	status, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		fmt.Printf("Error checking git status: %v\n", err)
		os.Exit(1)
	}
	if len(status) > 0 {
		fmt.Println("Error: Cannot bump because there are uncommitted changes!")
		os.Exit(1)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", absPath, err)
		os.Exit(1)
	}

	currentVersion := strings.TrimSpace(string(content))
	parts := strings.Split(currentVersion, ".")

	if len(parts) != 3 {
		fmt.Printf("Invalid version format in %s\n", absPath)
		os.Exit(1)
	}

	majorVer, _ := strconv.Atoi(parts[0])
	minorVer, _ := strconv.Atoi(parts[1])
	patchVer, _ := strconv.Atoi(parts[2])

	// Bump version
	if *major {
		majorVer++
		minorVer = 0
		patchVer = 0
	} else if *minor {
		minorVer++
		patchVer = 0
	} else {
		patchVer++
	}

	newVersion := fmt.Sprintf("%d.%d.%d", majorVer, minorVer, patchVer)

	err = os.WriteFile(absPath, []byte(newVersion), 0o644)
	if err != nil {
		fmt.Printf("Error writing to %s: %v\n", absPath, err)
		os.Exit(1)
	}

	_, err = exec.Command("git", "diff", "--exit-code", "--", absPath).CombinedOutput()
	if err == nil {
		fmt.Println("Error: Version has not changed!")
		os.Exit(1)
	}

	commitMsg := fmt.Sprintf("release: v%s", newVersion)
	cmd := exec.Command("git", "commit", "-m", commitMsg, "--", absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error committing the version change: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Version bumped from %s to %s and committed\n", currentVersion, newVersion)
}
