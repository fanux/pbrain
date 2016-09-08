package common

type CallBackFunc func(Command) error

var register map[string]CallBackFunc

func RegistCommand(command string, fun CallBackFunc) {
	register[command] = fun
}

func GetCommandHandle(command string) CallBackFunc {
	fun, ok := register[command]
	if ok {
		return fun
	}
	return nil
}

func UnRegistCommand(command string) {
	delete(register, command)
}

func init() {
	register = make(map[string]CallBackFunc)
}
