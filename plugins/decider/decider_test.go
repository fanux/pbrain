package decider

import (
	"testing"

	"github.com/nats-io/nats"
)

const CHANNEL = "plugin_decider"

/*
const (
	COMMAND_APP_METRICAL = "app-metrical"
)
*/

type Command struct {
	Command string
	Channel string // each plugin subscribe it plugin name, plugin name is channel
	Body    string
}

func TestEmitCommands(t *testing.T) {
	nc, _ := nats.Connect(nats.DefaultURL)
	mq, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer mq.Close()

	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"ats", "Metrical":90}`})
	mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"172.20.1.128:5000/nginx:latest", "Metrical":90}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"ats", "Metrical":30}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"ats", "Metrical":90}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"hadoop", "Metrical":20}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"hadoop", "Metrical":40}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"redis", "Metrical":70}`})
}
