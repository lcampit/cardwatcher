package ntfy

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

func (a *ntfyAdapter) Notify(ctx context.Context, message string) error {
	endpoint := fmt.Sprintf("%s/%s", a.ntfyHost, "cardwatcher-alerts-lc")
	err := requests.URL(endpoint).Post().BodyBytes([]byte(message)).Fetch(ctx)
	if err != nil {
		return fmt.Errorf("error in ntfy adapter when sending message %s: %w", message, err)
	}
	return nil
}
