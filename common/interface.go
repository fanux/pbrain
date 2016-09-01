package common

//each plugin need implementation this interface
type PluginInterface interface {
	Start() error
	Stop() error
	EnableStrategy(strategyName string)
	DisableStrategy(strategyName string) error
	UpdateDocument(strategyName string) error
}
