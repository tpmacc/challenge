package challenge

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called inappropriately")
	},
}

var layout = template.Must(
	template.
		New("layout.html").
		Funcs(layoutFuncs).
		ParseFiles("templates/layout.html"),
)

var Templates = template.Must(template.New("t").ParseGlob("templates/**/*.html"))

//var errorTemplate = `
//<html>
//	<body>
//		<h1>Error rendering template %s</h1>
//		<p>%s</p>
//	</body>
//</html>
//`

func RenderTemplate(w io.Writer, name string, data interface{}) error {

	funcs := buildTemplateFuncMap(name, data)

	layoutClone, _ := layout.Clone()
	layoutClone.Funcs(funcs)
	err := layoutClone.Execute(w, data)

	//if err != nil {
	//	http.Error(
	//		w,
	//		fmt.Sprintf(errorTemplate, name, err),
	//		http.StatusInternalServerError,
	//	)
	//}

	if err != nil {
		return TemplateError{name, err}
	}

	return nil
}

func buildTemplateFuncMap(name string, data interface{}) template.FuncMap {

	if err := templateDataError(data); err != nil {

		return template.FuncMap{
			"yield": func() (template.HTML, error) {
				buf := bytes.NewBuffer(nil)
				err := Templates.ExecuteTemplate(buf, "tables/error", err.Error())
				return template.HTML(buf.String()), err
			},
		}
	}

	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := Templates.ExecuteTemplate(buf, name, data.(TemplateData).data)
			return template.HTML(buf.String()), err
		},
	}

	return funcs
}

type TemplateData struct {
	data interface{}
	err error
}

func templateDataError(data interface{}) error {

	t, ok := data.(TemplateData)
	if ! ok {
		return fmt.Errorf("Unexpected template data type: %v", data)
	}

	return t.err
}

type TemplateError struct {
	TemplateName string
	Err error
}

func (err TemplateError) Error() string {
	return fmt.Sprintf("Error rendering template %s: %s", err.TemplateName, err.Err.Error())
}