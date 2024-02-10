package internal

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var (
	CatppuccinMacchiato = styles.Get("catppuccin-macchiato")
	HTMLFormatter       = html.New()
)

func Highlight(language, content string) (string, error) {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	if err := HTMLFormatter.Format(&buf, CatppuccinMacchiato, iterator); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func MustHighlight(language, content string) string {
	highlit, err := Highlight(language, content)
	if err != nil {
		panic(fmt.Sprintf("could not highlight: %v", err))
	}

	return highlit
}
