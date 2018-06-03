package ginamber

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/eknkc/amber"
	"github.com/gin-gonic/gin/render"
)

// ErrTemplateNotFound is returned when the template was not found in the directory
var ErrTemplateNotFound = errors.New("Error: template not found")

// AmberHTMLRender implements the render methods used by gin-gonic for
// Amber template engine
type AmberHTMLRender struct {
	DirOptions amber.DirOptions
	Options    amber.Options
	Dir        string

	templateMap map[string]*template.Template
	err         error
}

// NewDefaultOptions returns an AmberHTMLRender with all default options
func NewDefaultOptions() AmberHTMLRender {
	return AmberHTMLRender{
		DirOptions: amber.DefaultDirOptions,
		Options:    amber.DefaultOptions,
		Dir:        "templates",
	}
}

// Instance gets an instance of the amber renderer
func (a AmberHTMLRender) Instance(tplt string, args interface{}) render.Render {
	ret := AmberRender{}
	if a.templateMap == nil {
		a.Compile()
	}
	if a.templateMap != nil {
		t, ok := a.templateMap[tplt]
		if ok {
			ret.template = t
		}
	}
	ret.args = args
	return ret
}

// Compile compiles the templates in the configured directory
func (a *AmberHTMLRender) Compile() error {
	a.templateMap, a.err = amber.CompileDir(a.Dir, a.DirOptions, a.Options)
	return a.err
}

// AmberRender writes the response to the ResponseWriter
type AmberRender struct {
	template *template.Template
	args     interface{}
}

// Render writes the response to the ResponseWriter
func (a AmberRender) Render(w http.ResponseWriter) error {
	if a.template == nil {
		return ErrTemplateNotFound
	}
	a.template.Execute(w, a.args)
	return nil
}

// WriteContentType sets the response's content-type
func (a AmberRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
}
