package berlingo

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// Callback used to process an incoming HTTP request
func serveHttpRequest(ai AI, w http.ResponseWriter, r *http.Request) {

	log.Printf("HTTP: [%v] Processing %v %v", r.RemoteAddr, r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")

	var input io.Reader
	content_type := r.Header.Get("Content-Type")
	switch {
	case r.Method == "POST" && content_type == "application/json":
		input = r.Body
	case r.Method == "POST" && content_type == "application/x-www-form-urlencoded":
		// Detect & work-around bug https://github.com/thirdside/berlin-ai/issues/4
		r.ParseForm()
		j := `{
				"action": "` + r.Form["action"][0] + `",
				"infos": ` + r.Form["infos"][0] + `,
				"map": ` + r.Form["map"][0] + `,
				"state": ` + r.Form["state"][0] + `
			}`
		input = strings.NewReader(j)
	default:
		log.Printf("HTTP: Replying with error")
		w.Write([]byte(`{"error": "Invalid request"}`))
		return
	}

	_, response_json, err := do(ai, input)
	if err != nil {
		log.Printf("HTTP: Responding with error: %+v\n", err)
		w.Write([]byte("Error"))
	} else {
		log.Printf("HTTP: Responding with moves\n")
		w.Write(response_json)
	}

}

// ServeHttp serves the given AI over HTTP on the given port
func ServeHttp(ai AI, port string) {

	log.Println("Starting HTTP server on port", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHttpRequest(ai, w, r)
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("HTTP Serving Error:", err)
	}

}

// ServeHttp serves the given AI a single time
// Request is read from the given filename
// filename may be supplied as "-" to indicate STDIN
func ServeFile(ai AI, filename string) {

	var fh *os.File
	var err error

	if filename == "-" {
		fh = os.Stdin
	} else {
		fh, err = os.Open(filename)
		if err != nil {
			log.Println("Error opening", filename, ": ", err)
			return
		}
		defer fh.Close()
	}

	_, response_json, err := do(ai, fh)
	if err != nil {
		log.Println("Error processing request:", err)
		return
	}
	os.Stdout.Write(response_json)
}

// Serve will inspect the CLI arguments and automatically call either ServeHttp or ServeFile
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
