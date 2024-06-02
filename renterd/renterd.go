package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {

	url := "localhost:8080/api/worker/objects/foo/bar/baz"
	method := "PUT"

	payload := strings.NewReader("<file contents here>")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
