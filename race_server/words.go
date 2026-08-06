package main

import (
	"embed"
	"io/fs"
	"log"
	"math/rand"
	"path"
	"strings"
)

//go:embed data
var dataFS embed.FS

type Snippet struct {
	Topic   string
	Content string
}

type langData struct {
	Name        string
	EasyWords   []string
	MediumWords []string
	HardWords   []string
	AllWords    []string
	Snippets    []Snippet
}

var languages = map[string]*langData{}

func parseLesson(content string) []Snippet {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	var topic string
	codeStart := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		isComment := strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "--")
		if !isComment {
			break
		}
		if topic == "" {
			for _, prefix := range []string{"// Topic: ", "# Topic: ", "-- Topic: "} {
				if strings.HasPrefix(trimmed, prefix) {
					topic = strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				}
			}
		}
		codeStart = i + 1
	}

	for codeStart < len(lines) && strings.TrimSpace(lines[codeStart]) == "" {
		codeStart++
	}

	codeText := strings.TrimSpace(strings.Join(lines[codeStart:], "\n"))
	if len(codeText) == 0 {
		return nil
	}
	if topic == "" {
		topic = "Code Snippet"
	}
	return []Snippet{{Topic: topic, Content: codeText}}
}

func init() {
	entries, err := fs.ReadDir(dataFS, "data")
	if err != nil {
		return
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		ld := &langData{Name: name}

		if raw, err := fs.ReadFile(dataFS, "data/"+name+"/easy.txt"); err == nil {
			ld.EasyWords = strings.Fields(string(raw))
			ld.AllWords = append(ld.AllWords, ld.EasyWords...)
		}
		if raw, err := fs.ReadFile(dataFS, "data/"+name+"/medium.txt"); err == nil {
			ld.MediumWords = strings.Fields(string(raw))
			ld.AllWords = append(ld.AllWords, ld.MediumWords...)
		}
		if raw, err := fs.ReadFile(dataFS, "data/"+name+"/hard.txt"); err == nil {
			ld.HardWords = strings.Fields(string(raw))
			ld.AllWords = append(ld.AllWords, ld.HardWords...)
		}

		// Load snippets from both data/<lang>/ and data/<lang>/lessons/.
		// This matches the app-side content layout.
		loadSnippets := func(dir string) {
			files, err := fs.ReadDir(dataFS, dir)
			if err != nil {
				return
			}
			for _, f := range files {
				if f.IsDir() || strings.HasSuffix(f.Name(), ".txt") {
					continue
				}
				p := path.Join(dir, f.Name())
				raw, err := fs.ReadFile(dataFS, p)
				if err != nil {
					continue
				}
				snips := parseLesson(string(raw))
				ld.Snippets = append(ld.Snippets, snips...)
			}
		}
		loadSnippets("data/" + name)
		loadSnippets("data/" + name + "/lessons")

		if len(ld.AllWords) > 0 || len(ld.Snippets) > 0 {
			languages[name] = ld
		}
	}
	log.Printf("loaded %d languages into race server", len(languages))
}

func generateText(mode, lang, difficulty string) string {
	ld, ok := languages[lang]
	if !ok {
		ld = languages["english"]
	}
	if ld == nil {
		log.Printf("generateText fallback: no language data for mode=%s lang=%s diff=%s", mode, lang, difficulty)
		return "hello world"
	}

	if mode == "code" && len(ld.Snippets) > 0 {
		s := ld.Snippets[rand.Intn(len(ld.Snippets))]
		return s.Content
	}

	var pool []string
	switch difficulty {
	case "easy":
		pool = ld.EasyWords
	case "medium":
		pool = ld.MediumWords
	case "hard":
		pool = ld.HardWords
	}

	if len(pool) == 0 {
		pool = ld.AllWords
	}
	if len(pool) == 0 {
		log.Printf("generateText fallback: empty word pool for lang=%s diff=%s", lang, difficulty)
		return "hello world"
	}

	count := 40
	if difficulty == "hard" {
		count = 30
	} else if difficulty == "easy" {
		count = 50
	}

	res := make([]string, count)
	for i := range res {
		res[i] = pool[rand.Intn(len(pool))]
	}
	return strings.Join(res, " ")
}
