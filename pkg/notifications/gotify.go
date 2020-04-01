package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	t "github.com/jm9e/watchtower/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	gotifyType = "gotify"
)

type gotifyTypeNotifier struct {
	gotifyURL      string
	gotifyAppToken string
	logLevels      []log.Level
}

func newGotifyNotifier(c *cobra.Command, acceptedLogLevels []log.Level) t.Notifier {
	flags := c.PersistentFlags()

	gotifyURL, _ := flags.GetString("notification-gotify-url")
	if len(gotifyURL) < 1 {
		log.Fatal("Required argument --notification-gotify-url(cli) or WATCHTOWER_NOTIFICATION_GOTIFY_URL(env) is empty.")
	} else if !(strings.HasPrefix(gotifyURL, "http://") || strings.HasPrefix(gotifyURL, "https://")) {
		log.Fatal("Gotify URL must start with \"http://\" or \"https://\"")
	} else if strings.HasPrefix(gotifyURL, "http://") {
		log.Warn("Using an HTTP url for Gotify is insecure")
	}

	gotifyToken, _ := flags.GetString("notification-gotify-token")
	if len(gotifyToken) < 1 {
		log.Fatal("Required argument --notification-gotify-token(cli) or WATCHTOWER_NOTIFICATION_GOTIFY_TOKEN(env) is empty.")
	}

	n := &gotifyTypeNotifier{
		gotifyURL:      gotifyURL,
		gotifyAppToken: gotifyToken,
		logLevels:      acceptedLogLevels,
	}

	log.AddHook(n)

	return n
}

func (n *gotifyTypeNotifier) StartNotification() {}

func (n *gotifyTypeNotifier) SendNotification() {}

func (n *gotifyTypeNotifier) Levels() []log.Level {
	return n.logLevels
}

func (n *gotifyTypeNotifier) getURL() string {
	url := n.gotifyURL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url + "message?token=" + n.gotifyAppToken
}

func (n *gotifyTypeNotifier) Fire(entry *log.Entry) error {

	go func() {
		jsonBody, err := json.Marshal(gotifyMessage{
			Message:  "(" + entry.Level.String() + "): " + entry.Message,
			Title:    "Watchtower",
			Priority: 0,
		})
		if err != nil {
			fmt.Println("Failed to create JSON body for Gotify notification: ", err)
			return
		}

		jsonBodyBuffer := bytes.NewBuffer([]byte(jsonBody))
		resp, err := http.Post(n.getURL(), "application/json", jsonBodyBuffer)
		if err != nil {
			fmt.Println("Failed to send Gotify notification: ", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			fmt.Printf("Gotify notification returned %d HTTP status code", resp.StatusCode)
		}

	}()
	return nil
}

type gotifyMessage struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}
