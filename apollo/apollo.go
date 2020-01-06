package apollo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//apollo配置中心
const (
	//当前应用对应apollo地址
	apolloAppInfoUrl string = "%s/configs/%s/default/APP_ALL"
)

var ApolloResponseNil error = errors.New("*ApolloResponse is nil")

var ApolloHandlers []ApolloHandler
var ap *ApolloResponse

//	InitApolloHandlers 监听者赶紧到这里来，按时触发一系列动作
func InitApolloHandlers(a ApolloHandler) {
	ApolloHandlers = append(ApolloHandlers, a)
}

//	ApolloHandler 当ApolloTimer 定时刷新的时候，主动触发此类操作
type ApolloHandler interface {
	Handler(omap, nMap map[string]string) error
}

//	ApolloTimer apollo client 本地刷新器
type ApolloTimer interface {
	Start(cycle time.Duration)
	Close()
}

func NewApolloListener(hostUrl string, appName string) ApolloTimer {
	return &apolloListener{
		hostUrl: hostUrl,
		appName: appName,
		close:   make(chan struct{}),
	}
}

type apolloListener struct {
	hostUrl string
	appName string
	close   chan struct{}
}

//	循环间隔时间不能低于3秒,应用层需要判断返回值err，err!=nil时才能更新配置
func (l *apolloListener) Start(cycle time.Duration) {
	if cycle < 3*time.Second {
		cycle = 100 * time.Second
	}
	go func() {
		// 做一些初始化事情
		var respCache []byte
		respTemp, err := GetAppInfoByApollo(l.hostUrl, l.appName)
		respCache = make([]byte, len(respTemp))
		copy(respCache, respTemp)
		if err != nil {
			log.Println("apolloListener Start has error ->", err)
		}
		new1 := &ApolloResponse{}
		err = json.Unmarshal(respCache, &new1)
		if err != nil {
			log.Println("apolloListener Start Init has error ->", err)
		}
		ap = new1

	LOOP:
		for {
			//执行一次后，开始循环监听变化
			select {
			case <-l.close:
				break LOOP
			case <-time.After(cycle):
				respTemp, err := GetAppInfoByApollo(l.hostUrl, l.appName)
				if err != nil {
					log.Println(time.Now(), "apolloListener Start has error ->", err)
					continue
				}
				if bytes.Equal(respCache, respTemp) {
					//相等不需要更新操作
					continue
				}
				//不相等，执行handle,并且将respCache覆盖
				//	下面的监听者 触发动作
				l.Handle(respCache, respTemp)
				respCache = make([]byte, len(respTemp))
				copy(respCache, respTemp)
			}
		}
		log.Println("amazing apollo listener exit now")
	}()
}

//	Handle 通知监听者去处理
func (l *apolloListener) Handle(respOld, respNew []byte) error {
	old := &ApolloResponse{}
	err := json.Unmarshal(respOld, &old)
	if err != nil && len(respOld) != 0 { //第一次可能为空
		log.Println(respOld, "apolloListener Handle has error ->", err)
		return err
	}
	new1 := &ApolloResponse{}
	err = json.Unmarshal(respNew, &new1)
	if err != nil {
		log.Println("apolloListener Handle has error ->", err)
		return err
	}
	//	监听者 处理
	if len(ApolloHandlers) != 0 {
		go func(omap, nMap map[string]string) {
			for _, ah := range ApolloHandlers {
				err := ah.Handler(omap, nMap)
				if err != nil {
					log.Println(ah, "ah Handle has error ->", err)
				}
			}
		}(old.Configurations, new1.Configurations)
	}
	ap = new1
	return nil
}

func (l *apolloListener) Close() {
	close(l.close)
}

// GetAppInfoFromApollo 调用方需要将v结构体提前定义好。如果不知道，就定义成map[string]interface{}
func GetAppInfoFromApollo(hostUrl string, appName string, v interface{}) error {
	bs, err := GetAppInfoByApollo(hostUrl, appName)
	log.Println("bs:->", string(bs), ":err:", err)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, &v)
	return err
}

//从配置中心获取appInfo
func GetAppInfoByApollo(hostUrl string, appName string) ([]byte, error) {
	appInfoUrl := fmt.Sprintf(apolloAppInfoUrl, hostUrl, appName)
	resp, err := http.Get(appInfoUrl)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}

//	ApolloResponse 返回结构体
type ApolloResponse struct {
	AppID          string            `json:"appId"`
	Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

func GetValueByKey(key string) (string, error) {
	if ap == nil {
		return "", ApolloResponseNil
	}
	return ap.Configurations[key], nil
}
