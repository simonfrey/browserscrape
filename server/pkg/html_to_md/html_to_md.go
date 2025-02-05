package html_to_md

import (
	"fmt"
	"regexp"
	"server/pkg/md_table"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"strings"
)

func HTMLToMD(url, html string, cleanForLLM bool) (string, error) {
	// This function converts HTML to markdown
	// It uses the html-to-markdown library
	html = removeHTMLAttributes(html)
	html = getBody(html)

	mC := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)
	mC.Register.RendererFor("table", converter.TagTypeInline, md_table.RenderTable, converter.PriorityLate)

	markdown, err := mC.ConvertString(html, converter.WithDomain(strings.Split(url, "/")[0]))

	if err != nil {
		return "", fmt.Errorf("error converting html to markdown: %w", err)
	}

	if !cleanForLLM {
		// Return the full markdown
		return markdown, nil
	}

	// Clean the markdown for LLM
	return cleanMarkdown(markdown), nil
}

var htmlTagRegex = regexp.MustCompile(`<\s*([a-zA-Z0-9]+)[^>]*>`)

func removeHTMLAttributes(input string) string {
	// This regex matches HTML tags with attributes and replaces them with just the tag name
	return htmlTagRegex.ReplaceAllString(input, "<$1>")
}

var htmlBodyRegexp = regexp.MustCompile(`(?is)<body>(.*)</body>`)

func getBody(input string) string {
	// This regex matches HTML tags with attributes and replaces them with just the tag name
	return htmlBodyRegexp.FindString(input)
}

var markdownImageRegex = regexp.MustCompile(`\!\[.*?\]\(.+?\)`)
var markdownLinkRegex = regexp.MustCompile(`\[(.*?)\]\(.+?\)`)
var orderedListRegex = regexp.MustCompile(`^\d+\..+$`)
var urlRegex = regexp.MustCompile(`https?://[^\s]+?`)
var codeRegex = regexp.MustCompile(`\x60.*\x60`)

func cleanMarkdown(markdownIn string) string {
	markdown := markdownIn
	for {
		oldMarkdown := markdown

		// Remove code blocks
		markdown = codeRegex.ReplaceAllString(markdown, "")

		lines := strings.Split(markdown, "\n")
		markdown = ""
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			line = strings.ReplaceAll(line, `‎`, "")
			line = markdownImageRegex.ReplaceAllString(line, "")
			line = urlRegex.ReplaceAllString(line, "")
			lineWithoutLink := markdownLinkRegex.ReplaceAllString(line, "")
			line = markdownLinkRegex.ReplaceAllString(line, "$1")

			if strings.HasPrefix(line, "#") {
				markdown += "\n" + line + "\n"
				continue
			}

			if strings.HasPrefix(line, "|") ||
				strings.HasPrefix(line, "-") ||
				orderedListRegex.MatchString(line) {
				markdown += line + "\n"
				continue
			}

			// Check if line is too short to be a paragraph
			if len(lineWithoutLink) < 30 {
				continue
			}

			markdown += line + "\n"
		}
		if oldMarkdown == markdown {
			break
		}
	}

	return markdown
}
