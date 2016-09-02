package common

import (
	"testing"

	"github.com/nats-io/nats"
)

const CHANNEL = "plugin_pipeline"

const (
	COMMAND_START_PLUGIN     = "start-plugin"
	COMMAND_STOP_PLUGIN      = "stop-plugin"
	COMMAND_ENABLE_STRATEGY  = "enable-strategy"
	COMMAND_DISABLE_STRATEGY = "disable-strategy"
	COMMAND_UPDATE_DOCUMENT  = "update-document"
)

type Command struct {
	Command string
	Channel string // each plugin subscribe it plugin name, plugin name is channel
	Body    string
}

func TestEmitCommands(t *testing.T) {
	nc, _ := nats.Connect(nats.DefaultURL)
	mq, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer mq.Close()

	mq.Publish(CHANNEL, Command{COMMAND_START_PLUGIN, CHANNEL, ""})
	mq.Publish(CHANNEL, Command{COMMAND_STOP_PLUGIN, CHANNEL, ""})

	mq.Publish(CHANNEL, Command{COMMAND_ENABLE_STRATEGY, CHANNEL, "strategy-name"})
	mq.Publish(CHANNEL, Command{COMMAND_DISABLE_STRATEGY, CHANNEL, "strategy-name"})
	mq.Publish(CHANNEL, Command{COMMAND_UPDATE_DOCUMENT, CHANNEL, "strategy-name"})
}
