package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	fs := http.FileServer(http.Dir("."))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %q", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.UserAgent())
		fs.ServeHTTP(w, r)
	})
	log.Printf("Serving httpdocs on http://localhost:%s", *port)
	if err := http.ListenAndServe(":"+*port, handler); err != nil {
		log.Fatal(err)
	}
}
