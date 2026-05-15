package ntfy

import (
	"context"
	"fmt"
)

func (a *ntfyAdapter) Notify(ctx context.Context, message string) error {
	_, err := a.client.R().SetBody(message).Post(a.topic)
	if err != nil {
		return fmt.Errorf("sending notification message %s: %w", message, err)
	}
	return nil
}
