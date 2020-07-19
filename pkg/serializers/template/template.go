package template

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"text/template"

	"github.com/ncarlier/trackr/pkg/model"
)

const defaultDataFormatTemplate = "{{.ClientIP}} {{.HostName}} - [{{.FormattedTS}}] \"GET {{.DocumentPath}} {{.Protocol}}\" 200 1 \"{{.DocumentReferrer}}\" \"{{.UserAgent}}\""

var tplCounter int32

type serializer struct {
	template *template.Template
}

// NewSerializer create new JSON serializer
func NewSerializer(format string) (*serializer, error) {
	if format == "" {
		format = defaultDataFormatTemplate
	}
	atomic.AddInt32(&tplCounter, 1)

	tmpl, err := template.New(fmt.Sprintf("template-%d", tplCounter)).Parse(format)
	if err != nil {
		return nil, err
	}
	s := &serializer{
		template: tmpl,
	}
	return s, nil
}

func (s *serializer) Serialize(pageview model.PageView) ([]byte, error) {
	serialized := new(bytes.Buffer)
	err := s.template.Execute(serialized, pageview)
	if err != nil {
		return []byte{}, err
	}
	serialized.WriteString("\n")

	return serialized.Bytes(), nil
}
