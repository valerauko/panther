package main

import (
	"encoding/json"
	"fmt"
	"github.com/civo/civogo"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type presentRequest struct {
	Fqdn  string `json:"fqdn"`
	Value string `json:"value"`
}

func present(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Received present request %s", r.Body)

	decoder := json.NewDecoder(r.Body)

	var body presentRequest
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	client, err := civogo.NewClient(os.Getenv("CIVO_API_TOKEN"), os.Getenv("CIVO_API_REGION"))
	if err != nil {
		log.Errorf("Couldn't get Civo client: %s", err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	record := strings.Split(body.Fqdn, ".")[0]
	zone := body.Fqdn[len(record)+1 : len(body.Fqdn)-1]

	domain, err := client.GetDNSDomain(zone)
	if err != nil {
		log.Errorf("Couldn't get DNS zone: `%s`: %s", zone, err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
		return
	}

	config := &civogo.DNSRecordConfig{
		Name:     record,
		Value:    body.Value,
		Type:     civogo.DNSRecordTypeTXT,
		Priority: 10,
		TTL:      300,
	}

	_, err = client.CreateDNSRecord(domain.ID, config)
	if err != nil {
		log.Errorf("Failed to create DNS record `%s`: %s", record, err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	log.Infof("Successfully created DNS record `%s` in zone `%s`", record, zone)
}

func cleanup(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Received cleanup request %s", r.Body)

	decoder := json.NewDecoder(r.Body)

	var body presentRequest
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	client, err := civogo.NewClient(os.Getenv("CIVO_API_TOKEN"), os.Getenv("CIVO_API_REGION"))
	if err != nil {
		log.Errorf("Couldn't get Civo client: %s", err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	record := strings.Split(body.Fqdn, ".")[0]
	zone := body.Fqdn[len(record)+1 : len(body.Fqdn)-1]

	domain, err := client.GetDNSDomain(zone)
	if err != nil {
		log.Errorf("Couldn't get DNS zone: `%s`: %s", zone, err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
		return
	}

	records, err := client.ListDNSRecords(domain.ID)
	if err != nil {
		log.Errorf("Couldn't get DNS records: %s", err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
		return
	}

	var existingRecord civogo.DNSRecord
	for _, rec := range records {
		if rec.Name == record && rec.Value == body.Value {
			existingRecord = rec
		}
	}

	res, err := client.DeleteDNSRecord(&existingRecord)
	if err != nil {
		log.Errorf("Couldn't delete DNS record `%s`: %s", record, err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if res.Result != "success" {
		log.Errorf("Couldn't delete DNS record `%s`: %s", record, res)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	log.Infof("Successfully deleted DNS record `%s` from zone `%s`", record, zone)
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Infof("Health check OK")
}

func main() {
	http.HandleFunc("/present", present)
	http.HandleFunc("/cleanup", cleanup)
	http.HandleFunc("/health", health)
	http.ListenAndServe(":8080", nil)
}
