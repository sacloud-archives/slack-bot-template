package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
	"sort"
)

var MessageHandlers []SlackMessageHandler

func AppendMessageHandler(handler SlackMessageHandler) {
	MessageHandlers = append(MessageHandlers, handler)
}

const (
	HandlerPriorityHighest    = 5
	HandlerPriorityHigh       = 10
	HandlerPriorityNormalHign = 100
	HandlerPriorityNormal     = 200
	HandlerPriorityNormalLow  = 300
	HandlerPriorityLow        = 1000
	HandlerPriorityLowest     = 9999
)

type SlackListener struct {
	Client    *slack.Client
	BotID     string
	ChannelID string
}

type SlackMessageHandleFunc func(s *SlackListener, event *slack.MessageEvent) (handled bool, err error)
type SlackMessageHandler struct {
	Handler SlackMessageHandleFunc
	Order   int
}

// ListenAndResponse listens slack events and response
// particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.Client.NewRTM()

	// prepare handlers
	sort.Slice(MessageHandlers, func(i, j int) bool {
		return MessageHandlers[i].Order < MessageHandlers[j].Order
	})

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		//log.Printf("[DEBUG]: %#v", msg.Data)
		switch ev := msg.Data.(type) {
		case *slack.ConnectingEvent:
			log.Printf("[INFO] Slack RTM connecting")
		case *slack.ConnectedEvent:
			log.Printf("[INFO] Slack RTM connected")
		case *slack.ConnectionErrorEvent:
			log.Printf("[WARN] Slack RTM opening connection is failed: %s", ev.ErrorObj)
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
			//case *slack.UnmarshallingErrorEvent:
			//	log.Printf("[ErrorTrace] %s", ev.Error())
		}
	}
}

// handleMesageEvent handles message events.
func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {

	// Only response mention to bot. Ignore else.
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.BotID)) {
		return nil
	}

	err := s.applyHandlers(ev)
	if err != nil {
		return err
	}

	return nil
}

func (s *SlackListener) applyHandlers(ev *slack.MessageEvent) error {

	for _, h := range MessageHandlers {
		handled, err := h.Handler(s, ev)
		if handled {
			return err
		}
	}

	return fmt.Errorf("The message was not handled by any handlers")
}

func (s *SlackListener) HasCommandPrefix(ev *slack.MessageEvent, commands ...string) bool {
	tokens := s.MessageTokens(ev)
	if len(tokens) == 0 {
		return false
	}

	for _, cmd := range commands {
		if strings.HasPrefix(strings.ToLower(tokens[0]), cmd) {
			return true
		}
	}
	return false
}

func (s *SlackListener) HasMessageToken(ev *slack.MessageEvent) bool {
	return len(s.MessageTokens(ev)) > 0
}

func (s *SlackListener) MessageTokens(ev *slack.MessageEvent) []string {
	return strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
}

func (s *SlackListener) ResponseToUserIM(ev *slack.MessageEvent, msg string, params slack.PostMessageParameters) error {
	_, _, imc, err := s.Client.OpenIMChannel(ev.Msg.User)
	if err != nil {
		return err
	}

	if _, _, err := s.Client.PostMessage(imc, msg, params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}

func (s *SlackListener) ResponseToChannel(ev *slack.MessageEvent, msg string, params slack.PostMessageParameters) error {
	if _, _, err := s.Client.PostMessage(ev.Channel, msg, params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}
	return nil
}
