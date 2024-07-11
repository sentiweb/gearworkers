package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	gearman "github.com/mikespook/gearman-go/client"
	"github.com/sentiweb/gearworkers/pkg/types"
)

func GearmanResponseHandler(r *gearman.Response) {
	result, err := r.Result()
	if err != nil {
		log.Printf("Error in response %s", err)
	}
	log.Println(string(result))
}

func main() {

	node := flag.String("node", "", "Node")

	var err error

	flag.Parse()

	name := "expect"

	payload := types.HttpJobPayload{
		UrlParams:   map[string]string{"node": *node},
		QueryParams: map[string]string{"status": "0"},
	}

	body, _ := json.Marshal(payload)

	client, err := gearman.New("tcp", "127.0.0.1:4730")
	if err != nil {
		log.Fatalf("Error launching client %s", err)
	}
	var exitCode int = 0
	_, err = client.DoBg(name, []byte(body), 0)
	if err != nil {
		log.Printf("Error during send : %s", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
