package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	port  = flag.String("port", "8000", "port to listen on")
	name  = flag.String("name", "gocho", "server name")
	debug = flag.Bool("debug", false, "enable debug logging")
)

type Request struct {
	Datetime string
	Server   string
	Hostname string
	Method   string
	Path     string
	Query    map[string][]string
	Headers  map[string][]string
	Body     string
	Duration int64
	Status   int
}

func echo(w http.ResponseWriter, req *http.Request) {
	started := time.Now()
	if *debug {
		log.Print("[" + started.UTC().String() + "] " + req.Method + " " + req.URL.String())
	}

	r := Request{
		Datetime: started.UTC().String(),
		Server:   *name,
		Hostname: req.Host,
		Method:   req.Method,
		Path:     strings.Split(req.URL.String(), "?")[0],
		Query:    req.URL.Query(),
		Headers:  req.Header,
		Body:     "",
		Status:   200,
		Duration: time.Duration(0).Nanoseconds(),
	}
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Body reading error: %v", err)
			return
		}
		r.Body = string(bodyBytes)
		defer req.Body.Close()
	}

	doDelay(r)
	r.Status = getStatus(r)

	r.Duration = time.Since(started).Nanoseconds()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "HEAD, OPTIONS, GET, PUT, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if req.Method == "OPTIONS" || req.Method == "HEAD" || r.Status == 204 || r.Status == 304 {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(r.Status)
		w.Write([]byte{})
	} else {
		rb, _ := json.Marshal(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(r.Status)
		w.Write(rb)
	}

	log.Printf("Returning response: %v\n", r)

}

// getEchoSetting looks up an echo setting by name, expecting it to be
// specified via an HTTP header prefixed with "X-Echo-".
func getEchoSetting(r Request, name string) (string, bool) {
	key := http.CanonicalHeaderKey("X-Echo-" + name)
	if values, ok := r.Headers[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func doDelay(r Request) {
	value, ok := getEchoSetting(r, "Delay")
	if !ok {
		return
	}
	delay, err := time.ParseDuration(value)
	if err != nil {
		// Allow bare numbers to be treated as seconds, e.g. "2" == "2s".
		delay, err = time.ParseDuration(value + "s")
	}
	if err != nil {
		log.Printf("Invalid delay value %q: %v", value, err)
		return
	}
	log.Printf("Delaying response by %v\n", delay)
	time.Sleep(delay)
}

func getStatus(r Request) int {
	status := 200
	value, ok := getEchoSetting(r, "Status")
	if ok {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed >= 200 && parsed <= 599 {
			status = parsed
			if *debug {
				log.Printf("Returning explicit status: %v\n", status)
			}
		} else {
			log.Printf("Invalid status value %q, defaulting to %v", value, status)
		}
	}
	return status
}

func main() {
	flag.Parse()
	log.Print("Running echo-ish server on " + *port)
	http.HandleFunc("/", echo)
	err := http.ListenAndServe(":"+*port, nil)
	log.Printf("%v", err)
}
