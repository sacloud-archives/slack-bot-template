package handler

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/sacloud/slack-bot-template/bot"
)

func init() {
	bot.AppendMessageHandler(bot.SlackMessageHandler{
		Order:   bot.HandlerPriorityNormal,
		Handler: handleButtonInteraction,
	})
	bot.InteractionActionDef[actionButtonSelect] = handleButtonSelect
	bot.InteractionActionDef[actionButtonExec] = handleButtonExec

}

const (
	actionButtonSelect = "actionButtonSelect"
	actionButtonExec   = "actionButtonExec"
)

var (
	buttonExampleCommands = []string{
		"button",
	}
)

func handleButtonInteraction(s *bot.SlackListener, ev *slack.MessageEvent) (handled bool, err error) {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.ChannelID {
		return
	}

	if !s.HasCommandPrefix(ev, buttonExampleCommands...) {
		return
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			{
				Color:      "#e2e2e2",
				CallbackID: actionButtonSelect,
				Actions: []slack.AttachmentAction{
					{
						Name:  actionButtonSelect,
						Text:  "value1",
						Value: "value1",
						Type:  "button",
					},
					{
						Name:  actionButtonSelect,
						Text:  "value2",
						Value: "value2",
						Type:  "button",
					},
					{
						Name:  actionButtonSelect,
						Text:  "value3",
						Value: "value3",
						Type:  "button",
					},
					{
						Name:  bot.ActionCancel,
						Text:  "キャンセル",
						Type:  "button",
						Style: "danger",
					},
				},
			},
		},
	}
	err = s.ResponseToChannel(ev, "Please select\n", params)
	handled = true
	return
}

func handleButtonSelect(original slack.Message, action slack.AttachmentAction) (slack.Message, error) {
	value := action.Value

	original.Text = fmt.Sprintf("Selected: %s\nOK or Cannel?", value)
	original.Attachments[0].Actions = []slack.AttachmentAction{
		{
			Name:  actionButtonExec,
			Text:  "OK",
			Type:  "button",
			Style: "primary",
			Value: value,
		},

		{
			Name:  bot.ActionCancel,
			Text:  "キャンセル",
			Type:  "button",
			Style: "danger",
		},
	}
	return original, nil
}

func handleButtonExec(original slack.Message, action slack.AttachmentAction) (slack.Message, error) {
	value := action.Value
	original.Text = fmt.Sprintf("Submitted: %s", value)
	return original, nil
}
