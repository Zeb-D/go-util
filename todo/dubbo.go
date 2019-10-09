package todo

import (
	//dg "github.com/apache/dubbo-go/config"
	//"github.com/apache/dubbo-go/protocol/dubbo"
	"reflect"
	"strings"
)

const (
	Failover        = "failover"
	DefaultRegistry = "zk_1"
	UNDERLINE       = "_"
)

//todo dubbo 初始化如何接入apollo
//var GenericServicePool map[string]*dg.GenericService
//
////创建一个 接口级的dubbo泛化客户端
//func CreateDubboClient(interfaceName string, version string, protocol string) *dg.GenericService {
//	fmt.Println(dubbo.DUBBO)
//	key := interfaceName + UNDERLINE + version
//	var referenceConfig = dg.ReferenceConfig{
//		InterfaceName: interfaceName,
//		Cluster:       Failover,
//		Registry:      DefaultRegistry,
//		Protocol:      protocol,
//		Version:       version,
//		Generic:       true,
//		Retries:       1,
//	}
//	referenceConfig.GenericLoad(key)
//	time.Sleep(200 * time.Millisecond)
//	fmt.Println("create new DubboClient sleep 200 Millisecond")
//	clientService := referenceConfig.GetRPCService().(*dg.GenericService)
//	GenericServicePool[key] = clientService
//	return clientService
//}

//Struct2JavaMap golang的struct对象转化成Java的map，需要增加支持的范围
func Struct2JavaMap(obj interface{}) interface{} {
	if obj == nil {
		return obj
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Struct {
		result := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Struct {
				if v.Field(i).CanInterface() {
					setInMap(result, t.Field(i), Struct2JavaMap(v.Field(i).Interface()))
				}
			} else if v.Field(i).Kind() == reflect.Slice || v.Field(i).Kind() == reflect.Map {
				if v.Field(i).CanInterface() {
					setInMap(result, t.Field(i), Struct2JavaMap(v.Field(i).Interface()))
				}
			} else {
				if v.Field(i).CanInterface() {
					setInMap(result, t.Field(i), v.Field(i).Interface())
				}
			}
		}
		return result
	} else if t.Kind() == reflect.Slice {
		var newTemps = make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			newTemp := Struct2JavaMap(v.Index(i).Interface())
			newTemps = append(newTemps, newTemp)
		}
		return newTemps
	} else if t.Kind() == reflect.Map {
		var newTempMap = make(map[interface{}]interface{}, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			if !iter.Key().CanInterface() {
				continue
			}
			if !iter.Value().CanInterface() {
				continue
			}
			mapK := Struct2JavaMap(iter.Key().Interface())
			mapV := iter.Value().Interface()
			newTempMap[mapK] = Struct2JavaMap(mapV)
		}
		return newTempMap
	} else {
		return obj
	}
}

func setInMap(m map[string]interface{}, structField reflect.StructField, value interface{}) (result map[string]interface{}) {
	result = m
	if tagName := structField.Tag.Get("m"); tagName == "" {
		result[headerAtoa(structField.Name)] = value
	} else {
		result[tagName] = value
	}
	return
}

func headerAtoa(a string) (b string) {
	b = strings.ToLower(a[:1]) + a[1:]
	return
}
