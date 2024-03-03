package memory

type User struct {
	Id      int
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

type base struct {
	storage map[int]*User
}
