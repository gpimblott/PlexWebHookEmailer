package notify

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
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
Send notifications that a new item has been received
*/
func EmailNewItem( server EmailDetails, from, to string, mc event.MediaContainer) {
	var body string

	// SendBySMTP an notify for the new item
	header := fmt.Sprintf("New item added to library %s\n\n", mc.LibrarySectionTitle)
	if mc.Track != (event.Track{}) {
		body = fmt.Sprintf("\t%s\n\t\t%s\n\t\t%s", mc.Track.GrandParent, mc.Track.Parent, mc.Track.Title)
	}

	if mc.Video != (event.Video{}) {
		body = fmt.Sprintf("\t%s\n\t\t%s\n\t\t%s", mc.Video.GrandParent, mc.Video.Parent, mc.Video.Title)
	}


	err := SendBySMTP(server.Server, server.Port, server.Username, server.Password, from, to,
		"plex: New item added to library", header+body)

	if err != nil {
		log.Printf("Error sending mail : %s", err)
	}
}

/**
Send an SMTP email
 */
func SendBySMTP(host, port, username, password, fromEmail, toEmail, subject, body string) error {

	from := mail.Address{"", fromEmail}
	to := mail.Address{"", toEmail}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := host + ":" + port
	auth := smtp.PlainAuth("", username, password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}


