package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

//Launch a webserver that "handles" workflows whenever it receives a particular request
func main() {

	start := time.Now()
	ch := make(chan string)
	max_concurrency := 10

	workflows := make(string[], max_concurrency)

	for _, workflow := range workflows {
		go workflowHandler(workflow, ch)
	}

	for range os.Args[1:] {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/iwantaworkflow", workflowHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// Handler echoes the HTTP request
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}

	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)

	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}

	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

// counter echoes the number of calls so far
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}


func workflowHandler(url string, ch chan<- string) {
	start := time.Now()
	
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) //Send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // Don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}