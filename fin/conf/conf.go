package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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

type Config struct {
	Address  string `json:"addr"`
	Port     string `json:"port"`
	DBFile   string `json:"dbfile"`
	Template string `json:"tplfile"`
}

func (d *Config) Conf() *Config {

	fn := "conf.json"
	_, err := os.Stat(fn)
	if errors.Is(err, os.ErrNotExist) {
		log.Println(err)
		d.Address = "127.0.0.1"
		d.Port = "8080"
		d.DBFile = "sqlite.db"
		d.Template = "index.tpl"

		var b []byte

		b, _ = json.Marshal(d)

		if err := os.WriteFile(fn, b, 0644); err != nil {
			fmt.Println("Unable to save", err)
			return nil
		}
		fmt.Println("Default config saved to: conf.json")

	}

	configFile, err := os.Open(fn)
	if err != nil {
		fmt.Println("Can not open conf.json:", err)
		fmt.Println("Using default: 127.0.0.1:8080; sqlite.db; index.tpl")
		d.Address = "127.0.0.1"
		d.Port = "8080"
		d.DBFile = "sqlite.db"
		d.Template = "index.tpl"
		return d
	}
	fmt.Println("Successfully opened: conf.json")

	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	var c Config

	json.Unmarshal([]byte(byteValue), &c)
	fmt.Println("Using config:")
	fmt.Println("http://", c.Address, ":", c.Port)
	fmt.Println("SQLite: ", c.DBFile)
	fmt.Println("template: ", c.Template)
	fmt.Println("==================")

	return &c

}
