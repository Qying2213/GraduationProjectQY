package utils

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList, extension.DefinitionList, extension.Footnote),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

func MarkdownToHTML(in string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(in), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
