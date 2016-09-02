package cmd

import (
	"fmt"
	"log"

	"github.com/fanux/pbrain/common"
	"github.com/nats-io/nats"
)

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

type Message struct {
	Command chan Command
}

func (this *Message) MqCallBack(command *Command) {
	fmt.Println("get command: ", command)
	this.Command <- *command
}

func listenCommand(pluginName string, command chan Command) {
	// subscribe plugin channel

	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer c.Close()

	message := Message{Command: command}
	c.Subscribe(pluginName, message.MqCallBack)

	fmt.Println("subscribe channel: ", pluginName)

	// stop gorutine exit
	select {}
}

func RunPlugin(plugin common.PluginInterface) {
	command := make(chan Command)

	plugin.Start()

	go listenCommand(plugin.GetPluginName(), command)

	for {
		select {
		case c := <-command:
			plugin.Init() // update plugin info
			switch c.Command {
			case COMMAND_START_PLUGIN:
				plugin.Start()
			case COMMAND_STOP_PLUGIN:
				plugin.Stop()
			case COMMAND_ENABLE_STRATEGY:
				plugin.EnableStrategy(c.Body)
			case COMMAND_DISABLE_STRATEGY:
				plugin.DisableStrategy(c.Body)
			case COMMAND_UPDATE_DOCUMENT:
				plugin.UpdateDocument(c.Body)
			default:
				log.Printf("Unknown command [%s]\n", c.Command)
			}
		}
		fmt.Println("start listen a new command")
	}
}
