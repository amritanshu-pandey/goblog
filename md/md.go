package md

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Title string    `yaml:"title"`
	Date  time.Time `yaml:"date"`
	Draft bool      `yaml:"draft"`
	Tags  []string  `yaml:"tags"`
}

func parseMetadata(mdContent string) (Metadata, error) {
	metadata := Metadata{}
	err := yaml.Unmarshal([]byte(mdContent), &metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}

type Post struct {
	FileName string
	BodyHTML []byte
	Metadata Metadata
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func extractMetadata(mdFile *os.File) (string, string, error) {
	scanner := bufio.NewScanner(mdFile)
	metadata := []string{}
	body := []string{}
	sepCounter := 0
	line := 0
	hasMetadata := false
	for scanner.Scan() {
		current := scanner.Text()
		line += 1

		if line == 1 && current == "---" {
			hasMetadata = true
		}

		if sepCounter < 2 && hasMetadata {
			if current == "---" {
				sepCounter += 1
			} else {
				metadata = append(metadata, current)
			}
		} else {
			body = append(body, current)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Scanning complete")
			break
		}
	}
	if !hasMetadata {
		return "", "", fmt.Errorf("Metadata not  found in post: %s", mdFile.Name())
	}
	return strings.Join(metadata, "\n"), strings.Join(body, "\n"), nil
}

func Posts(path string) (map[string]Post, error) {
	source := os.DirFS(path)
	dirEntries, err := fs.Glob(source, "*.md")
	if err != nil {
		return nil, err
	}

	posts := make(map[string]Post)

	for _, p := range dirEntries {
		mdPath := fmt.Sprintf("%s/%s", path, p)
		mdFile, _ := os.Open(mdPath)

		metadataRaw, bodyRaw, err := extractMetadata(mdFile)
		if err != nil {
			fmt.Printf("Error: %s, Skipping post '%s'", err, p)
			continue
		}

		metadataParsed, err := parseMetadata(metadataRaw)
		if err != nil {
			fmt.Printf("Error: %s, Skipping post '%s'", err, p)
			continue
		}

		fileName, _ := strings.CutSuffix(p, ".md")
		post := Post{
			FileName: fileName,
			BodyHTML: mdToHTML([]byte(bodyRaw)),
			Metadata: metadataParsed,
		}

		posts[fileName] = post
	}

	return posts, nil
}

func ActivePosts(path string) (map[string]Post, error) {
	allPosts, err := Posts(path)
	if err != nil {
		return nil, err
	}

	activePosts := make(map[string]Post)

	for p := range allPosts {
		if !allPosts[p].Metadata.Draft {
			activePosts[p] = allPosts[p]
		}
	}

	return activePosts, nil
}

func SortedPostsByTitle(path string) ([]string, error) {
	activePosts, err := ActivePosts(path)
	if err != nil {
		return nil, err
	}

	titles := []string{}

	for k := range activePosts {
		titles = append(titles, k)
	}
	return sort.StringSlice(titles), nil
}

type Kv struct {
	Key   string
	Value Post
}

func SortedPostsByDate(path string) ([]Kv, error) {
	activePosts, err := ActivePosts(path)
	if err != nil {
		return nil, err
	}

	var titles []Kv

	for k, v := range activePosts {
		titles = append(titles, Kv{k, v})
	}

	sort.Slice(titles, func(i, j int) bool {
		return titles[i].Value.Metadata.Date.UnixMicro() > titles[j].Value.Metadata.Date.UnixMicro()
	})
	return titles, nil
}
