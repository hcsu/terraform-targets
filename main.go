package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tt [-t|-m] [plan|apply] [targets...]")
		os.Exit(1)
	}

	flag := os.Args[1]

	// Extract targets from terraform plan output
	if flag == "-t" || flag == "-m" {
		extractTargets(flag == "-m")
		return
	}

	// Execute terraform with targets
	cmd := os.Args[1]
	targets := os.Args[2:]

	args := []string{cmd}
	for _, target := range targets {
		args = append(args, "-target", target)
	}

	terraform := exec.Command("terraform", args...)
	terraform.Stdout = os.Stdout
	terraform.Stderr = os.Stderr
	terraform.Stdin = os.Stdin

	if err := terraform.Run(); err != nil {
		os.Exit(1)
	}
}

func extractTargets(moduleLevel bool) {
	cmd := exec.Command("terraform", "plan")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Read stderr in background and output to os.Stderr
	go io.Copy(os.Stderr, stderr)

	targets := parseTargets(stdout, moduleLevel)

	if err := cmd.Wait(); err != nil {
		// If terraform command failed, exit with error
		os.Exit(1)
	}

	// Output targets with quotes and backslash continuation
	for i, target := range targets {
		if i < len(targets)-1 {
			fmt.Printf("'%s' \\\n", target)
		} else {
			fmt.Printf("'%s'\n", target)
		}
	}
}

func parseTargets(r io.Reader, moduleLevel bool) []string {
	scanner := bufio.NewScanner(r)
	re := regexp.MustCompile(`#\s+(.+?)\s+will be`)
	moduleRe := regexp.MustCompile(`^module\.[^.]+`)

	var targets []string
	seen := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "will be") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}

		target := matches[1]

		// Handle module level extraction
		if moduleLevel && strings.HasPrefix(target, "module.") {
			if match := moduleRe.FindString(target); match != "" {
				target = match
			}
		}

		// Deduplicate
		if !seen[target] {
			targets = append(targets, target)
			seen[target] = true
		}
	}

	return targets
}
