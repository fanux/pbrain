package common

import "log"

type CallBackFunc func(Command) error

var register map[string]CallBackFunc

func RegistCommand(command string, fun CallBackFunc) {
	register[command] = fun
	log.Printf("regist command: [%s]", command)
}

func GetCommandHandle(command string) CallBackFunc {
	fun, ok := register[command]
	if ok {
		log.Printf("get command handler: [%s]", command)
		return fun
	}
	return nil
}

func UnRegistCommand(command string) {
	delete(register, command)
	log.Printf("unregist command: [%s]", command)
}

func init() {
	register = make(map[string]CallBackFunc)
}
