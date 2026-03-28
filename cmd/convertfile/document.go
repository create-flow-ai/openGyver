package convertfile

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/yuin/goldmark"
)

// mdToHTML converts Markdown to HTML using goldmark (CommonMark compliant).
func mdToHTML(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading markdown: %w", err)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(src, &buf); err != nil {
		return fmt.Errorf("converting markdown: %w", err)
	}

	html := wrapHTML(buf.String(), opts.InputPath)
	return os.WriteFile(opts.OutputPath, []byte(html), 0644)
}

// htmlToMD converts HTML to Markdown.
func htmlToMD(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading HTML: %w", err)
	}

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(string(src))
	if err != nil {
		return fmt.Errorf("converting HTML to Markdown: %w", err)
	}

	return os.WriteFile(opts.OutputPath, []byte(markdown+"\n"), 0644)
}

// mdToText converts Markdown to plain text by stripping formatting.
func mdToText(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading markdown: %w", err)
	}

	text := stripMarkdown(string(src))
	return os.WriteFile(opts.OutputPath, []byte(text), 0644)
}

// htmlToText converts HTML to plain text by stripping all tags.
func htmlToText(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading HTML: %w", err)
	}

	text := stripHTML(string(src))
	return os.WriteFile(opts.OutputPath, []byte(text), 0644)
}

// textToHTML wraps plain text in a basic HTML document with <pre>.
func textToHTML(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading text: %w", err)
	}

	escaped := htmlEscape(string(src))
	html := wrapHTML("<pre>"+escaped+"</pre>", opts.InputPath)
	return os.WriteFile(opts.OutputPath, []byte(html), 0644)
}

// textToMD wraps plain text in a Markdown code fence.
func textToMD(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading text: %w", err)
	}

	md := "```\n" + string(src)
	if !strings.HasSuffix(md, "\n") {
		md += "\n"
	}
	md += "```\n"
	return os.WriteFile(opts.OutputPath, []byte(md), 0644)
}

// renderMarkdownToHTML is a shared helper for markdown → HTML conversion.
func renderMarkdownToHTML(src []byte) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(src, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// --- Helpers ---

func wrapHTML(body, title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>%s</title>
</head>
<body>
%s
</body>
</html>
`, htmlEscape(title), body)
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

var (
	reHTMLTag     = regexp.MustCompile(`<[^>]*>`)
	reMDHeader    = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	reMDBold      = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reMDItalic    = regexp.MustCompile(`\*(.+?)\*`)
	reMDCode      = regexp.MustCompile("`([^`]+)`")
	reMDLink      = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	reMDImage     = regexp.MustCompile(`!\[([^\]]*)\]\([^)]+\)`)
	reMDListItem  = regexp.MustCompile(`(?m)^[\s]*[-*+]\s+`)
	reMDNumList   = regexp.MustCompile(`(?m)^[\s]*\d+\.\s+`)
	reMDBlockquote = regexp.MustCompile(`(?m)^>\s?`)
	reMDHR        = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	reMultiNewline = regexp.MustCompile(`\n{3,}`)
)

func stripHTML(s string) string {
	// Replace <br> and block elements with newlines
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")
	s = strings.ReplaceAll(s, "<br />", "\n")
	s = regexp.MustCompile(`</(?:p|div|h[1-6]|li|tr|blockquote)>`).ReplaceAllString(s, "\n")
	s = reHTMLTag.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = reMultiNewline.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s) + "\n"
}

func stripMarkdown(s string) string {
	s = reMDImage.ReplaceAllString(s, "$1")
	s = reMDLink.ReplaceAllString(s, "$1")
	s = reMDBold.ReplaceAllString(s, "$1")
	s = reMDItalic.ReplaceAllString(s, "$1")
	s = reMDCode.ReplaceAllString(s, "$1")
	s = reMDHeader.ReplaceAllString(s, "")
	s = reMDListItem.ReplaceAllString(s, "  ")
	s = reMDNumList.ReplaceAllString(s, "  ")
	s = reMDBlockquote.ReplaceAllString(s, "")
	s = reMDHR.ReplaceAllString(s, "")
	s = reMultiNewline.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s) + "\n"
}
