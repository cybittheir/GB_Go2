package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	header     = "<!DOCTYPE html>\n<html lang=\"ru\">\n<head>\n<title>%title%</title>\n<meta charset=\"utf-8\">\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">\n</head>\n"
	footer     = "</html>"
	Home       = "<br>\n<a href=\\>go home / add user</a>"
	List       = "<br>\n<a href=\\get>get users list</a>"
	AgeForm    = "<br>\n<a href=\\age>change user's age</a>"
	DeleteForm = "<br>\n<a href=\\delete>delete user</a>"
	Page       = header + `
<body id=body>
<h4>%title%</h4>
%body%
%menu%
</body>
` + footer
)

type WebPage struct {
	title  string
	header string
	body   string
	footer string
	data   string
	menu   string
}

func (w *WebPage) useTemplate() string {

	context := strings.Replace(Page, "%title%", w.title, -1)
	context = strings.Replace(context, "%body%", w.data, 1)
	context = strings.Replace(context, "%menu%", w.menu, 1)
	return context
}

func NewHttp() {
	mux := http.NewServeMux()
	srv := base{make(map[int]*User)}

	mux.HandleFunc("/", srv.Start)
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	mux.HandleFunc("/delete", srv.DelUser)
	mux.HandleFunc("/friendship/", srv.SetRelations)
	mux.HandleFunc("/age", srv.ChangeAge)

	err := http.ListenAndServe(httpServ, mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting HTTP server: %s\n", err)
		os.Exit(1)
	}

}
