package content

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var (
	markdown = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	// 双重保险：goldmark 不开启 unsafe，再按 UGC 策略净化一次。
	sanitizer = bluemonday.UGCPolicy()
)

// Render 把 Markdown 渲染成 HTML。不使用 html.WithUnsafe，原始 HTML 会被转义，
// 输出再经 bluemonday UGC 策略净化后才允许交给模板的 safeHTML。
func Render(source string) (string, error) {
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return sanitizer.Sanitize(buf.String()), nil
}

var (
	fencedCodePattern  = regexp.MustCompile("(?s)```.*?```")
	imagePattern       = regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`)
	linkPattern        = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	inlineCodePattern  = regexp.MustCompile("`([^`]*)`")
	headingPattern     = regexp.MustCompile(`(?m)^#{1,6}\s*`)
	listPattern        = regexp.MustCompile(`(?m)^\s*(?:[-*+]|\d+\.)\s+`)
	blockquotePattern  = regexp.MustCompile(`(?m)^>\s?`)
	htmlPattern        = regexp.MustCompile(`<[^>]+>`)
	whitespacePattern  = regexp.MustCompile(`\s+`)
	emphasisDelimiters = strings.NewReplacer("**", "", "__", "", "~~", "", "*", "", "_", "")
)

// Excerpt 从 Markdown 源码生成纯文本摘要。
func Excerpt(source string, limit int) string {
	text := fencedCodePattern.ReplaceAllString(source, " ")
	text = imagePattern.ReplaceAllString(text, "$1")
	text = linkPattern.ReplaceAllString(text, "$1")
	text = inlineCodePattern.ReplaceAllString(text, "$1")
	text = headingPattern.ReplaceAllString(text, "")
	text = listPattern.ReplaceAllString(text, "")
	text = blockquotePattern.ReplaceAllString(text, "")
	text = htmlPattern.ReplaceAllString(text, " ")
	text = emphasisDelimiters.Replace(text)
	text = whitespacePattern.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit]) + "…"
}
