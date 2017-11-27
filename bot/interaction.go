package bot

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var InteractionActionDef = map[string]InteractionHandleFunc{}

type InteractionHandleFunc func(original slack.Message, action slack.AttachmentAction) (msg slack.Message, err error)

const (
	ActionCancel = "cancel"
)

func init() {
	InteractionActionDef[ActionCancel] = interactionCancelAction
}

// interactionHandler handles interactive message response.
type InteractionHandler struct {
	Client            *slack.Client
	VerificationToken string
}

func (h InteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	if message.Token != h.VerificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	action := message.Actions[0]
	actionName := action.Name

	if handler, ok := InteractionActionDef[actionName]; ok {
		message.OriginalMessage.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
		msg, err := handler(message.OriginalMessage, action)
		if err != nil {
			log.Printf("[ERROR] Handling %q is failed: %s", actionName, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&msg)
		return
	}

	log.Printf("[ERROR] Invalid action was submitted: %s", actionName)
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func interactionCancelAction(original slack.Message, _ slack.AttachmentAction) (slack.Message, error) {
	title := ":x: キャンセルされました"
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Short: false,
		},
	}
	return original, nil
}
