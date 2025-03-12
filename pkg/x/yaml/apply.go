package yaml

import (
	"bytes"
	"fmt"
	"text/template"
)

func ApplyTmpl(tmplStr string, tmplData interface{}, debug bool) ([]byte, error) {
	tmpl, err := template.New("tmpManifest").Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		return nil, err
	}

	// debug
	if debug {
		fmt.Println(buf.String())
	}
	return buf.Bytes(), nil
}
