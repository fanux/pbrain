package common

//each plugin need implementation this interface
type PluginInterface interface {
	Start() error
	Stop() error

	EnableStrategy(strategyName string) error
	DisableStrategy(strategyName string) error
	UpdateDocument(strategyName string) error

	// Base plugin implementation
	GetPluginName() string
}
