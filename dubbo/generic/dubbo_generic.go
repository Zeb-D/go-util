package generic

import (
	"context"
	"github.com/Zeb-D/go-util/apollo"
	"github.com/Zeb-D/go-util/dubbo/common"
	"github.com/Zeb-D/go-util/log"
	"github.com/apache/dubbo-go/protocol/dubbo"
	"sync"
	"time"
)

import (
	dG "github.com/apache/dubbo-go/config"
)

var (
	servicePool = genericServicePool{
		genericServicePool: make(map[string]*dG.GenericService),
	}
)

const (
	Retries         = "0"
	Failover        = "failover"
	DefaultRegistry = "yd_zk_1"
	UNDERLINE       = "_"
)

type genericServicePool struct {
	rwLock             sync.RWMutex
	genericServicePool map[string]*dG.GenericService
}

func Init() {
	dConfig := dG.GetConsumerConfig()
	for i := range dConfig.Registries {
		dConfig.Registries[i].Address, _ = apollo.GetValueByKey("ZKAddr")
	}
	log.Info("dubboConfig", log.Any("dubboConfig", dConfig))
	dG.SetConsumerConfig(dConfig)
	dG.Load()
	log.Info("dubbo init finish")
}

type InvokeReq struct {
	InterfaceName  string
	Version        string
	Method         string
	ParameterTypes []string
}

func (p *genericServicePool) Invoke(req InvokeReq, reqData []interface{}) (resp interface{}, err error) {
	return InvokeGeneric(req.InterfaceName, req.Version, req.Method, req.ParameterTypes, reqData)
}

func Invoke(req InvokeReq, reqData []interface{}) (resp interface{}, err error) {
	return servicePool.Invoke(req, reqData)
}

func InvokeGeneric(interfaceName string, version string, method string,
	parameterTypes []string, reqData []interface{}) (resp interface{}, err error) {

	service, _ := servicePool.getGenericService(interfaceName, version)
	resp, err = service.Invoke(context.TODO(),
		[]interface{}{method, parameterTypes, reqData})
	if err != nil {
		log.Warn("invoke ", log.Any("req", interfaceName),
			log.Any("version", version), log.Any("method", method),
			log.Any("parameterTypes", parameterTypes), log.Any("reqData", reqData),
			log.ErrorField(err))
		return
	}
	resp, err = common.DealGerResp(resp, false)
	if err != nil {
		return
	}
	return
}

func (p *genericServicePool) add(appName string, c *dG.GenericService) {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	p.genericServicePool[appName] = c
	log.Info("add interface", log.Any("len", len(p.genericServicePool)))
}

func (p *genericServicePool) get(appName string) *dG.GenericService {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.genericServicePool[appName]
}

func (p *genericServicePool) check(pattern string) bool {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	if _, ok := p.genericServicePool[pattern]; ok {
		return true
	} else {
		return false
	}
}

func (p *genericServicePool) createGenericService(interfaceName, version string) *dG.GenericService {
	key := interfaceName + UNDERLINE + version
	var referenceConfig = dG.ReferenceConfig{
		InterfaceName: interfaceName,
		Cluster:       Failover,
		Registry:      DefaultRegistry,
		Protocol:      dubbo.DUBBO,
		Version:       version,
		Generic:       true,
		Retries:       Retries,
	}
	//referenceConfig.Load(interfaceName) //appName是GetService的唯一标识不可缺少
	referenceConfig.GenericLoad(interfaceName)
	time.Sleep(200 * time.Millisecond) //第一次生成客户端时等待zk连接，目前等待200毫秒 如果需要不等待 请在启动时init
	genericService := referenceConfig.GetRPCService().(*dG.GenericService)
	p.add(key, genericService)
	return genericService
}

func (p *genericServicePool) getGenericService(interfaceName, version string) (*dG.GenericService, int) {
	key := interfaceName + UNDERLINE + version
	if p.check(key) {
		return p.get(key), 0
	} else {
		service := p.createGenericService(interfaceName, version)
		return service, 1
	}
}
