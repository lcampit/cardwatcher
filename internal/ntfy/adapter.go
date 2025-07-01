package ntfy

import "context"

type NtfyAdapter interface {
	Notify(ctx context.Context, message string) error
}

type ntfyAdapter struct {
	ntfyHost string
	ntfyPort string
}

func NewNtfyAdapter(host, port string) NtfyAdapter {
	return &ntfyAdapter{
		ntfyHost: host,
		ntfyPort: port,
	}
}
