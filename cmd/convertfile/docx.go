package convertfile

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// --- DOCX Reader ---

// docxBody represents the simplified structure of word/document.xml.
type docxBody struct {
	Paragraphs []docxParagraph `xml:"body>p"`
}

type docxParagraph struct {
	Runs []docxRun `xml:"r"`
}

type docxRun struct {
	Text []docxText `xml:"t"`
}

type docxText struct {
	Content string `xml:",chardata"`
}

// readDOCXParagraphs extracts paragraph text from a DOCX file.
func readDOCXParagraphs(path string) ([]string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("opening DOCX: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			var body docxBody
			if err := xml.NewDecoder(rc).Decode(&body); err != nil {
				return nil, fmt.Errorf("parsing document.xml: %w", err)
			}

			var paragraphs []string
			for _, p := range body.Paragraphs {
				var line strings.Builder
				for _, run := range p.Runs {
					for _, t := range run.Text {
						line.WriteString(t.Content)
					}
				}
				paragraphs = append(paragraphs, line.String())
			}
			return paragraphs, nil
		}
	}
	return nil, fmt.Errorf("word/document.xml not found in DOCX")
}

// --- DOCX Writer ---

// writeDOCX creates a minimal valid DOCX file from paragraphs of text.
func writeDOCX(path string, paragraphs []string) error {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// [Content_Types].xml
	writeZipFile(w, "[Content_Types].xml", contentTypesXML)

	// _rels/.rels
	writeZipFile(w, "_rels/.rels", relsXML)

	// word/_rels/document.xml.rels
	writeZipFile(w, "word/_rels/document.xml.rels", documentRelsXML)

	// word/document.xml
	var bodyParts strings.Builder
	for _, p := range paragraphs {
		escaped := xmlEscape(p)
		bodyParts.WriteString(fmt.Sprintf(`<w:p><w:r><w:t xml:space="preserve">%s</w:t></w:r></w:p>`, escaped))
	}
	docXML := fmt.Sprintf(documentXMLTemplate, bodyParts.String())
	writeZipFile(w, "word/document.xml", docXML)

	if err := w.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func writeZipFile(w *zip.Writer, name, content string) {
	f, _ := w.Create(name)
	f.Write([]byte(content))
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

const contentTypesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`

const relsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`

const documentRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`

const documentXMLTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>%s</w:body>
</w:document>`

// --- Conversion functions ---

func docxToText(opts ConvertOpts) error {
	paragraphs, err := readDOCXParagraphs(opts.InputPath)
	if err != nil {
		return err
	}
	text := strings.Join(paragraphs, "\n") + "\n"
	return os.WriteFile(opts.OutputPath, []byte(text), 0644)
}

func docxToHTML(opts ConvertOpts) error {
	paragraphs, err := readDOCXParagraphs(opts.InputPath)
	if err != nil {
		return err
	}

	var body strings.Builder
	for _, p := range paragraphs {
		if p == "" {
			body.WriteString("<br>\n")
		} else {
			body.WriteString(fmt.Sprintf("<p>%s</p>\n", htmlEscape(p)))
		}
	}

	html := wrapHTML(body.String(), opts.InputPath)
	return os.WriteFile(opts.OutputPath, []byte(html), 0644)
}

func docxToMD(opts ConvertOpts) error {
	paragraphs, err := readDOCXParagraphs(opts.InputPath)
	if err != nil {
		return err
	}
	text := strings.Join(paragraphs, "\n\n") + "\n"
	return os.WriteFile(opts.OutputPath, []byte(text), 0644)
}

func textToDOCX(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading text: %w", err)
	}
	paragraphs := strings.Split(string(src), "\n")
	return writeDOCX(opts.OutputPath, paragraphs)
}

func mdToDOCX(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading markdown: %w", err)
	}
	// Strip markdown formatting, split into paragraphs
	text := stripMarkdown(string(src))
	paragraphs := strings.Split(text, "\n")
	return writeDOCX(opts.OutputPath, paragraphs)
}

func htmlToDOCX(opts ConvertOpts) error {
	src, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return fmt.Errorf("reading HTML: %w", err)
	}
	text := stripHTML(string(src))
	paragraphs := strings.Split(text, "\n")
	return writeDOCX(opts.OutputPath, paragraphs)
}
