package translate

import (
	"fmt"
	"io"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

type TemplateExecuter interface {
	Execute(wr io.Writer, data interface{}) error
	Parse(text string) (TemplateExecuter, error)
	Clone() (TemplateExecuter, error)
}

type Template struct {
	x TemplateExecuter
}

func (t Template) Ext(v map[string]interface{}, file string) ([]byte, error) {
	templateEngine, err := t.x.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone text template error: %w", err)
	}

	tmp, err := templateEngine.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("parse text template error: %w", err)
	}

	return Execute(v, tmp)
}

func NewTemplate() *Template {
	return &Template{
		x: &TextAdapter{textTemplate.New("txt").Funcs(sprig.TxtFuncMap())},
	}
}
