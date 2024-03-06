package web

import (
	"GB_Study_02/fin/conf"
	"GB_Study_02/fin/funcs/f"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	Home       = "<br>\n<a href=\\>go home / add user</a>"
	List       = "<br>\n<a href=\\get>get users list</a>"
	AgeForm    = "<br>\n<a href=\\age>change user's age</a>"
	DeleteForm = "<br>\n<a href=\\delete>delete user</a>"

	HomeMenu   = "<br>\n<a href=\\>go home / add user</a>"
	ListMenu   = "<br>\n<a href=\\get>get users list</a>"
	AgeMenu    = "<br>\n<a href=\\age>change user's age</a>"
	DeleteMenu = "<br>\n<a href=\\delete>delete user</a>"
)

type User struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

type Stora struct {
	Stora map[int]*User
}

func (s *Stora) Count() int {

	return len(s.Stora)
}

// func (s *Stora) User() *User {
// 	return s.Stora[NewUser().UserId()]
// }

// func (u *User) UserId() int {
// 	return u.Id
// }

func NewUser() *User {
	return &(User{})
}

func NewRec() *Stora {
	return &(Stora{})
}

func (s Stora) NextUser(content []byte) (string, error) {
	var u User
	if err := json.Unmarshal(content, &u); err != nil {
		if st, err := QueryMarshal(string(content)); err != nil {
			log.Println(err.Error())
		} else {
			if err := json.Unmarshal(st, &u); err != nil {
				return "", err
			}

			s.Stora[u.Id] = &u

			return fmt.Sprintf("User %s (%d) was created and got ID%d", u.Name, u.Age, u.Id), nil
		}
	} else {
		LastId++
		u.Id = LastId
		s.Stora[u.Id] = &u

		return fmt.Sprintf("json:id_%d", u.Id), nil
	}

	return "", nil
}

type WebPage struct {
	template string
	title    string
	data     string
	menu     string
	Storage  Stora
}

// var wp WebPage

func (wp *WebPage) Menu(r string) string {

	hm := HomeMenu
	am := AgeMenu
	dm := DeleteMenu
	lm := ListMenu

	l := 0
	if wp.Storage.Stora != nil {
		l = wp.Storage.Count()

	}

	switch l {
	case 0:
		return hm
	default:
		switch r {

		case "/":
			hm = ""
		case "/delete":
			dm = ""
		case "/age":
			am = ""
		case "/get":
			lm = ""
		}
		return hm + am + dm + lm
	}

}

func (wp *WebPage) useTemplate() string {

	context := strings.Replace(string(wp.template), "%title%", wp.title, -1)
	context = strings.Replace(context, "%body%", wp.data, 1)
	context = strings.Replace(context, "%menu%", wp.menu, 1)
	return context
}

var LastId int

func NewHttp(a, p string) (string, error) {
	mux := http.NewServeMux()

	srv := WebPage{template: conf.Template(), Storage: Stora{make(map[int]*User)}}

	mux.HandleFunc("/", srv.Start)
	mux.HandleFunc("/create", srv.Create)
	mux.HandleFunc("/get", srv.GetAll)
	mux.HandleFunc("/delete", srv.DelUser)
	mux.HandleFunc("/friendship/", srv.SetRelations)
	mux.HandleFunc("/age", srv.ChangeAge)

	err := http.ListenAndServe(a+":"+p, mux)

	if errors.Is(err, http.ErrServerClosed) {
		return "", fmt.Errorf("server closed")
	} else if err != nil {
		return "", fmt.Errorf("error starting HTTP server: %s", err)
	}
	return "HTTP-server started", nil

}

func (wp *WebPage) Start(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		form := `<form method=POST action=/create>
		Username <input type=text id=name name=name size=20><br>
		Age <input type=text id=age name=age size=5><br>
		<input type=submit value="Add User" name=ok>
		</form>`
		defer r.Body.Close()

		menu := wp.Menu(r.RequestURI)
		form = strings.Replace(form, "%menu%", menu, 1)
		wp.title = "Home page - add user"
		wp.menu = menu
		wp.data = form

		w.Write([]byte((wp).useTemplate()))
		return
	}

}

func (wp *WebPage) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("1: " + err.Error()))
			fmt.Println("1: " + err.Error())
			return
		}
		defer r.Body.Close()

		if text, err := wp.Storage.NextUser(content); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
			if !strings.Contains(text, "json:id_") {
				menu := wp.Menu(r.RequestURI)

				wp.title = "New user added"
				wp.menu = menu
				wp.data = text
				w.Write([]byte((wp).useTemplate()))
			} else {
				jp := strings.Split(text, "id_")
				if len(jp) > 1 {
					j, _ := strconv.Atoi(jp[1])
					if wp.Storage.Stora[j] != nil {
						if js, err := json.Marshal(wp.Storage.Stora[j]); err != nil {
							w.Write([]byte("Error. Unable to marshal JSON"))
						} else {
							w.Write(js)
						}
					}

				}
			}

		}

	}
}

func (wp *WebPage) GetAll(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method == "GET" {

		//		q := "SELECT * FROM users"

		// q:="select id1 as l,users.user_name,users.age from(
		// select friends.id1 from friends
		// where id2=%i
		// union
		// select friends.id2 from friends
		// where id1=%i
		// )s
		// left join users on l=id"

		response := ""
		for _, u := range wp.Storage.Stora {
			response += u.ToString() + wp.Storage.AllFriends(u) + ".\n"
		}
		w.WriteHeader(http.StatusOK)

		menu := wp.Menu(r.RequestURI)

		wp.title = "Users list"
		wp.menu = menu
		wp.data = f.Nl2br(response)

		w.Write([]byte((wp).useTemplate()))
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (wp *WebPage) DelUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		response := ""
		for _, u := range wp.Storage.Stora {
			response += u.Numbered()
		}
		w.WriteHeader(http.StatusOK)

		delForm := `<form method=POST action=/delete>
		User ID <input type=text id=id name=id size=3><br>
		<input type=submit value="Delete" name=ok>
		</form>`

		menu := wp.Menu(r.RequestURI)

		wp.title = "Choose user to delete"
		wp.menu = menu
		wp.data = f.Nl2br(response + delForm)

		w.Write([]byte((wp).useTemplate()))

		return
	} else if r.Method == "POST" {

		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()

		var u User

		if st, err := QueryMarshal(string(content)); err != nil {
			fmt.Println(err.Error())
		} else {
			if err := json.Unmarshal(st, &u); err != nil {

				menu := wp.Menu(r.RequestURI)
				wp.title = "Error"
				wp.menu = menu
				wp.data = f.Nl2br(fmt.Sprintf("User with ID%d not found. Get list and Try again", u.Id))

				w.Write([]byte((wp).useTemplate()))

			} else if wp.Storage.Stora != nil {
				delete(wp.Storage.Stora, u.Id)
				wp.Storage.RemoveFriend(u.Id)

				menu := wp.Menu(r.RequestURI)

				wp.title = "User was deleted"
				wp.menu = menu
				wp.data = f.Nl2br(fmt.Sprintf("User with ID%d was deleted", u.Id))
				w.Write([]byte((wp).useTemplate()))

			}
		}
		return
	} else if r.Method == "DELETE" {
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			fmt.Println(err.Error())
			return
		}
		defer r.Body.Close()
		var u User

		if err := json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		//			var u User
		w.WriteHeader(http.StatusOK)
		name := wp.Storage.Stora[u.Id].Name
		delete(wp.Storage.Stora, u.Id)
		wp.Storage.RemoveFriend(u.Id)

		w.Write([]byte(fmt.Sprintf("removed: %s", name)))

	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusBadRequest)

}

func (wp *WebPage) ChangeAge(w http.ResponseWriter, r *http.Request) {

	menu := wp.Menu(r.RequestURI)

	if r.Method == "GET" {

		response := ""
		for _, u := range wp.Storage.Stora {
			response += u.Numbered()
		}
		w.WriteHeader(http.StatusOK)

		aForm := `<form method=POST action=/age>
		User ID <input type=text id=id name=id size=3><br>
		Correct Age to <input type=text id=age name=age size=3><br>
		<input type=submit value="Change" name=ok>
		</form>`

		wp.title = "Choose user to change his Age"
		wp.menu = menu
		wp.data = f.Nl2br(response + aForm)
		w.Write([]byte((wp).useTemplate()))

		return
	}

	if r.Method == "POST" {

		var u User
		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
			return
		}
		defer r.Body.Close()

		if st, err := QueryMarshal(string(content)); err != nil {
			log.Println(err.Error())
		} else {
			if err := json.Unmarshal(st, &u); err != nil {
				wp.title = "Error"
				wp.menu = menu
				wp.data = f.Nl2br("Wrong query. Try again")
				w.Write([]byte((wp).useTemplate()))
				return
			}

			if u.Id < 0 {
				wp.title = "Error"
				wp.menu = menu
				wp.data = f.Nl2br(fmt.Sprintf("wrong user ID%d. Get list and Try again", u.Id))
				w.Write([]byte((wp).useTemplate()))

				return
			} else if u.Age < 1 {
				wp.title = "Error"
				wp.menu = menu
				wp.data = f.Nl2br(fmt.Sprintf("wrong User Age%d. Get list and Try again", u.Age))
				w.Write([]byte((wp).useTemplate()))
				return
			} else if wp.Storage.Stora[u.Id] != nil {
				wp.Storage.Stora[u.Id].Age = u.Age
				wp.title = "User's Age was changed"
				wp.menu = menu
				wp.data = f.Nl2br(fmt.Sprintf("New User's Age (ID%d) is %d now", u.Id, u.Age))
				w.Write([]byte((wp).useTemplate()))

			}
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (wp *WebPage) SetRelations(w http.ResponseWriter, r *http.Request) {

	menu := wp.Menu(r.RequestURI)

	qs := strings.Split(r.RequestURI, "/id")
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()

	if len(qs) < 2 {
		wp.title = "Error"
		wp.menu = menu
		wp.data = f.Nl2br("First choose user in the list for set relations")
		w.Write([]byte((wp).useTemplate()))
		return
	}

	if r.Method == "GET" {
		response := ""

		if q, err := strconv.Atoi(qs[1]); err != nil {
			wp.title = "Error"
			wp.menu = menu
			wp.data = f.Nl2br("Wrong ID format")
			w.Write([]byte((wp).useTemplate()))
			return
		} else {
			for i, u := range wp.Storage.Stora {
				if i != q {
					response += u.Checked(slices.Contains(wp.Storage.Stora[q].Friends, fmt.Sprint(i)))
				}
			}
			form := fmt.Sprintf("<form method=POST action=%s>%s<input type=submit value='Link Friends' name=ok>\n</form>", r.RequestURI, response)

			wp.title = "Choose users to set friendship with " + wp.Storage.Stora[q].Name
			wp.menu = menu
			wp.data = f.Nl2br(form)
			w.Write([]byte((wp).useTemplate()))
		}
		return
	}

	if r.Method == "POST" {

		if q, err := strconv.Atoi(qs[1]); err != nil {
			wp.title = "Error"
			wp.menu = menu
			wp.data = f.Nl2br("User not found")
			w.Write([]byte((wp).useTemplate()))
			return
		} else {
			content, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				fmt.Println(err.Error())
				return
			}

			defer r.Body.Close()
			if st, err := f.QueryParseLite(content); err != nil {
				fmt.Println(err.Error())
				wp.title = "Error"
				wp.menu = menu
				wp.data = f.Nl2br("User not found")
				w.Write([]byte((wp).useTemplate()))
			} else {
				friends := make([]string, 0)
				for i, v := range st {
					if v == "on" {
						fn := strings.Split(i, "_")
						friends = append(friends, fn[1])
						ff, _ := strconv.Atoi(fn[1])
						if _, err := f.Contains(wp.Storage.Stora[ff].Friends, q); err != nil {
							wp.Storage.Stora[ff].Friends = append(wp.Storage.Stora[ff].Friends, strconv.Itoa(q))
						}

					}
				}
				oldFriends := wp.Storage.Stora[q].Friends
				if len(oldFriends) > 0 {
					for _, va := range oldFriends {
						if v, err := strconv.Atoi(va); err != nil {
							log.Println("Err. convert:", err)
						} else {
							if i, err := f.Contains(wp.Storage.Stora[v].Friends, q); err != nil {
								//								log.Println("Err. nothing to do:", err)
							} else if _, err := f.Contains(friends, v); err != nil {
								wp.Storage.Stora[v].Friends[i] = wp.Storage.Stora[v].Friends[len(wp.Storage.Stora[v].Friends)-1]
								wp.Storage.Stora[v].Friends[len(wp.Storage.Stora[v].Friends)-1] = ""
								wp.Storage.Stora[v].Friends = wp.Storage.Stora[v].Friends[:len(wp.Storage.Stora[v].Friends)-1]
							}
						}
					}
				}

				wp.Storage.Stora[q].Friends = friends
				wp.title = "Relationship set"
				wp.menu = menu
				wp.data = f.Nl2br(wp.Storage.Stora[q].Name + " friendship updated")
				w.Write([]byte((wp).useTemplate()))

			}
			return
		}
	}
}

func QueryMarshal(query string) ([]byte, error) {

	q := strings.Split(query, "&")

	var u User

	for _, r := range q {
		t := strings.Split(r, "=")
		if t[0] == "friends" {
			u.Friends = make([]string, 0)
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

func (u *User) ToString() string {

	b := "User <b>%s</b> with <u>id%d</u> and age of %d y.o. has "

	fUrl := fmt.Sprintf("<a href=/friendship/id%d>", u.Id)

	if len(u.Friends) == 0 {
		return fmt.Sprintf(b+"no %sfriends</a>", u.Name, u.Id, u.Age, fUrl)
	} else if len(u.Friends) == 1 {
		return fmt.Sprintf(b+"only %sfriend</a> ", u.Name, u.Id, u.Age, fUrl)
	} else {
		return fmt.Sprintf(b+"%d %sfriends</a>: ", u.Name, u.Id, u.Age, len(u.Friends), fUrl)
	}
}

func (s *Stora) AllFriends(u *User) string {
	res := make([]string, 0)

	for _, v := range s.Stora[u.Id].Friends {
		if i, err := strconv.Atoi(v); err == nil {
			if s.Stora[i] != nil {
				res = append(res, fmt.Sprintf("%s<i>(id%d; age %d)</i>", s.Stora[i].Name, i, s.Stora[i].Age))
			}
		}
	}

	return strings.Join(res, ", ")
}

func (u *User) Numbered() string {

	bp := "%d. %s (%d)"

	return fmt.Sprintf(bp+"\n", u.Id, u.Name, u.Age)
}

func (u *User) Checked(c bool) string {

	b := ""
	if c {
		b = "%d.<input type=checkbox name=id_%d id=id_%d checked><label for=id_%d>%s (%d)</label>"

	} else {
		b = "%d.<input type=checkbox name=id_%d id=id_%d><label for=id_%d>%s (%d)</label>"
	}
	return fmt.Sprintf(b+"\n", u.Id, u.Id, u.Id, u.Id, u.Name, u.Age)
}

func (s *Stora) RemoveFriend(uid int) {
	for i, v := range s.Stora {
		if i != uid {
			if ind, err := f.Contains(v.Friends, uid); err != nil {
				log.Println("Err. nothing to do:", err)
			} else {
				v.Friends[ind] = v.Friends[len(v.Friends)-1]
				v.Friends[len(v.Friends)-1] = ""
				v.Friends = v.Friends[:len(v.Friends)-1]
			}
		}
	}
}
