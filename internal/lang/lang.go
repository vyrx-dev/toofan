package lang

import (
	"embed"
	"io/fs"
	"math/rand"
	"sort"
	"strings"
)

//go:embed data
var dataFS embed.FS

type Snippet struct {
	Topic   string
	Content string
}

type langData struct {
	Name     string
	Words    []string
	Snippets []Snippet
}

var languages = map[string]*langData{}

// Names holds code language names (excludes english), sorted
var Names []string

// parseLesson extracts a snippet from a lesson file.
// Leading comments are stripped from the typed content.
// Supports // (go/js/dart), # (shell), and -- (lua) comment styles.
// The first "Topic:" comment becomes the display heading.
func parseLesson(content string) []Snippet {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	var topic string
	codeStart := 0

	// scan leading comments for metadata, find where code begins
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// check if this is a comment line (any supported prefix)
		isComment := strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "--")

		if !isComment {
			break
		}

		// extract "Topic:" from any comment style
		if topic == "" {
			for _, prefix := range []string{"// Topic: ", "# Topic: ", "-- Topic: "} {
				if strings.HasPrefix(trimmed, prefix) {
					topic = strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				}
			}
		}

		codeStart = i + 1
	}

	// skip blank lines between comments and code
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
	// Discover languages by looking at data/* directories
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

		// Read words.txt
		if raw, err := fs.ReadFile(dataFS, "data/"+name+"/words.txt"); err == nil {
			words := strings.Fields(string(raw))
			ld.Words = append(ld.Words, words...)
		}

		lessons, _ := fs.ReadDir(dataFS, "data/"+name)
		for _, lesson := range lessons {
			if lesson.IsDir() {
				continue
			}
			path := "data/" + name + "/" + lesson.Name()
			if raw, err := fs.ReadFile(dataFS, path); err == nil {
				snips := parseLesson(string(raw))
				ld.Snippets = append(ld.Snippets, snips...)
			}
		}

		if len(ld.Words) > 0 || len(ld.Snippets) > 0 {
			languages[name] = ld
			if name != "english" {
				Names = append(Names, name)
			}
		}
	}
	sort.Strings(Names)
}

// RandomWords picks random words for the word-mode typing test
func RandomWords(name string, count int) []string {
	ld, ok := languages[name]
	if !ok || len(ld.Words) == 0 {
		ld = languages["english"]
	}

	if len(ld.Words) == 0 {
		return []string{"hello", "world"}
	}

	out := make([]string, count)
	for i := range out {
		out[i] = ld.Words[rand.Intn(len(ld.Words))]
	}
	return out
}

// RandomSnippet picks a random code snippet for code-mode typing.
func RandomSnippet(name string) Snippet {
	ld, ok := languages[name]
	if !ok || len(ld.Snippets) == 0 {
		return Snippet{
			Topic:   "Fallback Words",
			Content: strings.Join(RandomWords(name, 50), " "),
		}
	}

	return ld.Snippets[rand.Intn(len(ld.Snippets))]
}

// GetSnippets returns all snippets for a given language
func GetSnippets(name string) []Snippet {
	if ld, ok := languages[name]; ok {
		return ld.Snippets
	}
	return nil
}


// Difficulty levels
type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

// GetHardWords loads difficult words from hard_words.txt
func GetHardWords() []string {
	raw, err := fs.ReadFile(dataFS, "data/english/hard_words.txt")
	if err != nil {
		// Fallback hard words if file not found
		return []string{
			"beautiful", "beginning", "believe", "business", "challenge",
			"character", "committee", "definitely", "embarrass", "necessary",
			"privilege", "recommend", "separate", "successful", "therefore",
		}
	}
	return strings.Fields(string(raw))
}

// RandomWordsWithDifficulty returns words based on difficulty level
func RandomWordsWithDifficulty(name string, count int, difficulty Difficulty) []string {
	ld, ok := languages[name]
	if !ok || len(ld.Words) == 0 {
		ld = languages["english"]
	}

	easyWords := ld.Words
	hardWords := GetHardWords()

	switch difficulty {
	case Easy:
		return randomFromSlice(easyWords, count)
	case Hard:
		if len(hardWords) > 0 {
			return randomFromSlice(hardWords, count)
		}
		return randomFromSlice(easyWords, count)
	case Medium:
		// Mix: 50% easy, 50% hard
		easyCount := count / 2
		hardCount := count - easyCount
		words := randomFromSlice(easyWords, easyCount)
		words = append(words, randomFromSlice(hardWords, hardCount)...)
		return words
	default:
		return randomFromSlice(easyWords, count)
	}
}

// randomFromSlice returns n random words from a slice
func randomFromSlice(slice []string, count int) []string {
	if len(slice) == 0 {
		return []string{"hello", "world"}
	}
	out := make([]string, count)
	for i := range out {
		out[i] = slice[rand.Intn(len(slice))]
	}
	return out
}

// GetWordCountForDifficulty returns total available words for a difficulty
func GetWordCountForDifficulty(difficulty Difficulty) int {
	switch difficulty {
	case Easy:
		if ld, ok := languages["english"]; ok {
			return len(ld.Words)
		}
	case Hard:
		return len(GetHardWords())
	case Medium:
		if ld, ok := languages["english"]; ok {
			return len(ld.Words) + len(GetHardWords())
		}
	}
	return 0
}
