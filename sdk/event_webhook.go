package sendgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// EventWebhook is a Sendgrid event webhook settings.
type EventWebhook struct { //nolint:maligned
	ID                string `json:"id,omitempty"`
	Enabled           bool   `json:"enabled"`
	URL               string `json:"url,omitempty"`
	FriendlyName      string `json:"friendly_name,omitempty"`       //nolint:tagliatelle
	GroupResubscribe  bool   `json:"group_resubscribe"`             //nolint:tagliatelle
	Delivered         bool   `json:"delivered"`
	GroupUnsubscribe  bool   `json:"group_unsubscribe"`             //nolint:tagliatelle
	SpamReport        bool   `json:"spam_report"`                   //nolint:tagliatelle
	Bounce            bool   `json:"bounce"`
	Deferred          bool   `json:"deferred"`
	Unsubscribe       bool   `json:"unsubscribe"`
	Processed         bool   `json:"processed"`
	Open              bool   `json:"open"`
	Click             bool   `json:"click"`
	Dropped           bool   `json:"dropped"`
	OAuthClientID     string `json:"oauth_client_id,omitempty"`     //nolint:tagliatelle
	OAuthClientSecret string `json:"oauth_client_secret,omitempty"` //nolint:tagliatelle
	OAuthTokenURL     string `json:"oauth_token_url,omitempty"`     //nolint:tagliatelle
}

// EventWebhookList is the response from the list all webhooks endpoint.
type EventWebhookList struct {
	Webhooks   []EventWebhook `json:"webhooks"`
	MaxAllowed int            `json:"max_allowed"` //nolint:tagliatelle
}

type EventWebhookSigning struct {
	Enabled   bool   `json:"enabled"`
	PublicKey string `json:"public_key"` //nolint:tagliatelle
}

func parseEventWebhook(respBody string) (*EventWebhook, RequestError) {
	var body EventWebhook
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing event webhook: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseEventWebhookList(respBody string) (*EventWebhookList, RequestError) {
	var body EventWebhookList
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing event webhook list: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

func parseEventWebhookSigning(respBody string) (*EventWebhookSigning, RequestError) {
	var body EventWebhookSigning
	if err := json.Unmarshal([]byte(respBody), &body); err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed parsing event webhook: %w", err),
		}
	}

	return &body, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// CreateEventWebhook creates a new event webhook and returns it with its assigned ID.
func (c *Client) CreateEventWebhook(ctx context.Context, webhook EventWebhook) (*EventWebhook, RequestError) {
	if webhook.URL == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrURLRequired,
		}
	}

	respBody, statusCode, err := c.Post(ctx, "POST", "/user/webhooks/event/settings", webhook)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed creating event webhook: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedCreatingEventWebhook, statusCode, respBody),
		}
	}

	return parseEventWebhook(respBody)
}

// ReadEventWebhook retrieves a specific event webhook by ID.
// SendGrid doesn't have a single-webhook GET endpoint, so we fetch all and filter.
func (c *Client) ReadEventWebhook(ctx context.Context, id string) (*EventWebhook, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrEventWebhookIDRequired,
		}
	}

	webhooks, _, reqErr := c.ListEventWebhooks(ctx)
	if reqErr.Err != nil {
		return nil, reqErr
	}

	for i := range webhooks {
		if webhooks[i].ID == id {
			return &webhooks[i], RequestError{StatusCode: http.StatusOK, Err: nil}
		}
	}

	return nil, RequestError{
		StatusCode: http.StatusNotFound,
		Err:        fmt.Errorf("event webhook not found: %s", id),
	}
}

// UpdateEventWebhook updates an existing event webhook by ID.
func (c *Client) UpdateEventWebhook(ctx context.Context, id string, webhook EventWebhook) (*EventWebhook, RequestError) {
	if id == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrEventWebhookIDRequired,
		}
	}

	if webhook.URL == "" {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrURLRequired,
		}
	}

	respBody, statusCode, err := c.Post(ctx, "PATCH", "/user/webhooks/event/settings/"+id, webhook)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed updating event webhook: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedUpdatingEventWebhook, statusCode, respBody),
		}
	}

	return parseEventWebhook(respBody)
}

// DeleteEventWebhook deletes an event webhook by ID.
func (c *Client) DeleteEventWebhook(ctx context.Context, id string) RequestError {
	if id == "" {
		return RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        ErrEventWebhookIDRequired,
		}
	}

	_, statusCode, err := c.Get(ctx, "DELETE", "/user/webhooks/event/settings/"+id)
	if err != nil {
		return RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed deleting event webhook: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices && statusCode != http.StatusNotFound {
		return RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d", ErrFailedDeletingEventWebhook, statusCode),
		}
	}

	return RequestError{StatusCode: http.StatusOK, Err: nil}
}

// ListEventWebhooks retrieves all event webhooks.
func (c *Client) ListEventWebhooks(ctx context.Context) ([]EventWebhook, int, RequestError) {
	respBody, statusCode, err := c.Get(ctx, "GET", "/user/webhooks/event/settings/all")
	if err != nil {
		return nil, 0, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("failed listing event webhooks: %w", err),
		}
	}

	list, parseErr := parseEventWebhookList(respBody)
	if parseErr.Err != nil {
		return nil, 0, parseErr
	}

	return list.Webhooks, list.MaxAllowed, RequestError{StatusCode: http.StatusOK, Err: nil}
}

// ConfigureEventWebhookSigning enables or disables webhook signing for a specific webhook.
func (c *Client) ConfigureEventWebhookSigning(ctx context.Context, id string, enabled bool) (*EventWebhookSigning, RequestError) {
	endpoint := "/user/webhooks/event/settings/signed"
	if id != "" {
		endpoint = "/user/webhooks/event/settings/" + id + "/signed"
	}

	respBody, statusCode, err := c.Post(ctx, "PATCH", endpoint, EventWebhookSigning{
		Enabled: enabled,
	})
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed configuring event webhook signing: %w", err),
		}
	}

	if statusCode >= http.StatusMultipleChoices {
		return nil, RequestError{
			StatusCode: statusCode,
			Err:        fmt.Errorf("%w, status: %d, response: %s", ErrFailedPatchingEventWebhook, statusCode, respBody),
		}
	}

	return parseEventWebhookSigning(respBody)
}

// ReadEventWebhookSigning retrieves the signing configuration for a specific webhook.
func (c *Client) ReadEventWebhookSigning(ctx context.Context, id string) (*EventWebhookSigning, RequestError) {
	endpoint := "/user/webhooks/event/settings/signed"
	if id != "" {
		endpoint = "/user/webhooks/event/settings/" + id + "/signed"
	}

	respBody, _, err := c.Get(ctx, "GET", endpoint)
	if err != nil {
		return nil, RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return parseEventWebhookSigning(respBody)
}
