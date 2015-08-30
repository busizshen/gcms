package messages

import "github.com/gin-gonic/contrib/sessions"

// Message is used for storing Flash messages with extra metadata,
// rather than just storing flash messages strings.
type Message struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Add is a helper function that adds a Message struct as flash message.
func Add(session sessions.Session, message, messageType string) {
	session.AddFlash(&Message{
		Message: message,
		Type:    messageType,
	})
}

// GetMessages is a Helper function to return all flash messages.
func GetMessages(session sessions.Session) []interface{} {
	return session.Flashes()
}
