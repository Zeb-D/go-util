package apollo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//apollo配置中心
const (
	//当前应用对应apollo地址
	apolloAppInfoUrl = "%s/configs/%s/default/APP_ALL"
)

// GetAppInfoFromApollo 调用方需要将v结构体提前定义好。如果不知道，就定义成map[string]interface{}
func GetAppInfoFromApollo(hostUrl string, appName string, v interface{}) error {
	bs, err := GetAppInfoByApollo(hostUrl, appName)
	fmt.Println("bs:->", string(bs), ":err:", err)
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
	fmt.Println("url:->", appInfoUrl, "resp:", resp, ":err:", err)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}
