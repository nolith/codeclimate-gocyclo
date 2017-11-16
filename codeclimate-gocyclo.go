package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/codeclimate/cc-engine-go/engine"
)

const (
	defaultMaxComplexity = 9
	complexityField      = 0
	packageField         = 1
	functionField        = 2
	locationField        = 3
	fieldsCount          = 4
	locationFileField    = 0
	locationRowField     = 1
	locationColumnField  = 2
	locationFieldsCount  = 3
)

func main() {
	rootPath := "/code/"

	config, err := engine.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(1)
	}

	analysisFiles, err := engine.GoFileWalk(rootPath, engine.IncludePaths(rootPath, config))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing: %s", err)
		os.Exit(1)
	}

	maxComplexity := getMaxComplexity(config)

	for _, path := range analysisFiles {
		relativePath := strings.SplitAfter(path, rootPath)[1]

		lintFile(relativePath, maxComplexity)
	}
}

func lintFile(relativePath string, maxComplexity int) {
	gocyclo := exec.Command("gocyclo", "-over", strconv.Itoa(maxComplexity), relativePath)
	out, _ := gocyclo.CombinedOutput()

	reader := bytes.NewReader(out)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		// The output fields for each line are:
		//  <complexity> <package> <function> <file:row:column>

		fields := strings.SplitN(line, " ", fieldsCount)
		location := strings.SplitN(fields[locationField], ":", locationFieldsCount)
		row, err := strconv.ParseInt(location[locationRowField], 10, strconv.IntSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse gocyclo output: %v\n", line)
			os.Exit(1)
		}

		issue := &engine.Issue{
			Type:        "issue",
			Check:       "Gocyclo/Complexity/Function",
			Description: fmt.Sprintf("function `%s` in package `%s` has cyclomatic complexity of %s", fields[functionField], fields[packageField], fields[complexityField]),
			Categories:  []string{"Complexity"},
			Location: &engine.Location{
				Path:  relativePath,
				Lines: &engine.LinesOnlyPosition{Begin: int(row), End: int(row)},
			},
		}
		engine.PrintIssue(issue)

	}
}

func getMaxComplexity(config engine.Config) int {
	if subConfig, ok := config["config"].(map[string]interface{}); ok {
		if maxComplexity, ok := subConfig["over"].(string); ok {
			val, err := strconv.ParseInt(maxComplexity, 10, strconv.IntSize)
			if err == nil {
				return int(val)
			}
		}
	}
	return defaultMaxComplexity
}
