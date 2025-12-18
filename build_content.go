package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gomarkdown/markdown"
)

const placeholder = "/* __DIGEST_DATA__ */"

func main() {
	// Load markdown files
	files, err := filepath.Glob("docs/*.md")
	if err != nil {
		panic(err)
	}

	digests := map[string]string{}

	for _, path := range files {
		b, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		html := markdown.ToHTML(b, nil, nil)
		date := strings.TrimSuffix(filepath.Base(path), ".md")
		digests[date] = string(html)
	}

	// Sort dates descending
	dates := make([]string, 0, len(digests))
	for d := range digests {
		dates = append(dates, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	// JSON encode digests
	data, err := json.Marshal(digests)
	if err != nil {
		panic(err)
	}

	// Build inline JS app
	var js bytes.Buffer
	js.WriteString("const DIGESTS = ")
	js.Write(data)
	js.WriteString(";\n\n")

	js.WriteString(`
const datesEl = document.getElementById("dates");
const contentEl = document.getElementById("content");

const dates = Object.keys(DIGESTS).sort().reverse();

for (const date of dates) {
  const div = document.createElement("div");
  div.textContent = date;
  div.className = "date";
  div.onclick = () => {
    document.querySelectorAll(".date").forEach(d => d.classList.remove("active"));
    div.classList.add("active");
    contentEl.innerHTML = DIGESTS[date];
  };
  datesEl.appendChild(div);
}

if (dates.length > 0) {
  datesEl.firstChild.click();
}
`)

	// Read template HTML
	htmlTemplate, err := os.ReadFile("template/index.html")
	if err != nil {
		panic(err)
	}

	if !bytes.Contains(htmlTemplate, []byte(placeholder)) {
		panic("placeholder not found in index.html")
	}

	// Inject JS
	finalHTML := bytes.Replace(
		htmlTemplate,
		[]byte(placeholder),
		js.Bytes(),
		1,
	)

	// Overwrite index.html
	err = os.WriteFile("pre-docs/index.html", finalHTML, 0644)
	if err != nil {
		panic(err)
	}
}
