package notify

import (
	"fmt"
	"log"
	"pimblott.com/plexWebhookServer/plex/event"
)

// Email details
type EmailDetails struct {
	Server   string
	Port     string
	Username string
	Password string
}

/**
Perform notifications that a new item has been received
*/
func EmailNewItem( server EmailDetails, from, to string, mc event.MediaContainer) {
	var body string

	// Send an notify for the new item
	header := fmt.Sprintf("New item added to library %s\n\n", mc.LibrarySectionTitle)
	if mc.Track != (event.Track{}) {
		body = fmt.Sprintf("\t%s\n\t\t%s\n\t\t%s", mc.Track.GrandParent, mc.Track.Parent, mc.Track.Title)
	}

	if mc.Video != (event.Video{}) {
		body = fmt.Sprintf("\t%s\n\t\t%s\n\t\t%s", mc.Video.GrandParent, mc.Video.Parent, mc.Video.Title)
	}

	err := Send(server.Server, server.Port, server.Username, server.Password, from, to,
		"plex: New item added to library", header+body)

	if err != nil {
		log.Printf("Error sending mail : %s", err)
	}
}

