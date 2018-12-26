package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// WebAPI represents the REST API server.
type WebAPI struct {
	server  *http.Server
	devices DeviceLibrary
}

func CreateWebAPI(port int, devices DeviceLibrary) *WebAPI {
	portStr := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: portStr}
	webAPI := &WebAPI{
		server:  server,
		devices: devices}
	http.Handle("/", webAPI)
	return webAPI
}

// Start starts the HTTP server. Stop it using the Stop function. Non-blocking.
// Returns a channel that is written to when the HTTP server has stopped.
func (wa *WebAPI) Start() chan bool {
	done := make(chan bool)

	go func() {
		log.Printf("Starting Web API on port %s\n", wa.server.Addr)
		if err := wa.server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("WebAPI: ListenAndServe() shutdown reason: %s", err)
		}
		done <- true // Signal that http server has stopped
	}()
	return done
}

// Stop stops the HTTP server.
func (wa *WebAPI) Stop() {
	wa.server.Shutdown(context.Background())
}

// ServeHTTP handles incoming HTTP requests
func (wa *WebAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "devices" {
		wa.handleDevices(w, r)
	} else if head == "shutdown" && r.Method == "POST" {
		wa.Stop()
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: %s or method %s!", r.URL.Path, r.Method)
	}
}

// handleDevices handles url devices/*
func (wa *WebAPI) handleDevices(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	id, idErr := strconv.Atoi(head)
	if head == "" && r.Method == "GET" {
		deviceIDs, err := wa.devices.GetDeviceIds()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		toJSON(deviceIDs, w)
	} else if idErr == nil {
		// Check that device exists
		deviceIDs, _ := wa.devices.GetDeviceIds()
		for deviceID := range deviceIDs {
			if deviceID == id {
				wa.handleDeviceID(id, w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Device with id %d does not exist", id)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: devices/%s or method %s!", r.URL.Path, r.Method)
	}
}

// handleDeviceID handles url devices/<id>/*
func (wa *WebAPI) handleDeviceID(id int, w http.ResponseWriter, r *http.Request) {
	var head string
	originalPath := r.URL.Path
	head, r.URL.Path = shiftPath(r.URL.Path)
	if (head == "on" || head == "off") && r.URL.Path == "/" && r.Method == "POST" {
		if wa.devices.SupportsOnOff(id) {
			var err error
			if head == "on" {
				err = wa.devices.TurnOn(id)
			} else {
				err = wa.devices.TurnOff(id)
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Devices with id %d don't support on/off", id)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: devices/%d%s or method %s!", id, originalPath, r.Method)
	}
}

// toJSON converts the v object to JSON and writes result to the response
func toJSON(v interface{}, w http.ResponseWriter) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// shiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
