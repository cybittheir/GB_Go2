package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	header     = "<!DOCTYPE html>\n<html lang=\"ru\">\n<head>\n<title>%title%</title>\n<meta charset=\"utf-8\">\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">\n</head>\n"
	footer     = "</html>"
	home       = "<br>\n<a href=\\>go home / add user</a>"
	list       = "<br>\n<a href=\\get>get users list</a>"
	ageForm    = "<br>\n<a href=\\age>change user's age</a>"
	deleteForm = "<br>\n<a href=\\delete>delete user</a>"
	page       = header + `
<body id=body>
<h4>%title%</h4>
%body%
%menu%
</body>
` + footer
)

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
