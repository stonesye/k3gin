package ws

import "context"

type options struct {
	configFile string
}

func WithConfigFile(configFile string) func(*options) {
	return func(o *options) {
		o.configFile = configFile
	}
}

// Run WebSocket Server
func Run(ctx context.Context, opts ...func(*options)) error {
	return nil
}
