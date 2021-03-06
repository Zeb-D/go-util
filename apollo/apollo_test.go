package apollo

import (
	"fmt"
	"testing"
	"time"
)

type apolloResonse struct {
	AppID          string `json:"appId"`
	Cluster        string `json:"cluster"`
	NamespaceName  string `json:"namespaceName"`
	Configurations struct {
		RuntimeRegion string `json:"runtime.region"`
	} `json:"configurations"`
	ReleaseKey string `json:"releaseKey"`
}

type apollo404Resonse struct {
	Timestamp string `json:"timestamp"`
	Status    int8   `json:"status"`
	Error     string `json:"error"`
	Msg       string `json:"message"`
	Path      string `json:"path"`
}

func TestGetAppInfoFromApollo(t *testing.T) {
	hostUrl := "http://192.1.1.2:1111"
	appName := "test2"
	v := &apolloResonse{}
	//v := &map[string]interface{}
	err := GetAppInfoFromApollo(hostUrl, appName, v)
	if err != nil {
		t.Error(err)
	}
	println("rr:->", v.Configurations.RuntimeRegion)
	if get, want := v.AppID, "test2"; get != want {
		t.Errorf("except %v, but get %v", want, get)
	}
}

func TestApolloListener_Handle(t *testing.T) {
	hostUrl := "http://192.1.1.2:1111"
	appName := "test2"
	InitApolloHandlers(&apolloHandler1{})
	l := NewApolloListener(hostUrl, appName)
	l.Start(10 * time.Second)
	time.Sleep(10 * time.Second)
	s, _ := GetValueByKey("mysql.username")
	fmt.Println("value->", s)
}

type apolloHandler1 struct {
}

func (ah *apolloHandler1) Handler(omap, nMap map[string]string) error {
	key := "runtime.region"
	ov := omap[key]
	nv := nMap[key]
	fmt.Println("在处理某些事情", ov, nv)
	return nil
}
