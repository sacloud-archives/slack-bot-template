package handler

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/sacloud/slack-bot-template/bot"
)

func init() {
	bot.AppendMessageHandler(bot.SlackMessageHandler{
		Order:   bot.HandlerPriorityNormal,
		Handler: handleListInteraction,
	})
	bot.InteractionActionDef[actionListSelect] = handleListSelect
	bot.InteractionActionDef[actionListExec] = handleListExec

}

const (
	actionListSelect = "actionListSelect"
	actionListExec   = "actionListExec"
)

var (
	listExampleCommands = []string{
		"list",
	}
)

func handleListInteraction(s *bot.SlackListener, ev *slack.MessageEvent) (handled bool, err error) {
	// Only response in specific channel. Ignore else.
	if ev.Channel != s.ChannelID {
		return
	}

	if !s.HasCommandPrefix(ev, listExampleCommands...) {
		return
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			{
				Color:      "#e2e2e2",
				CallbackID: actionListSelect,
				Actions: []slack.AttachmentAction{
					{
						Name: actionListSelect,
						Type: "select",
						Options: []slack.AttachmentActionOption{
							{
								Text:  "value1",
								Value: "value1",
							},
							{
								Text:  "value2",
								Value: "value2",
							},
							{
								Text:  "value3",
								Value: "value3",
							},
						},
					},
					{
						Name:  bot.ActionCancel,
						Text:  "キャンセル",
						Type:  "select",
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

func handleListSelect(original slack.Message, action slack.AttachmentAction) (slack.Message, error) {
	value := action.SelectedOptions[0].Value

	original.Text = fmt.Sprintf("Selected: %s\nOK or Cannel?", value)
	original.Attachments[0].Actions = []slack.AttachmentAction{
		{
			Name:  actionListExec,
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

func handleListExec(original slack.Message, action slack.AttachmentAction) (slack.Message, error) {
	value := action.Value
	original.Text = fmt.Sprintf("Submitted: %s", value)
	return original, nil
}
