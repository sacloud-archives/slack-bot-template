package main

import (
	"log"
	"net/http"
	"os"

	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"github.com/sacloud/slack-bot-template/bot"
	_ "github.com/sacloud/slack-bot-template/bot/handler" // for run init()
)

// https://api.slack.com/slack-apps
// https://api.slack.com/internal-integrations
type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"3000"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID string `envconfig:"CHANNEL_ID" required:"true"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)
	slackListener := &bot.SlackListener{
		Client:    client,
		BotID:     env.BotID,
		ChannelID: env.ChannelID,
	}

	// message handler(RTM)
	go slackListener.ListenAndResponse()

	// Register handler to receive interactive message
	// responses from slack (kicked by user action)
	http.Handle("/interaction", bot.InteractionHandler{
		Client:            client,
		VerificationToken: env.VerificationToken,
	})

	// for health-check
	http.HandleFunc("/status.html", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "OK")
	})

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0
}
