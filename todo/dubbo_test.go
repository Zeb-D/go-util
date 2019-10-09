package todo

import (
	"encoding/json"
	"fmt"
	hessian "github.com/apache/dubbo-go-hessian2"
	"reflect"
	"testing"
)

var loc = struct {
	Province string
	Load     string
}{
	Load: "Sz nanShan",
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
	//result, err := DubboClientStu.DoGer("com.xxx.provider.service.IHelloService", "", "sleep", []string{"com.tuya.provider.common.domain.User"}, reqData)
	//t.Log(result, err)
}
