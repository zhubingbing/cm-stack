// Package prometheus provides bindings to the Prometheus HTTP API v1:
// http://prometheus.io/docs/querying/api/
package prometheus

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type (

	// Timestamp is a helper for (un)marhalling time
	Timestamp time.Time

	Message struct {
		Text string `json:"text"`
	}

	// HookMessage is the message we receive from Alertmanager
	HookMessage struct {
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
		Status            string            `json:"status"`
		Receiver          string            `json:"receiver"`
		GroupLabels       map[string]string `json:"groupLabels"`
		CommonLabels      map[string]string `json:"commonLabels"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
		Alerts            []Alert           `json:"alerts"`
	}

	// Alert is a single alert.
	Alert struct {
		Status       string            `json:"status"`
		Labels       map[string]string `json:"labels"`
		Annotations  map[string]string `json:"annotations"`
		StartsAt     string            `json:"startsAt,omitempty"`
		EndsAt       string            `json:"EndsAt,omitempty"`
		GeneratorURL string            `json:"generatorURL"`
	}

	// just an example alert store. in a real hook, you would do something useful
	AlertStore struct {
		sync.Mutex
		capacity int
		alerts   []*HookMessage
	}
)

func GetAlerts() {

	capacity := 64


	s := &AlertStore{
		capacity: capacity,
	}




	//http.HandleFunc("/healthz", healthzHandler)
	//http.HandleFunc("/alerts", s.alertsHandler)
	//log.Fatal(http.ListenAndServe(*addr, nil))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func (s *AlertStore) alertsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getHandler(w, r)
	case http.MethodPost:
		s.postHandler(w, r)
	default:
		http.Error(w, "unsupported HTTP method", 400)
	}
}

func (s *AlertStore) getHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(enc)
	s.Lock()
	defer s.Unlock()

	if err := enc.Encode(s.alerts); err != nil {
		log.Printf("error encoding messages: %v", err)
	}
}

func (s *AlertStore) postHandler(w http.ResponseWriter, r *http.Request) {

	a, err := ioutil.ReadAll(r.Body)
	var msg Message

	var hook HookMessage

	err = json.Unmarshal(a, &hook)
	if err != nil {
		return
	}

	defer r.Body.Close()
	test := hook.Alerts


	s_url := "https://hook.bearychat.com/=bw8Sf/incoming/b42fcb5cdbbda95d831445f071a58ab2"
	contentType := "application/json;charset=utf-8"

	for i := 0; i < len(test); i++ {
		fmt.Println(test[i].Annotations["description"])
		msg.Text = test[i].Annotations["description"]
		b, err := json.Marshal(msg)

		if err != nil {
			fmt.Println("ss")
		}

		body := bytes.NewBuffer(b)

		resp, err := http.Post(s_url, contentType, body)
		if err != nil {
			log.Println("Post failed:", err)
			return
		}

		defer resp.Body.Close()
		bo, err := ioutil.ReadAll(resp.Body)
		fmt.Println("post: \n", string(bo))
	}

}