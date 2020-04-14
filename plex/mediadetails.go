package plex

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"pimblott.com/plexWebhookServer/plex/event"
)

/*
 Query the PLex server for the details of the new item added
 */
func GetMediaDetails(plexServer, authToken, item string) (mc event.MediaContainer, err error) {
	// Get the detailed item information from the plex server
	url := fmt.Sprintf("%s%s?X-Plex-Token=%s", plexServer, item, authToken)
	resp, getErr := http.Get(url)
	if getErr != nil {
		log.Printf("Error getting item info : %s", getErr)
		return event.MediaContainer{}, getErr
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP Status error %d for %s",resp.StatusCode, url)
		return mc, nil
	}

	defer resp.Body.Close()

	// Decode the XML response
	details := event.MediaContainer{}
	DecodeMediaContainer(resp.Body, &details)

	return details, nil
}

/*
 Map the XML structure returned from Plex to a MediaContainer structure
 */
func DecodeMediaContainer(reader io.Reader, mc *event.MediaContainer) {
	decoder := xml.NewDecoder(reader)
	decoder.Decode(mc)
}
