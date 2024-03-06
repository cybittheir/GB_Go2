package base

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

var LastId int

type Base interface {
	AddUser(ctx context.Context, u *User) error
	DElPerson() error
	Info() (*User, error)
	IsExists() (bool, error)
	Count() int
}

type User struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

type Stora struct {
	Storage map[int]*User
}

func (s *Stora) Count() int {
	return len(s.Storage)
}

func (s *Stora) AllFriends(u *User) string {
	res := make([]string, 0)

	for _, v := range s.Storage[u.Id].Friends {
		if vi, err := strconv.Atoi(v); err == nil {
			res = append(res, s.Storage[vi].Name+"<i>(id"+fmt.Sprint(vi)+"; age "+fmt.Sprint(s.Storage[vi].Age)+")</i>")
		}
	}

	return strings.Join(res, ", ")
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
