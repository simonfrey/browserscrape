package md_table

import (
	"bytes"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"golang.org/x/net/html"
	"strings"
)

// renderTable converts an HTML table node to a Markdown table and writes it.
func RenderTable(ctx converter.Context, w converter.Writer, node *html.Node) converter.RenderStatus {
	// Define a helper function to extract text content from nodes.
	var extractText func(*html.Node) string
	extractText = func(n *html.Node) string {
		if n.Type == html.TextNode {
			return strings.TrimSpace(n.Data)
		}
		var result string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result += extractText(c)
		}
		return strings.TrimSpace(result)
	}

	// Traverse the HTML nodes to extract table data.
	var table [][]string
	var currentRow []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "tr":
				currentRow = []string{}
			case "td", "th":
				text := extractText(n)
				if text != "" {
					currentRow = append(currentRow, text)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}

		if n.Type == html.ElementNode && n.Data == "tr" && len(currentRow) > 0 {
			table = append(table, currentRow)
		}
	}

	traverse(node)

	// Convert the table data to a Markdown format.
	if len(table) == 0 {
		return converter.RenderSuccess
	}

	// Calculate the maximum width for each column.
	columnWidths := make([]int, len(table[0]))
	for _, row := range table {
		for i, cell := range row {
			if i < len(columnWidths) && len(cell) > columnWidths[i] {
				columnWidths[i] = len(cell)
			}
		}
	}

	var markdown bytes.Buffer

	// Write the header.
	for i, cell := range table[0] {
		if i > 0 {
			markdown.WriteString(" | ")
		}
		markdown.WriteString(cell)
		padding := columnWidths[i] - len(cell)
		markdown.WriteString(strings.Repeat(" ", padding))
	}
	markdown.WriteString("\n")

	// Write the separator.
	for i, width := range columnWidths {
		if i > 0 {
			markdown.WriteString("-|-")
		}
		markdown.WriteString(strings.Repeat("-", width))
	}
	markdown.WriteString("\n")

	// Write the data rows.
	for i := 1; i < len(table); i++ {
		for j, cell := range table[i] {
			if j > 0 {
				markdown.WriteString(" | ")
			}
			markdown.WriteString(cell)
			if j < len(columnWidths) {
				padding := columnWidths[j] - len(cell)
				markdown.WriteString(strings.Repeat(" ", padding))
			}
		}
		markdown.WriteString("\n")
	}

	// Write the converted Markdown content
	w.WriteString("\n" + markdown.String() + "\n")

	// Return successful rendering status.
	return converter.RenderSuccess
}
