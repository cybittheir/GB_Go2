package main

import (
	"GB_Study_02/fin/conf"
	"GB_Study_02/fin/storage/sqlite"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// var LastId int

type User struct {
	Id      int
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

func (u *User) ToString() string {

	bp := "User <b>%s</b> with <u>id%d</u> and age of %d y.o. has "

	fUrl := fmt.Sprintf("<a href=/friendship/id%d>", u.Id)

	if len(u.Friends) == 0 {
		return fmt.Sprintf(bp+"no %sfriends</a>", u.Name, u.Id, u.Age, fUrl)
	} else if len(u.Friends) == 1 {
		return fmt.Sprintf(bp+"only %sfriend</a> ", u.Name, u.Id, u.Age, fUrl)
	} else {
		return fmt.Sprintf(bp+"%d %sfriends</a>: ", u.Name, u.Id, u.Age, len(u.Friends), fUrl)
	}
}

func (b *base) AllFriends(u *User) string {
	res := make([]string, 0)

	for _, v := range b.storage[u.Id].Friends {
		if vi, err := strconv.Atoi(v); err == nil {
			res = append(res, b.storage[vi].Name+"<i>(id"+fmt.Sprint(vi)+"; age "+fmt.Sprint(b.storage[vi].Age)+")</i>")
		}
	}

	return strings.Join(res, ", ")
}

func (u *User) Numbered() string {

	bp := "%d. %s (%d)"

	return fmt.Sprintf(bp+"\n", u.Id, u.Name, u.Age)
}

func (u *User) Checked(c bool) string {

	bp := ""
	if c {
		bp = "%d.<input type=checkbox name=id_%d id=id_%d checked><label for=id_%d>%s (%d)</label>"

	} else {
		bp = "%d.<input type=checkbox name=id_%d id=id_%d><label for=id_%d>%s (%d)</label>"
	}
	return fmt.Sprintf(bp+"\n", u.Id, u.Id, u.Id, u.Id, u.Name, u.Age)
}

type base struct {
	storage map[int]*User
}

func main() {

	//	LastId := 0

	var cf conf.Config

	c := (*conf.Config).Conf(&cf)

	s, err := sqlite.New(c.DBFile)

	if err != nil {
		log.Fatalf("can't connect to DB: %s", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatalf("can't init DB: %s", err)
	}

	mux := http.NewServeMux()
	// srv := base{make(map[int]*User)}

	// mux.HandleFunc("/", srv.Start)
	// mux.HandleFunc("/create", srv.Create)
	// mux.HandleFunc("/get", srv.GetAll)
	// mux.HandleFunc("/delete", srv.DelUser)
	// mux.HandleFunc("/friendship/", srv.SetRelations)
	// mux.HandleFunc("/age", srv.ChangeAge)

	if err := http.ListenAndServe(c.Address+":"+c.Port, mux); errors.Is(err, http.ErrServerClosed) {
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

/*
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

func (b *base) Start(w http.ResponseWriter, r *http.Request) {

	form := `<form method=POST action=/create>
Username <input type=text id=name name=name size=20><br>
Age <input type=text id=age name=age size=5><br>
<input type=submit value="Add User" name=ok>
</form>`
	defer r.Body.Close()
	menu := ""
	switch len(b.storage) {
	case 0:
		menu = ""
	default:
		menu = ageForm + deleteForm + list
	}
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
			log.Println(http.DetectContentType(content))
		}

		if st, err := u.QueryParse(content); err != nil {
			log.Println(err.Error())
		} else {
			if err := json.Unmarshal(st, &u); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			b.storage[u.Id] = &u

			w.WriteHeader(http.StatusCreated)
			text := "User " + u.Name + " was created and got ID" + strconv.Itoa(u.Id)

			menu := ""
			switch len(b.storage) {
			case 0:
				menu = home
			default:
				menu = home + ageForm + deleteForm + list
			}

			w.Write([]byte(useTemplate("User added", text, menu)))
		}

	}
}

func (b *base) GetAll(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	if r.Method == "GET" {

		response := ""
		for _, user := range b.storage {
			response += user.ToString() + b.AllFriends(user) + ".\n"
		}
		w.WriteHeader(http.StatusOK)
		menu := ""
		switch len(b.storage) {
		case 0:
			menu = home
		default:
			menu = home + ageForm + deleteForm
		}

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
	menu := ""

	switch len(b.storage) {
	case 0:
		menu = home
	default:
		menu = home + ageForm + list
	}

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
		defer r.Body.Close()

		if st, err := QueryParse(content); err != nil {
			fmt.Println(err.Error())
		} else {
			if delId, err := strconv.Atoi(st["id"]); err != nil {
				w.Write([]byte(useTemplate("Error", fmt.Sprintf("User with ID%d not found. Get list and Try again", delId), menu)))
			} else if b.storage != nil {
				delete(b.storage, delId)
				b.RemoveFriend(delId)
				switch len(b.storage) {
				case 0:
					menu = home
				default:
					menu = home + ageForm + list
				}
				w.Write([]byte(useTemplate("User was deletes", fmt.Sprintf("User with ID%d was deleted", delId), menu)))

			}
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (b *base) ChangeAge(w http.ResponseWriter, r *http.Request) {
	aForm := `<form method=POST action=/age>
	User ID <input type=text id=id name=id size=3><br>
	Correct Age to <input type=text id=age name=age size=3><br>
	<input type=submit value="Change" name=ok>
	</form>`

	menu := ""
	switch len(b.storage) {
	case 0:
		menu = home
	default:
		menu = home + deleteForm + list
	}

	if r.Method == "GET" {

		response := ""
		for _, user := range b.storage {
			response += user.Numbered()
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(useTemplate("Choose user to change his Age", nl2br(response+aForm), menu)))
		return
	}

	if r.Method == "POST" {

		content, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
			return
		}
		defer r.Body.Close()

		if st, err := QueryParse(content); err != nil {
			log.Println(err.Error())
		} else {
			menu = home + ageForm + list
			if uId, err := strconv.Atoi(st["id"]); err != nil {
				w.Write([]byte(useTemplate("Error", fmt.Sprintf("wrong user ID%d. Get list and Try again", uId), menu)))
				return
			} else if uAge, err := strconv.Atoi(st["age"]); err != nil {
				w.Write([]byte(useTemplate("Error", fmt.Sprintf("wrong User Age%d. Get list and Try again", uAge), menu)))
				return
			} else if b.storage[uId] != nil {
				b.storage[uId].Age = uAge
				w.Write([]byte(useTemplate("User's Age was changed", fmt.Sprintf("New User's Age (ID%d) is %d now", uId, uAge), menu)))

			}
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (b *base) SetRelations(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	menu := ""
	switch len(b.storage) {

	case 0:
		menu = home
	default:
		menu = home + ageForm + deleteForm + list
	}

	qs := strings.Split(r.RequestURI, "/id")
	if len(qs) < 1 {
		return
	}

	if r.Method == "GET" {
		response := ""
		w.WriteHeader(http.StatusOK)
		if q, err := strconv.Atoi(qs[1]); err != nil {
			w.Write([]byte(useTemplate("Error", "User not found", menu)))
			return
		} else {
			for ind, user := range b.storage {
				if ind != q {
					response += user.Checked(slices.Contains(b.storage[q].Friends, strconv.Itoa(ind)))
				}
			}
			form := "<form method=POST action=" + r.RequestURI + ">" + response +
				"<input type=submit value='Link Friends' name=ok>\n</form>"

			w.Write([]byte(useTemplate("Choose users to set friendship with "+b.storage[q].Name, nl2br(form), menu)))
		}
		return
	}

	if r.Method == "POST" {

		if q, err := strconv.Atoi(qs[1]); err != nil {
			w.Write([]byte(useTemplate("Error", "User not found", menu)))
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
			if st, err := QueryParse(content); err != nil {
				fmt.Println(err.Error())
				w.Write([]byte(useTemplate("Error", "User not found", menu)))
			} else {
				menu = home + ageForm + list
				friends := make([]string, 0)
				for i, v := range st {
					if v == "on" {
						fn := strings.Split(i, "_")
						friends = append(friends, fn[1])
						ff, _ := strconv.Atoi(fn[1])
						if _, err := contains(b.storage[ff].Friends, q); err != nil {
							b.storage[ff].Friends = append(b.storage[ff].Friends, strconv.Itoa(q))
						}

					}
				}
				oldFriends := b.storage[q].Friends
				if len(oldFriends) > 0 {
					for _, va := range oldFriends {
						if v, err := strconv.Atoi(va); err != nil {
							log.Println("Err. convert:", err)
						} else {
							if i, err := contains(b.storage[v].Friends, q); err != nil {
								log.Println("Err. nothing to do:", err)
							} else if _, err := contains(friends, v); err != nil {
								b.storage[v].Friends[i] = b.storage[v].Friends[len(b.storage[v].Friends)-1]
								b.storage[v].Friends[len(b.storage[v].Friends)-1] = ""
								b.storage[v].Friends = b.storage[v].Friends[:len(b.storage[v].Friends)-1]
							}
						}
					}
				}

				b.storage[q].Friends = friends

				w.Write([]byte(useTemplate("Relationship set", b.storage[q].Name+" friendship updated", menu)))
			}
			return
		}
	}
}

func (b *base) RemoveFriend(uid int) {
	for i, v := range b.storage {
		if i != uid {
			if ind, err := contains(v.Friends, uid); err != nil {
				log.Println("Err. nothing to do:", err)
			} else {
				v.Friends[ind] = v.Friends[len(v.Friends)-1]
				v.Friends[len(v.Friends)-1] = ""
				v.Friends = v.Friends[:len(v.Friends)-1]
			}
		}
	}
}

func contains(s []string, a int) (int, error) {
	if len(s) > 0 {
		for i, v := range s {
			if vi, err := strconv.Atoi(v); err != nil {
				return -1, fmt.Errorf("element not found:%s", v)
			} else {
				if vi == a {
					return i, nil
				}
			}
		}
		return -1, fmt.Errorf("element not found")
	} else {
		return -1, fmt.Errorf("list is empty")
	}
}
*/
