package notify

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var webhookURL string

func init() {
	webhookURL = os.Getenv("WEBHOOKURL")
}

type body struct {
	Tenant string `json:"tenant"`
}

func Notify(tenant string) {
	if webhookURL == "" {
		log.Println("webhook URL not set. Not sending notification.")
		return
	}

	data, err := json.Marshal(body{Tenant: tenant})
	if err != nil {
		log.Println("webhook notify failed to marshal json")
		return
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Println("failed to POST to webhook URL")
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("webhook notify got status code %d", resp.StatusCode)
		return
	}
}
