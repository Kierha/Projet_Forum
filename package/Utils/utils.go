package Utils

import (
	"html/template"

	"github.com/gorilla/sessions"
)

var Tmpl *template.Template
var Store = sessions.NewCookieStore([]byte("secret"))
