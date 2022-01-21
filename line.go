package dbbr

import (
	"bufio"
	"os"
	"strings"
)

// Line repressents a line in DBBR configration file
type Line struct {
	Number int
	Value  string
}

// ConvertToLines converts a slice of strings to Lines.
func ConvertToLines(strLines []string) []Line {
	//var lines []Line
	//for i, v := range strLines {
	//	append(lines, Line{i, v})
	//}
	//return lines
	lines := make([]Line, len(strLines))
	for i, v := range strLines {
		lines[i] = Line{i, v}
	}
	return lines
}

// ReadLinesFromFile converts textual file to lines.
func ReadLinesFromFile(path string) ([]Line, error) {
	lines := []Line{}
	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i += 1 {
		line := Line{i, scanner.Text()}
		lines = append(lines, line)
	}

	if err = scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

// SectionLines creates sections of rules from lines.
func SectionLines(lines []Line) [][]Line {
	sections := [][]Line{}
	currSection := []Line{}
	for _, line := range lines {
		if len(strings.TrimSpace(line.Value)) == 0 {
			if len(currSection) != 0 {
				sections = append(sections, currSection)
				currSection = []Line{}
			}
			continue
		}
		currSection = append(currSection, line)
	}
	if len(currSection) != 0 {
		sections = append(sections, currSection)
	}
	return sections
}
