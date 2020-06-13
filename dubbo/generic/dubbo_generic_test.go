package generic

import (
	"context"
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
	"testing"
	"time"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"

	_ "github.com/apache/dubbo-go/protocol/dubbo"
)

func init() {
	log.SetGlobalLog("dubbo-go", false)
}

//local test
//cd /go-util/dubbo/generic
//
//export CONF_CONSUMER_FILE_PATH=$PROJECT_PWD"/todo/dubbo_conf/client.yml"
//export APP_LOG_CONF_FILE=$PROJECT_PWD"/todo/dubbo_conf/log.yml"
//go test -v -run TestGenericInvoke
func TestGenericInvoke(t *testing.T) {
	//config.Load()
	var referenceConfig = config.ReferenceConfig{
		InterfaceName: "com.yd.scala.dubbo.client.IHelloService",
		Cluster:       "failover",
		Registry:      DefaultRegistry,
		Protocol:      dubbo.DUBBO,
		Generic:       true,
	}
	appName := referenceConfig.InterfaceName + UNDERLINE + ""
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
	fmt.Println("res:", resp)

	var invokeReq = InvokeReq{
		InterfaceName:  "com.yd.scala.dubbo.client.IHelloService",
		Method:         "sayHello",
		ParameterTypes: []string{"java.lang.String"},
	}

	resp, err = Invoke(invokeReq, []interface{}{"Yd1"})
	fmt.Println("resp: ", resp)

	resp, err = Invoke(invokeReq, []interface{}{"1Yd"})
	fmt.Println("resp: ", resp)
}
