package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	Address  string `json:"addr"`
	Port     string `json:"port"`
	DBFile   string `json:"dbfile"`
	JSFile   string `json:"jsfile"`
	Template string `json:"tplfile"`
}

var Cfg Config

func (c *Config) Conf() *Config {

	fn := "conf.json"
	_, err := os.Stat(fn)
	if errors.Is(err, os.ErrNotExist) {
		log.Println(err)
		c.Address = "127.0.0.1"
		c.Port = "8080"
		c.DBFile = "sqlite.db"
		c.JSFile = "storage.json"
		c.Template = "index.tpl"

		var b []byte

		b, _ = json.Marshal(c)

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
		c.Address = "127.0.0.1"
		c.Port = "8080"
		c.DBFile = "sqlite.db"
		c.JSFile = "storage.json"
		c.Template = "index.tpl"
		return c
	}
	fmt.Println("Successfully opened: conf.json")

	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	var cf Config

	json.Unmarshal([]byte(byteValue), &cf)
	fmt.Println("Using config:")
	fmt.Println("http://", cf.Address, ":", cf.Port)
	fmt.Println("SQLite: ", cf.DBFile)
	fmt.Println("Storage: ", cf.JSFile)
	fmt.Println("template: ", cf.Template)
	fmt.Println("==================")
	Cfg.Template = cf.Template
	Cfg.JSFile = cf.JSFile
	return &cf

}

func Template() string {

	tf := Cfg.Template
	templateFile, err := os.Open(tf)
	if err != nil {
		header := "<!DOCTYPE html>\n<html lang=\"ru\">\n<head>\n<title>%title%</title>\n<meta charset=\"utf-8\">\n<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">\n</head>\n"
		footer := "</html>"
		page := header + `
	<body id=body>
	<h4>%title%</h4>
	%body%
	%menu%
	</body>
	` + footer
		return page
	}
	defer templateFile.Close()
	byteValue, _ := io.ReadAll(templateFile)
	return string(byteValue)
}

func JSStor() string {

	return Cfg.JSFile
}
