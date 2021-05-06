package main

import (
	"encoding/json"
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
	decoder := json.NewDecoder(r.Body)

	var body presentRequest
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	client, err := civogo.NewClient(os.Getenv("CIVO_API_TOKEN"), os.Getenv("CIVO_API_REGION"))
	if err != nil {
		log.Errorf("Couldn't get Civo client: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	record := strings.Split(body.Fqdn, ".")[0]
	zone := body.Fqdn[len(record)+1 : len(body.Fqdn)-1]

	domain, err := client.GetDNSDomain(zone)
	if err != nil {
		log.Errorf("Couldn't get DNS zone: `%s`: %s", zone, err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	config := &civogo.DNSRecordConfig{
		Name:     record,
		Value:    body.Value,
		Type:     civogo.DNSRecordTypeTXT,
		Priority: 10,
		TTL:      300,
	}

	log.Infof("Creating DNS record `%s` in zone `%s`", record, zone)

	_, err = client.CreateDNSRecord(domain.ID, config)
	if err != nil {
		log.Errorf("Failed to create DNS record `%s`: %s", record, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	log.Infof("Successfully created DNS record `%s` in zone `%s`", record, zone)
}

func cleanup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var body presentRequest
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	client, err := civogo.NewClient(os.Getenv("CIVO_API_TOKEN"), os.Getenv("CIVO_API_REGION"))
	if err != nil {
		log.Errorf("Couldn't get Civo client: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	record := strings.Split(body.Fqdn, ".")[0]
	zone := body.Fqdn[len(record)+1 : len(body.Fqdn)-1]

	domain, err := client.GetDNSDomain(zone)
	if err != nil {
		log.Errorf("Couldn't get DNS zone: `%s`: %s", zone, err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	records, err := client.ListDNSRecords(domain.ID)
	if err != nil {
		log.Errorf("Couldn't get DNS records: %s", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if res.Result != "success" {
		log.Errorf("Couldn't delete DNS record `%s`: %s", record, res)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Infof("Successfully deleted DNS record `%s` from zone `%s`", record, zone)
}

func main() {
	http.HandleFunc("/present", present)
	http.HandleFunc("/cleanup", cleanup)
	http.ListenAndServe(":8080", nil)
}
