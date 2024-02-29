package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const header = "<!DOCTYPE html>\n<html lang=\"ru\">\n<head>\n<title>%title%</title>\n<script src=\"https://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js\"></script>\n</head>\n"
const footer = "</html>"
const addForm = "<br>\n<a href=\\>go home</a>"
const list = "<br>\n<a href=\\get>get users list</a>"
const friendsForm = "<br>\n<a href=\\friendship>friendship</a>"
const ageForm = "<br>\n<a href=\\age>change user's age</a>"
const deleteForm = "<br>\n<a href=\\delete>delete user</a>"
const page = header + `
<body id=body>
<h4>%title%</h4>
%body%
%menu%
</body>
` + footer

var LastId int

type User struct {
	Id      int
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

func (u *User) ToString() string {

	bp := "User %s with id%d and age of %d y.o. has "

	if len(u.Friends) == 0 {
		return fmt.Sprintf(bp+"no friends.\n", u.Name, u.Id, u.Age)
	} else if len(u.Friends) == 1 {
		return fmt.Sprintf(bp+"only friend %s.\n", u.Name, u.Id, u.Age, strings.Join(u.Friends, ", "))
	} else {
		return fmt.Sprintf(bp+"friends: %s.\n", u.Name, u.Id, u.Age, strings.Join(u.Friends, ", "))
	}
}

func (u *User) Numbered() string {

	bp := "%d. %s (%d)"

	return fmt.Sprintf(bp+"\n", u.Id, u.Name, u.Age)
}

func nl2br(in string) string {
	return strings.Replace(in, "\n", "<br>\n", -1)
}

func useTemplate(title, data, menu string) string {
	context := strings.Replace(page, "%title%", title, -1)
	context = strings.Replace(context, "%body%", data, 1)
	context = strings.Replace(context, "%menu%", menu, 1)
	return context
}

type base struct {
	storage map[int]*User
}

func main() {
	LastId = 0
	mux := http.NewServeMux()
	srv := base{make(map[int]*User)}

	mux.HandleFunc("/", Start)
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	mux.HandleFunc("/delete", srv.DelUser)
	mux.HandleFunc("/friendship", srv.SetRelations)
	mux.HandleFunc("/friendbreak", srv.DelRelations)
	mux.HandleFunc("/age", srv.ChangeAge)
	mux.HandleFunc("/agechange", srv.ChangeAge)

	err := http.ListenAndServe("127.0.0.1:8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting HTTP server: %s\n", err)
		os.Exit(1)
	}

}

func QueryParse(query []byte) (map[string]string, error) {
	q := strings.Split(string(query), "&")
	m := make(map[string]string)

	for _, r := range q {
		t := strings.Split(r, "=")
		m[t[0]] = t[1]
	}
	if len(m) > 0 {
		return m, nil
	}
	return nil, fmt.Errorf("wrong query")
}

func (u *User) QueryParse(query []byte) ([]byte, error) {

	q := strings.Split(string(query), "&")
	var fr []string
	for _, r := range q {
		t := strings.Split(r, "=")
		if t[0] == "friends" {
			u.Friends = fr
		} else if t[0] == "age" && t[1] != "" {
			u.Age, _ = strconv.Atoi(t[1])
		} else if t[0] == "id" && t[1] != "0" {
			u.Id, _ = strconv.Atoi(t[1])
		} else if t[0] == "name" && t[1] != "" {
			u.Name = t[1]
			LastId++
			u.Id = LastId
		}
	}

	return json.Marshal(u)
}

func Start(w http.ResponseWriter, r *http.Request) {

	form := `<form method=POST action=/create>
Username <input type=text id=name name=name size=20><br>
Age <input type=text id=age name=age size=5><br>
<!-- Friends list <input type=text name=friends size=50><br> -->
<input type=submit value="Add User" name=ok>
<button onclick="sendJSON()">Проверить JSON</button>
</form>`
	defer r.Body.Close()
	menu := friendsForm + ageForm + deleteForm + list
	form = strings.Replace(form, "%menu%", menu, 1)
	if r.Method == "GET" {
		w.Write([]byte(useTemplate("Home page - add user", form, menu)))
		return
	}

}

func (b *base) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("1: " + err.Error()))
			fmt.Println("1: " + err.Error())
			return
		}
		defer r.Body.Close()

		var u User
		//TODO: detect Content-type: application/json
		if http.DetectContentType(content) != "" {
			fmt.Println(http.DetectContentType(content))
		}

		if st, err := u.QueryParse(content); err != nil {
			fmt.Println(err.Error())
		} else {
			if err := json.Unmarshal(st, &u); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b.storage[u.Id] = &u

			w.WriteHeader(http.StatusCreated)
			text := "User " + u.Name + " was created and got ID" + strconv.Itoa(u.Id)
			menu := addForm + friendsForm + ageForm + deleteForm + list

			w.Write([]byte(useTemplate("User added", text, menu)))
		}

	}
}

func (b *base) GetAll(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method == "GET" {

		response := ""
		for _, user := range b.storage {
			response += user.ToString()
		}
		w.WriteHeader(http.StatusOK)
		menu := addForm + friendsForm + ageForm + deleteForm
		w.Write([]byte(useTemplate("Users list", nl2br(response), menu)))
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (b *base) DelUser(w http.ResponseWriter, r *http.Request) {
	delForm := `<form method=POST action=/delete>
	User ID <input type=text id=id name=id size=3><br>
	<input type=submit value="Delete" name=ok>
	</form>`
	menu := addForm + friendsForm + ageForm + list

	if r.Method == "GET" {

		response := ""
		for _, user := range b.storage {
			response += user.Numbered()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(useTemplate("Choose user to delete", nl2br(response+delForm), menu)))
		return
	} else if r.Method == "POST" {

		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			fmt.Println(err.Error())
			return
		}

		if st, err := QueryParse(content); err != nil {
			fmt.Println(err.Error())
		} else {
			if delId, err := strconv.Atoi(st["id"]); err != nil {
				w.Write([]byte(useTemplate("Error", fmt.Sprintf("User with ID%d not found. Get list and Try again", delId), menu)))
			} else if b.storage != nil {
				delete(b.storage, delId)
				w.Write([]byte(useTemplate("User was deletes", fmt.Sprintf("User with ID%d was deleted", delId), menu)))

			}
		}
		defer r.Body.Close()
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (b *base) SetRelations(w http.ResponseWriter, r *http.Request) {

	form := `<form method=POST action=/link>
Username <input type=text name=name size=20><br>
Age <input type=text name=age size=5><br>
<!-- Friends list <input type=text name=friends size=50><br> -->
<input type=submit value="Add User" name=ok>
</form>`
	defer r.Body.Close()
	menu := friendsForm + ageForm + deleteForm + list
	form = strings.Replace(form, "%menu%", menu, 1)
	if r.Method == "GET" {
		w.Write([]byte(useTemplate("Home page - add user", form, menu)))
		return
	}
}

func (b *base) DelRelations(w http.ResponseWriter, r *http.Request) {
	return
}

func (b *base) ChangeAge(w http.ResponseWriter, r *http.Request) {
	return
}
