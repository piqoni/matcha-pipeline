package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gomarkdown/markdown"
)

func main() {
	files, err := filepath.Glob("docs/*.md")
	if err != nil {
		panic(err)
	}

	content := map[string]string{}

	for _, path := range files {
		b, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		html := markdown.ToHTML(b, nil, nil)
		date := strings.TrimSuffix(filepath.Base(path), ".md")
		content[date] = string(html)
	}

	keys := make([]string, 0, len(content))
	for k := range content {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	out, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("docs/content.js",
		[]byte("window.DIGESTS = "+string(out)),
		0644,
	)
	if err != nil {
		panic(err)
	}
}
