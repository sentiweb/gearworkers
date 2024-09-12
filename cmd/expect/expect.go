package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	gearman "github.com/mikespook/gearman-go/client"
	"github.com/sentiweb/gearworkers/pkg/types"
)

func main() {

	node := flag.String("node", "", "Node")
	//wait := flag.Bool("wait", false, "Wait for completion")

	var err error

	flag.Parse()

	name := "expect"

	if *node == "" {
		fmt.Println("Error -node must be provided")
		flag.Usage()
		return
	}

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
	var handle string
	handle, err = client.DoBg(name, []byte(body), 0)
	if err != nil {
		log.Printf("Error during send : %s", err)
		exitCode = 1
	}
	fmt.Println(handle)
	start := time.Now()
	maxWait := 20 * time.Second
	for {
		status, err := client.Status(handle)
		if err != nil {
			fmt.Println("Error", err)
			break
		}
		if status != nil {
			fmt.Printf("Final Known: %t Running: %t", status.Known, status.Running)
			fmt.Printf("Completed %d / %d\n", status.Numerator, status.Denominator)
			if !(status.Running || status.Known) {
				break
			}
		}
		<-time.After(time.Second)
		if time.Since(start) > maxWait {
			break
		}
	}
	os.Exit(exitCode)
}
