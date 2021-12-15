package translate

import (
	"io"
	textTemplate "text/template"
)

type TextAdapter struct {
	template *textTemplate.Template
}

func (a *TextAdapter) Execute(wr io.Writer, data interface{}) error {
	return a.template.Execute(wr, data) //nolint: wrapcheck
}

func (a *TextAdapter) Parse(text string) (TemplateExecuter, error) {
	tmp, err := a.template.Parse(text)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return &TextAdapter{tmp}, nil
}

func (a *TextAdapter) Clone() (TemplateExecuter, error) {
	tmp, err := a.template.Clone()
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return &TextAdapter{tmp}, nil
}
