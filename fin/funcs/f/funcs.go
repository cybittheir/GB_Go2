package f

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func Nl2br(in string) string {
	return strings.Replace(in, "\n", "<br>\n", -1)
}

func CheckPanic(e error) {
	if e != nil {
		panic(e)
	}
}

func Check(e error) {
	if e != nil {
		log.Println("Error:", e)
	}
}

func QueryParseLite(query []byte) (map[string]string, error) {
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

func Contains(s []string, a int) (int, error) {
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
