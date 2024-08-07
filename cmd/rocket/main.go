package main

import (
	"encoding/json"
	"flag"
	"io"
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

type SimpleMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func loadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func main() {

	channel := flag.String("channel", "", "Channel")
	text := flag.String("text", "", "Text")
	file := flag.String("file", "", "JSON file containing message")
	wait := flag.Bool("wait", false, "Wait for response")
	quiet := flag.Bool("quiet", false, "Do not log message")
	name := flag.String("name", "chat", "Gearman function name to call to send message")

	var body []byte
	var err error

	flag.Parse()

	if *file != "" {
		body, err = loadFile(*file)
		if err != nil {
			log.Fatalf("Error loading file %s : %s", *file, err)
		}
	} else {
		if *channel == "" || *text == "" {
			log.Fatalf("channel and text must not be empty")
		}
		m := SimpleMessage{Channel: *channel, Text: *text}
		payload := types.HttpJobPayload{Body: m}
		body, _ = json.Marshal(payload)
	}

	if *name == "" {
		log.Fatalf("name must not be empty")
	}

	if !*quiet {
		log.Println("Connecting to gearman")
	}
	client, err := gearman.New("tcp", "127.0.0.1:4730")
	if err != nil {
		log.Fatalf("Error launching client %s", err)
	}
	var exitCode int = 0
	if *wait {
		_, err = client.Do(*name, []byte(body), 0, func(r *gearman.Response) {
			if !*quiet {
				result, result_err := r.Result()
				if result_err != nil {
					log.Printf("Gearman response error : %s", result_err)
					exitCode = 1
				} else {
					log.Println(result)
				}
			}
		})
	} else {
		_, err = client.DoBg(*name, []byte(body), 0)
	}
	if err != nil {
		log.Printf("Error during send : %s", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
