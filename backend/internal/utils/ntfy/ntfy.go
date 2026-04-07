package ntfy

import (
	"fmt"
	"net/http"
	"strings"
)

// ========================================
// Types
// ========================================
type SendNotificationParams struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority string `json:"priority,omitempty"` // "min", "low", "default", "high", "urgent"
	Tags     string `json:"tags,omitempty"`     // comma-separated list of tags, e.g. "warning,skull"
}

// ========================================
// Errors
// ========================================

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error %s: %s - %v", e.Code, e.Message, e.Cause)
}

const (
	ErrFailedToSendNotification = "FAILED_TO_SEND_NOTIFICATION"
)

// ========================================
// Functions
// ========================================

// SendNotification sends a notification to the ntfy server.
func SendNotification(notification SendNotificationParams) error {
    req, err := http.NewRequest(
        "POST",
        "https://ntfy.xerverlab.com/alerts",
        strings.NewReader(notification.Message),
    )
    if err != nil {
        return &Error{
            Code:    ErrFailedToSendNotification,
            Message: "Failed to create HTTP request",
            Cause:   err,
        }
    }

    req.Header.Set("Title", notification.Title)
    req.Header.Set("Priority", notification.Priority)
    req.Header.Set("Tags", notification.Tags)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return &Error{
            Code:    ErrFailedToSendNotification,
            Message: "Failed to send notification",
            Cause:   err,
        }
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return &Error{
            Code:    ErrFailedToSendNotification,
            Message: fmt.Sprintf("Unexpected HTTP status: %s", resp.Status),
        }
    }

    return nil
}
