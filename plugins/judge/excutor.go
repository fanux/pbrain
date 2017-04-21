package judge

import (
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cli"
	"github.com/docker/swarm/common"
	"github.com/fanux/pbrain/cmd"
)

//Excute is
func Excute(cmds []string) error {
	if len(cmds) == 0 {
		return nil
	}
	cli.SendRequest(ParseCommand(cmds), cmd.SwarmURL)
}

//ParseCommand is
func ParseCommand(cmds []string) (body common.ScaleAPI) {
	var flag string
	var err error

	for _, cmd := range cmds {
		item := common.ScaleItem{Labels: make(map[string]string)}
		ss := strings.Split(cmd, " ")
		logrus.Debugf("got cmd slice:%s", ss)
		//      fmt.Printf("got cmd slice:%s", ss)
		for _, s := range ss {
			//          fmt.Printf("cmd slice is:--%s--\n", s)
			switch s {
			case " ":
			case "":
			case "-f":
				flag = "f"
			case "-e":
				flag = "e"
			case "-n":
				flag = "n"
			case "-l":
				flag = "l"
			default:
				switch flag {
				case "f":
					item.Filters = append(item.Filters, s)
					flag = ""
				case "e":
					item.ENVs = append(item.ENVs, s)
					flag = ""
				case "n":
					item.Number, err = strconv.Atoi(s)
					if err != nil {
						logrus.Errorf("got error scale number : %s", s)
						return
					}
					flag = ""
				case "l":
					vSlice := strings.SplitN(s, "=", 2)
					if len(vSlice) == 2 {
						item.Labels[vSlice[0]] = vSlice[1]
					} else {
						logrus.Printf("invalid label: %s", s)
						return
					}
					flag = ""
				default:
				}
			}
		}

		body.Items = append(body.Items, item)
	}

	logrus.Debugf("got scale body: %s", body)

	return
}
