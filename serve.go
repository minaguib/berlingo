package berlingo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func do(ai AI, r io.Reader) (response *Response, response_json []byte, err error) {

	game, err := NewGame(ai, r)
	if err != nil {
		return nil, nil, err
	}

	game.Do()

	response_json, err = game.Response.ToJson()
	if err != nil {
		return nil, nil, err
	}

	return response, response_json, nil
}

func ServeHttp(ai AI, port string) {

	fmt.Println("Starting HTTP server on port", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling HTTP request from", r.RemoteAddr)
		_, response_json, err := do(ai, r.Body)
		if err != nil {
			fmt.Printf("Sending errors: %+v\n", err)
			w.Write([]byte("Error"))
		} else {
			fmt.Printf("Sending moves respons\n")
			w.Write(response_json)
		}
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("HTTP Error:", err)
	}
}

func ServeFile(ai AI, filename string) {

	var fh *os.File
	var err error

	if filename == "-" {
		fh = os.Stdin
	} else {
		fh, err = os.Open(filename)
		if err != nil {
			fmt.Println("Error opening", filename, ": ", err)
			return
		}
		defer fh.Close()
	}

	_, response_json, err := do(ai, fh)
	if err != nil {
		fmt.Println("Error processing request:", err)
		return
	}
	os.Stdout.Write(response_json)
}

func Serve(ai AI) {

	port_or_filename := "-"
	if len(os.Args) >= 2 {
		port_or_filename = os.Args[1]
	}

	_, err := strconv.Atoi(port_or_filename)
	if err == nil {
		ServeHttp(ai, port_or_filename)
	} else {
		ServeFile(ai, port_or_filename)
	}
}
