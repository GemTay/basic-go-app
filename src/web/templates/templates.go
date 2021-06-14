package templates

import (
	"html/template"
	"log"
	"net/http"
)

var tpls = template.Must(template.ParseGlob("./web/templates/*"))

func Render(w http.ResponseWriter, filename string, data interface{}) {

	err := tpls.ExecuteTemplate(w, filename, data)
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}
