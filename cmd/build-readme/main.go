package main

import (
	"bytes"
	"os"
	"text/template"

	"app/internal/core"
	_ "app/internal/tests"
)

type TestInfo struct {
	Name        string
	Description string
}

func main() {
	tests := core.All()

	var testInfos []TestInfo
	for _, t := range tests {
		testInfos = append(testInfos, TestInfo{
			Name:        t.Name(),
			Description: t.Description(),
		})
	}

	tmplData, err := os.ReadFile("README.md.tmpl")
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("readme").Parse(string(tmplData))
	if err != nil {
		panic(err)
	}

	data := map[string]any{
		"Tests": testInfos,
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		panic(err)
	}

	if err := os.WriteFile("README.md", out.Bytes(), 0644); err != nil {
		panic(err)
	}
}
