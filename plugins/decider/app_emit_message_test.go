package decider

import (
	"testing"

	"github.com/nats-io/nats"
)

/*
const (
	COMMAND_APP_METRICAL = "app-metrical"
)
*/

var MESSAGE = `{
    "App":"172.20.1.128:5000/httpd:2.4",
    "Number":-6,
    "MinNum":2,
    "Hosts":[
        "172.20.1.106",
        "172.20.1.119",
        "172.20.1.107",
        "172.20.1.105",
        "172.20.1.104",
        "172.20.1.118",
        "172.20.1.108"
    ],
    "ScaleUp":"172.20.1.128:5000/nginx:latest"
}`

func TestEmitCommandsWithHost(t *testing.T) {
	nc, _ := nats.Connect(nats.DefaultURL)
	mq, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer mq.Close()

	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"ats", "Metrical":90}`})
	mq.Publish(CHANNEL, Command{COMMAND_APP_SCALE, CHANNEL, MESSAGE})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"ats", "Metrical":30}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"ats", "Metrical":90}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"hadoop", "Metrical":20}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"hadoop", "Metrical":40}`})
	//mq.Publish(CHANNEL, Command{COMMAND_APP_METRICAL, CHANNEL, `{"App":"redis", "Metrical":70}`})
}
