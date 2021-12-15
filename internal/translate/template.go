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

var GlobalTemplate *Template

type Template struct {
	TXT TemplateExecuter
}

func (t Template) Ext(v map[string]interface{}, file string, funcList []string) ([]byte, error) {
	templateEngine, err := t.TXT.Clone()
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
		TXT: &TextAdapter{textTemplate.New("txt").Funcs(sprig.TxtFuncMap())},
	}
}

func SetGlobalTemplate(t *Template) {
	GlobalTemplate = t
}
