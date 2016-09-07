package cmd

import (
	"fmt"
	"log"

	"github.com/fanux/pbrain/common"
	"github.com/nats-io/nats"
)

type Message struct {
	Command chan common.Command
}

func (this *Message) MqCallBack(command *common.Command) {
	fmt.Println("get command: ", command)
	this.Command <- *command
}

func listenCommand(pluginName string, command chan common.Command) {
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
	command := make(chan common.Command)

	plugin.Start()

	go listenCommand(plugin.GetPluginName(), command)

	for {
		select {
		case c := <-command:
			plugin.Init() // update plugin info
			switch c.Command {
			case common.COMMAND_START_PLUGIN:
				plugin.Start()
			case common.COMMAND_STOP_PLUGIN:
				plugin.Stop()
			case common.COMMAND_ENABLE_STRATEGY:
				plugin.EnableStrategy(c.Body)
			case common.COMMAND_DISABLE_STRATEGY:
				plugin.DisableStrategy(c.Body)
			case common.COMMAND_UPDATE_DOCUMENT:
				plugin.UpdateDocument(c.Body)
			default:
				fun := common.GetCommandHandle(c.Command)
				if fun != nil {
					fun(c)
				} else {
					log.Printf("Unknown command [%s]\n", c.Command)
				}
			}
		}
		fmt.Println("start listen a new command")
	}
}
