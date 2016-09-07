package common

const (
	PLUGIN_ENABLE  = "enable"
	PLUGIN_DISABLE = "disable"

	STRATEGY_ENABLE  = "enable"
	STRATEGY_DISABLE = "disable"
)

type Plugin struct {
	Name        string
	Kind        string
	Status      string
	Description string
	Spec        string
	Manual      string
}

type Strategy struct {
	//witch plugin strategy belongs to
	PluginName string
	Name       string
	Status     string
	Document   string
}

/*
	define the app name and scale to witch number
*/
type AppScale struct {
	App    string
	Number int
}

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
