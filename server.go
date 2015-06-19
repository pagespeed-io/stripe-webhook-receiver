package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/thorduri/pushover"
)

var (
	log            = logrus.New()
	pushoverClient *pushover.Pushover
	port           = flag.String("p", os.Getenv("PORT"), "Port to listen on")
	poToken        = flag.String("ptoken", os.Getenv("PUSHOVER_TOKEN"), "Pushover application token")
	poUser         = flag.String("puser", os.Getenv("PUSHOVER_USER"), "Pushover user (or group)")
)

type Payload struct {
	Data struct {
		Object struct {
			Name  string `json:"name"`
			Brand string `json:"brand"`
		} `json:"object"`
	} `json:"data"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

func init() {
	flag.Parse()
	if *poToken == "" || *poUser == "" {
		log.Fatal("Please provide pushover credentials.")
	}
	pushoverClient, _ = pushover.NewPushover(*poToken, *poUser)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", webhookHandler).Methods("POST")

	fmt.Printf("Listening on :%s\n", *port)
	err := http.ListenAndServe(":"+*port, r)
	if err != nil {
		log.Fatal("Couldn’t start server", err)
	}
}

func webhookHandler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Couldn’t read request body")
		http.Error(rw, "Internal Server Error", 500)
		return
	}

	var r Payload
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Warn("Error parsing JSON body")

		http.Error(rw, "Invalid JSON body", 400)
		return
	}

	log.WithFields(logrus.Fields{
		"id": r.Id,
	}).Info("Received incoming request")

	switch r.Type {
	case "customer.source.created":
		m := &pushover.Message{
			Message: fmt.Sprintf("%s added a %s card.", r.Data.Object.Name, r.Data.Object.Brand),
			Title:   fmt.Sprintf("%s", r.Type),
		}

		go sendNotification(m)
		fmt.Fprintf(rw, "%s\n", r.Type)

	default:
		http.NotFound(rw, req)
	}
}

func sendNotification(m *pushover.Message) {
	req, _, err := pushoverClient.Push(m)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"id":    req,
		"error": err,
	}).Info("Sent notification.")
}
