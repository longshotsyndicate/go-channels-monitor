package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/longshotsyndicate/go-channels-monitor/channel"
)

type Service struct {
	url     string
	errc    chan error
	name    string
	monitor *channel.Monitor
}

func New(monitor *channel.Monitor, serviceName string, url string, errc chan error) *Service {

	return &Service{
		url:     url,
		errc:    errc,
		name:    serviceName,
		monitor: monitor,
	}
}

func (this *Service) Start() {
	http.HandleFunc("/channels", this.chanHandler)
	go this.start()
}

func (this *Service) start() {
	if err := http.ListenAndServe(this.url, nil); err != nil {
		this.errc <- err
	}
}

func (this *Service) chanHandler(w http.ResponseWriter, r *http.Request) {
	chStats := this.monitor.GetAll()

	resp := &ServiceChannelsStatus{
		Service:  this.name,
		Channels: chStats,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Printf("Error: %#v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonResp)
}

type ServiceChannelsStatus struct {
	Service  string `json:"service"`
	Channels map[string]*channel.ChanState
}

type Config struct {
	Name string
	Url  string
}
