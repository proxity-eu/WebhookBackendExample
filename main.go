// Copyright 2021 Proxity SA
// https://www.proxity.eu/

package main

import (
	"container/list"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type webhook struct {
	Timestamp  time.Time `json:"timestamp"`
	ID         string    `json:"id"`
	RegionID   string    `json:"region_id"`
	Data       string    `json:"data"`
	DeviceData string    `json:"device_data"`
}

func main() {
	list := list.New()

	mux := &http.ServeMux{}
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			return
		}

		wh := webhook{Timestamp: time.Now()}
		if err := json.NewDecoder(req.Body).Decode(&wh); err != nil {
			return
		}

		log.Printf("Invoking webhook %+v", wh)

		list.PushFront(wh)
		if list.Len() > 100 {
			list.Remove(list.Back())
		}
	})
	mux.HandleFunc("/stats", func(w http.ResponseWriter, req *http.Request) {
		var items []webhook
		for e := list.Front(); e != nil; e = e.Next() {
			items = append(items, e.Value.(webhook))
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(items)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Running webhook test backend on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
