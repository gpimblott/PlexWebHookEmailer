package main

import (
	"log"
	"net/http"
	"pimblott.com/plexWebhookServer/environment"
	"pimblott.com/plexWebhookServer/notify"
	"pimblott.com/plexWebhookServer/plex"
	"pimblott.com/plexWebhookServer/plex/event"
)

var authToken string
var plexServer = environment.GetEnvOrStop("PLEX_SERVER")
var emailDetails = notify.EmailDetails{
	Server:   environment.GetEnvOrStop("MAIL_SERVER"),
	Port:     environment.GetEnvOrStop("MAIL_PORT"),
	Username: environment.GetEnvOrStop("MAIL_USERNAME"),
	Password: environment.GetEnvOrStop("MAIL_PASSWORD"),
}
var mailFrom = environment.GetEnvOrStop("MAIL_FROM")
var mailTo = environment.GetEnvOrStop("MAIL_TO")

/*
Handle a webhook
*/
func handleWebHook(w http.ResponseWriter, req *http.Request) {
	log.Printf("Received webhook")

	multiPartReader, err := req.MultipartReader()
	if err != nil {
		// Detect error type for the http answer
		if err == http.ErrNotMultipart || err == http.ErrMissingBoundary {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Println("can't create a multipart reader from request:", err)
		return
	}

	// Use the multipart reader to parse the request body
	payload, _, err := event.Extract(multiPartReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("can't decode POST request from plex server", err)
		return
	}

	log.Printf("%+v\n", *payload)
	switch payload.Event {
	case "library.new":
		item := payload.Metadata.Key
		log.Printf("New media [%s]", item)

		details, getErr := plex.GetMediaDetails(plexServer, authToken, item)
		if getErr != nil {
			log.Printf("Error getting item [%s] details: %s", item, getErr)
			return
		}

		notify.EmailNewItem(emailDetails, mailFrom, mailTo, details)

		break
	case "device.new":
		log.Printf("New device detected")
		break
	default:
		log.Printf("Event %s not implemented", payload.Event)
	}

	w.WriteHeader(200)
}

/*
Wrapper to login to plex and log any errors
THe application is stopped if a login occurs
*/
func login(username, password string) {
	log.Printf("plex login [%s]", username)
	token, authError := plex.Login(username, password)
	if authError != nil {
		log.Fatalf("Error authenticating with plex %s", authError)
	}
	authToken = token
}

/*
Application entry point.
Define the HTTP handler and start the server on the port defined by
environment variable PORT or 8090 if not defined.
*/
func main() {
	port := environment.GetEnvWithFallback("PORT", "8090")
	log.Printf("plex webhook server [%s] starting on %s...", plexServer, port)

	// Retrieve the credentials and login to plex
	login(
		environment.GetEnvOrStop("PLEX_USER"),
		environment.GetEnvOrStop("PLEX_PASSWORD"))

	// Create the server
	log.Printf("Server running")
	http.HandleFunc("/", handleWebHook)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
