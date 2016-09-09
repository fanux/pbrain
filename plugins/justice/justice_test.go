package justice

import (
	"testing"

	"github.com/nats-io/nats"
)

const CHANNEL = "plugin_justice"

const (
	COMMAND_APP_METRICAL = "app-metrical"
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

	mq.Publish(CHANNEL, Command{, CHANNEL, ``})
	mq.Publish(CHANNEL, Command{, CHANNEL, ``})

	mq.Publish(CHANNEL, Command{, CHANNEL, ``})
	mq.Publish(CHANNEL, Command{, CHANNEL, ``})
	mq.Publish(CHANNEL, Command{, CHANNEL, ``})
}
