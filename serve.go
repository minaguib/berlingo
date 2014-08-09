// Package berlingo is a Go framework for writing AIs for berlin-ai.com
package berlingo

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func do(ai AI, r io.Reader) (response *Response, responseJSON []byte, err error) {

	game, err := NewGame(ai, r)
	if err != nil {
		return nil, nil, err
	}

	game.DoAction()

	responseJSON, err = game.Response.ToJSON()
	if err != nil {
		return nil, nil, err
	}

	return response, responseJSON, nil
}

// Callback used to process an incoming HTTP request
func serveHTTPRequest(ai AI, w http.ResponseWriter, r *http.Request) {

	log.Printf("HTTP: [%v] Processing %v %v", r.RemoteAddr, r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")

	var input io.Reader
	contentType := r.Header.Get("Content-Type")
	switch {
	case r.Method == "POST" && contentType == "application/json":
		input = r.Body
	case r.Method == "POST" && contentType == "application/x-www-form-urlencoded":
		// Detect & work-around bug https://github.com/thirdside/berlin-ai/issues/4
		r.ParseForm()
		j := `{
				"action": "` + r.Form.Get("action") + `",
				"infos": ` + r.Form.Get("infos") + `,
				"map": ` + r.Form.Get("map") + `,
				"state": ` + r.Form.Get("state") + `
			}`
		input = strings.NewReader(j)
	default:
		log.Printf("HTTP: Replying with error: Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request"}`))
		return
	}

	_, responseJSON, err := do(ai, input)
	if err != nil {
		log.Printf("HTTP: Responding with error: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	} else {
		log.Printf("HTTP: Responding with moves\n")
		w.Write(responseJSON)
	}

}

// InitAppEngine allows usage on Google AppEngine
func InitAppEngine(ai AI) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHTTPRequest(ai, w, r)
	})
}

// ServeHTTP serves the given AI over HTTP on the given port
func ServeHTTP(ai AI, port string) {

	log.Println("Starting HTTP server on port", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHTTPRequest(ai, w, r)
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("HTTP Serving Error:", err)
	}

}

// ServeFile serves the given AI a single time
// JSON request is read from the given filename
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

	_, responseJSON, err := do(ai, fh)
	if err != nil {
		log.Println("Error processing request:", err)
		return
	}
	os.Stdout.Write(responseJSON)
}

// Serve will inspect the CLI arguments and automatically call either ServeHTTP or ServeFile
func Serve(ai AI) {

	portOrFilename := "-"
	if len(os.Args) >= 2 {
		portOrFilename = os.Args[1]
	}

	_, err := strconv.Atoi(portOrFilename)
	if err == nil {
		ServeHTTP(ai, portOrFilename)
	} else {
		ServeFile(ai, portOrFilename)
	}
}
