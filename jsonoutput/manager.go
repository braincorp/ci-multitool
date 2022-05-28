package jsonoutput

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Manager struct {
	output *PulumiJSONOutput

	URNPrefix string
}

func NewManagerFromFile(path string) (*Manager, error) {
	m := &Manager{
		output: &PulumiJSONOutput{},
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	err = json.Unmarshal(fileContents, m.output)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal: %w", err)
	}

	return m, nil
}

// ShortSummaryString returns short one line summary of the changes
func (m *Manager) ShortSummaryString() string {
	warningCount := 0
	for _, d := range m.output.Diagnostics {
		if d.Severity == "error" {
			return "error"
		}
		if d.Severity == "warning" {
			warningCount++
		}
	}
	changeSummary := m.output.ChangeSummary
	if changeSummary.Create == 0 && changeSummary.Update == 0 {
		return "unchanged"
	}

	// TODO replaced

	var resParts []string
	if changeSummary.Create != 0 {
		resParts = append(resParts, fmt.Sprintf("created %d", changeSummary.Create))
	}
	if changeSummary.Update != 0 {
		resParts = append(resParts, fmt.Sprintf("updated %d", changeSummary.Update))
	}
	if warningCount != 0 {
		resParts = append(resParts, fmt.Sprintf("warn %d", warningCount))
	}

	return strings.Join(resParts, " || ")
}

func (m *Manager) stripURN(urn string) string {
	res := strings.TrimPrefix(urn, "urn:")
	res = strings.TrimPrefix(res, m.URNPrefix)
	return strings.TrimPrefix(res, "::")
}

func (m *Manager) urnParts(urn string) []string {
	parts := strings.Split(m.stripURN(urn), "::")
	var acc []string
	for _, part := range parts {
		acc = append(acc, strings.Split(part, "$")...)
	}
	return acc
}

func (m *Manager) urnParent(urn string) string {
	parts := m.urnParts(urn)
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(append(parts[:0], parts[:len(parts)-1]...), "::")
}

// TreeString returns a tree that looks like the `pulumi preview` console output
func (m *Manager) TreeString() string {
	urnToNode := make(map[string]Tree)

	tree := NewTree(m.URNPrefix)
	urnToNode[""] = tree
	for _, step := range m.output.Steps {
		strippedURN := m.stripURN(step.Urn)
		urnParts := m.urnParts(strippedURN)
		for i, urnPart := range urnParts {
			name := strings.Join(urnParts[:i+1], "::")
			_, ok := urnToNode[name]
			if ok {
				continue
			}
			parentName := m.urnParent(name)
			parentNode, ok := urnToNode[parentName]
			if !ok {
				fmt.Printf("missing parent for %s", name)
				continue
			}
			isLeaf := len(urnParts)-1 == i
			if isLeaf {
				break
			}
			node := parentNode.Add(urnPart)
			urnToNode[name] = node
		}
	}
	/*
		internalApps := tree.Add("InternalApps")
		docs := internalApps.Add("Docs")
		docs.Add("kubernetes:apps/v1:Deployment")
	*/
	res := tree.Print()

	/*
		maxLen := 0

			scanner := bufio.NewScanner(strings.NewReader(res))
			for scanner.Scan() {
				l := utf8.RuneCountInString(scanner.Text())
				if l > maxLen {
					maxLen = l
				}
			}
			nextLen := maxLen + 2

			scanner = bufio.NewScanner(strings.NewReader(res))
			for scanner.Scan() {
				text := scanner.Text()
				l := utf8.RuneCountInString(text)
				paddingLen := nextLen - l
				fmt.Printf("%s%s %s\n", text, strings.Repeat(" ", paddingLen), "asdf")
			}
	*/

	return res
}
