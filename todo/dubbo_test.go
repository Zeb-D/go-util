package todo

import (
	"context"
	"encoding/json"
	"fmt"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
	"reflect"
	"testing"
	"time"
)

var loc = struct {
	Province string
	Load     string
}{
	Load: "Sz nanShan",
}

//本地测试
//cd /go-util/todo
//
//export CONF_CONSUMER_FILE_PATH=$PWD"/dubbo_conf/client.yml"
//export APP_LOG_CONF_FILE=$PWD"/dubbo_conf/log.yml"
//go test -v -run TestDubbo
func TestDubbo(t *testing.T) {
	var appName = "HelloProviderGer"
	var referenceConfig = config.ReferenceConfig{
		InterfaceName: "com.yd.scala.dubbo.client.IHelloService",
		Cluster:       "failover",
		Registry:      DefaultRegistry,
		Protocol:      dubbo.DUBBO,
		Generic:       true,
	}
	referenceConfig.GenericLoad(appName) //appName is the unique identification of RPCService

	time.Sleep(3 * time.Second)
	println("\n\n\nstart to generic invoke")
	resp, err := referenceConfig.GetRPCService().(*config.GenericService).Invoke(context.TODO(),
		[]interface{}{"sayHello",
			[]string{"java.lang.String"},
			[]interface{}{"Yd"}})
	if err != nil {
		panic(err)
	}
	println("res: %+v\n", resp)
	println("succ!")

}

func TestHessian(t *testing.T) {
	user := struct {
		Age        uint8
		Name       string
		Attchments map[interface{}]interface{}
	}{
		Age:  18,
		Name: "go-util",
	}
	user.Attchments = make(map[interface{}]interface{})
	user.Attchments["sex"] = "male"
	user.Attchments["locl"] = loc

	v, _ := json.Marshal(user)
	fmt.Printf("user:%s \n", string(v))

	cc := Struct2JavaMap(user)

	attchments := cc.(map[string]interface{})["attchments"]
	fmt.Printf("type:%s,val:%s \n", reflect.TypeOf(attchments), reflect.ValueOf(attchments))
	hessianData := []hessian.Object{cc}
	fmt.Printf("+v:%+v,#v:%#v,T:%T \n", hessianData, hessianData, hessianData)
	fmt.Printf("+v:%v", hessianData)

	//%+v 打印包括字段在内的实例的完整信息
	//%#v 打印包括字段和限定类型名称在内的实例的完整信息
	//%T 打印某个类型的完整说明
}
