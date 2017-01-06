package markdown

import (
	"bytes"
	"text/template"
)

const tpl = `
<html>
<head>
<style>
{{.__css__}}
</style>
</head>

<body>
<article class="markdown-body">
{{.__body__}}
</article>
</body>

</html>
`

const (
	githubCommonHTMLFlags = 0 |
		HTML_USE_XHTML |
		HTML_USE_SMARTYPANTS |
		HTML_SMARTYPANTS_FRACTIONS |
		HTML_SMARTYPANTS_LATEX_DASHES

	githubCommonExtensions = 0 |
		EXTENSION_NO_INTRA_EMPHASIS |
		EXTENSION_TABLES |
		EXTENSION_FENCED_CODE |
		EXTENSION_AUTOLINK |
		EXTENSION_STRIKETHROUGH |
		EXTENSION_SPACE_HEADERS |
		EXTENSION_HEADER_IDS |
		EXTENSION_BACKSLASH_LINE_BREAK |
		EXTENSION_DEFINITION_LISTS
)

// GithubForRender github render
type GithubForRender struct {
	Tpl string
	Map map[string]interface{}
}

// GithubMarkdown renders markdown to github style
func GithubMarkdown(in []byte, hasCatalog bool) ([]byte, error) {
	flg := githubCommonHTMLFlags
	if hasCatalog {
		flg |= HTML_TOC
	}
	render := HtmlRenderer(flg, "", css)
	body := MarkdownOptions(in, render, Options{
		Extensions: githubCommonExtensions,
	})
	m := map[string]interface{}{
		"__css__":  css,
		"__body__": string(body),
	}
	tpl, err := template.New("markdown").Parse(tpl)
	if err != nil {
		return nil, err
	}
	var out = new(bytes.Buffer)
	err = tpl.Execute(out, m)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
