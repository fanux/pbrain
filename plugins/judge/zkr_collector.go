package judge

/*
 *连接zkr的thrift
 *调用thrift的OnCommand接口
 *获取引擎实例数
 */

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"

	"git.7s.iflytek.com/sivs_team/sivs_tool.git/sivs_cli/thrift/tes"
	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	zkrServer = []string{"192.168.86.60:9095"} //可配置多台zkr
)

//ConnectThrift 连接zkr的thrift
func ConnectThrift(ipPort string) (transport *thrift.TSocket, client *tes.TesClient, err error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err = thrift.NewTSocket(ipPort)

	if err != nil {
		return nil, nil, err
	}

	useTransport := transportFactory.GetTransport(transport)
	client = tes.NewTesClientFactory(useTransport, protocolFactory)

	return transport, client, nil
}

//GetInst 获取不同的路由服务实例数信息(如：iat、tts....)
func GetInst(router string) (r string, err error) {
	var (
		succ      = false
		client    *tes.TesClient
		transport *thrift.TSocket
	)

	for _, ipPort := range zkrServer {
		transport, client, err = ConnectThrift(ipPort)
		if err != nil {
			fmt.Printf("connect thrift error\n")
			continue
		}
		if err = transport.Open(); err != nil {
			fmt.Printf("transport open error\n")
			continue
		}
		succ = true
		defer transport.Close()
		break
	}

	if !succ {
		fmt.Printf("connect thrift error\n")
		return
	}

	var param = "status  -router  " + router //传递给zkr接口的参数，需用空格分隔开
	tesMsg := tes.NewTesNotifyMsg()

	tesMsg.Data = []byte(param)
	r, _ = client.OnCommand(tesMsg)

	return r, nil

}

//EngineInst 引擎实例数的信息
type EngineInst struct {
	TotalInst int32 //总实例数
	IdleInst  int32 //空闲实例数
}

//EngineInfo  引擎信息
type EngineInfo struct {
	Svc  string //服务器
	Inst EngineInst
}

//EngineInfoSet 路由服务的所有引擎的信息
type EngineInfoSet struct {
	Error     string //返回错误
	EngineSet []EngineInfo
}

//ZkrCollector is
type ZkrCollector struct {
	Metric int
}

//Collect is
func (c *ZkrCollector) Collect() int {
	var router = "iat" //路由类型
	var totalInst int
	var idleInst int

	r, _ := GetInst(router)
	if r == "" {
		fmt.Printf("router %s getInst error\n", router)
		return
	}
	fmt.Printf("router %s getInst info is: %s\n", router, r)

	var engineInfoSet EngineInfoSet
	err := json.Unmarshal([]byte(r), &engineInfoSet)
	if err != nil {
		fmt.Printf("json Unmarshal error")
		return
	}

	if engineInfoSet.Error != "" {
		fmt.Printf("call GetInst error")
		return
	}

	//获取引擎的总实例数和空闲实例数，可做相应的处理
	for _, engine := range engineInfoSet.EngineSet {
		var svc = engine.Svc //引擎（ip:port）
		fmt.Printf("svc is %s\n", svc)

		totalInst += engine.Inst.TotalInst
		idleInst += engine.Inst.IdleInst
		fmt.Printf("totalInst is %d,idleInst is %d\n", totalInst, idleInst)
	}

	c.Metric = (IdleInst * 100) / totalInst

	logrus.Debugf("got zkr metric : %d")

	return c.Metric
}
