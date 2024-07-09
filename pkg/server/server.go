package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
	"github.com/sentiweb/gearworkers/pkg/admin"
	"github.com/sentiweb/gearworkers/pkg/config"
)

type HttpServer struct {
	addr          string
	gearmanServer string
}

func NewHttpServer(conf *config.AppConfig) *HttpServer {
	return &HttpServer{addr: conf.Server.Addr, gearmanServer: conf.GearmanServer}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(err string, w http.ResponseWriter, status int) {
	e := ErrorResponse{Error: err}
	r, _ := json.Marshal(e)
	w.Write(r)
	w.WriteHeader(status)
}

func (h *HttpServer) Start() error {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {

		stats, err := admin.Load(h.gearmanServer)
		if err != nil {
			WriteError(fmt.Sprintf("Unable to get status : %s", err), w, 501)
		}
		b, err := json.Marshal(stats)
		if err != nil {
			WriteError(fmt.Sprintf("Unable to serialize status : %s", err), w, 501)
		}
		w.Write(b)
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})

	log.Printf("Starting server at port %s ...", h.addr)
	return http.ListenAndServe(h.addr, nil)
}
