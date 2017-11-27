package handler

import (
	"github.com/nlopes/slack"
	"github.com/sacloud/slack-bot-template/bot"
)

func init() {
	bot.AppendMessageHandler(bot.SlackMessageHandler{
		Order:   bot.HandlerPriorityNormal,
		Handler: handleMessage,
	})
}

var (
	exampleResponse = "メッセージへの応答サンプルです。\n\n"
	exampleTitle    = "メッセージ応答サンプルのタイトル"
	exampleMsg      = "メッセージ応答サンプルの本文:+1:"

	exampleCommands = []string{
		"message",
	}
)

func handleMessage(s *bot.SlackListener, ev *slack.MessageEvent) (handled bool, err error) {

	if s.HasCommandPrefix(ev, exampleCommands...) {

		// message to channel
		params := slack.PostMessageParameters{
			Markdown: true,
			Attachments: []slack.Attachment{
				{
					Color:      "#e2e2e2",
					Pretext:    exampleTitle,
					Text:       exampleMsg,
					MarkdownIn: []string{"text"},
				},
			},
		}
		err = s.ResponseToChannel(ev, exampleResponse, params)
		handled = true
	}

	return
}
