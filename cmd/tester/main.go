package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	gearman "github.com/mikespook/gearman-go/client"
)

func TestServer() {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.UserAgent())
		fmt.Println(r.URL)
		fmt.Println(r.Header)
		w.WriteHeader(200)
	})

	log.Println("Starting server...")
	http.ListenAndServe("127.0.0.1:3202", nil)
}

func main() {

	go TestServer()

	log.Println("Connecting to gearman")
	client, err := gearman.New("tcp", "127.0.0.1:4730")
	if err != nil {
		log.Fatalf("Error launching client %s", err)
	}
	log.Println("Testing ")
	body := `{"query":{"test":"toto"},"headers":{"X-Tester":"TestValue"}}`
	_, err = client.Do("local_http", []byte(body), 0, func(r *gearman.Response) {
		log.Println("Response")
		log.Println(r)
	})

	time.Sleep(20 * time.Second)
}
