package processor

import (
	"bytes"
	"fmt"
	"html/template"
	"pagegen/vv"
	"path/filepath"
	"strings"
)

type CheckFn func(file *vv.File) bool

func NewHtmlPart(data map[string]any, isPart CheckFn) vv.ProcessorFn {
	return func(config *vv.Ctx, file *vv.File) error {
		if !isPart(file) {
			return nil
		}

		dir := strings.ReplaceAll(file.Dir, string(filepath.Separator), ".")
		name := file.Name
		if dir != "" {
			name = dir + "_" + name
		}
		name = "html_" + name

		tmpl, err := template.New(name).Parse(string(file.Body))
		if err != nil {
			return fmt.Errorf("error parsing template %q (%v)", file.SourcePath, err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return fmt.Errorf("error executing template %q (%v)", file.SourcePath, err)
		}

		data[name] = template.HTML(buf.String())
		file.Ignored = true
		file.Hidden = true
		return nil
	}
}

func NewHtml(data map[string]any, isHtml CheckFn) vv.ProcessorFn {
	return func(config *vv.Ctx, file *vv.File) error {
		if !isHtml(file) {
			return nil
		}

		tmpl, err := template.New(file.Name).Parse(string(file.Body))
		if err != nil {
			return fmt.Errorf("error parsing template %q (%v)", file.SourcePath, err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return fmt.Errorf("error executing template %q (%v)", file.SourcePath, err)
		}

		file.Body = buf.Bytes()
		return nil
	}
}
