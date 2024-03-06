package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const proxyAddr string = "localhost:8080"

type User struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

var (
	counter            int    = 0
	firstInstanceHost  string = "http://localhost:8081"
	secondInstanceHost string = "http://localhost:8082"
)

func main() {
	http.HandleFunc("/", handleProxy)
	log.Fatalln(http.ListenAndServe(proxyAddr, nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {

	textBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var u User

	if err := json.Unmarshal(textBytes, &u); err != nil {
		fmt.Println("0>", u)
	} else {
		fmt.Println("1>", string(textBytes))
	}
	text, err := json.Marshal(&u)

	if err != nil {
		log.Fatalln(err)
	}

	//	fmt.Println("2>", string(text))

	client := http.Client{Timeout: 2 * time.Second}

	if counter == 0 {
		//		resp, err := http.Post(firstInstanceHost, "application/json; charset=utf-8", bytes.NewBuffer(text))
		resp, err := http.NewRequest(r.Method, firstInstanceHost+r.RequestURI, bytes.NewBuffer(text))
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")

		if err != nil {
			log.Fatalln(err)
		}
		//		fmt.Println(firstInstanceHost, "<<", string(textBytes))
		counter++
		res, err := client.Do(resp)

		defer resp.Body.Close()

		resBytes, err := io.ReadAll(res.Body)

		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(firstInstanceHost+r.RequestURI, "<<", string(textBytes), ";", r.Method, ">>", string(resBytes))

		return
	}

	//	resp, err := http.Post(secondInstanceHost, "application/json; charset=utf-8", bytes.NewBuffer(text))
	resp, err := http.NewRequest(r.Method, secondInstanceHost+r.RequestURI, bytes.NewBuffer(text))
	resp.Header.Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		log.Fatalln(err)
	}

	counter--
	res, err := client.Do(resp)

	defer resp.Body.Close()

	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(secondInstanceHost+r.RequestURI, "<<", string(textBytes), ";", r.Method, ">>", string(resBytes))

}
