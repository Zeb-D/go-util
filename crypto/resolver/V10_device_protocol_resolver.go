package resolver

import (
	"encoding/json"
	"errors"
	"github.com/Zeb-D/go-util/crypto"
	"github.com/Zeb-D/go-util/crypto/decode"
	"log"
	"reflect"
	"strings"
)

var SignFailed error = errors.New("sign not equals")
var UnsupportedType = errors.New("unsupported type params type")

func init() {
	InitResolver(&V10DeviceProtocolResolver{})
}

type V10DeviceProtocolResolver struct {
}

func (r *V10DeviceProtocolResolver) ProtocolVersion() string {
	return V_1_0
}

func (r *V10DeviceProtocolResolver) Decode(protocolData interface{}, devId, localKey, protocolVersion string) (string, error) {
	var decodedProtocolData map[string]string
	switch protocolData.(type) {
	case map[string]string: //数据上报是JSON格式
		decodedProtocolData = protocolData.(map[string]string)
	case string: //数据下发是字符串格式
		decodedProtocolData = make(map[string]string)
		err := json.Unmarshal([]byte(protocolData.(string)), &decodedProtocolData)
		if err != nil {
			log.Println("V10DeviceProtocolResolver Decode error ", err)
			return "", err
		}
	default:
		log.Println("unsupported type params type:%v", reflect.TypeOf(protocolData).Kind().String())
		return "", UnsupportedType
	}

	comparedSign := decodedProtocolData[KEY_SIGN]
	if len(comparedSign) == 0 && HARDWARE_UPGRADE_PROGRESS == decodedProtocolData[KEY_PROTOCOL] {
		//设备固件升级消息没有进行加密处理,故直接返回
		bs, err := json.Marshal(decodedProtocolData)
		if err != nil {
			log.Println("V10DeviceProtocolResolver json.Marshal(decodedProtocolData) error ", err, " decodedProtocolData:", decodedProtocolData)
		}
		return string(bs), err
	}

	signInput, err := decode.CreateSignInput(decodedProtocolData, localKey)
	if err != nil {
		log.Println("V10DeviceProtocolResolver decode.CreateSignInput error ", err, " decodedProtocolData:", decodedProtocolData)
		return "", err
	}

	sign := crypto.MD5(signInput)
	if !strings.EqualFold(comparedSign, sign) {
		log.Println("sign not equals,compared=", comparedSign,
			",computed=", sign,
			",signInput=", signInput,
			",msg=", decodedProtocolData)
		return "", SignFailed
	}

	bs, err := json.Marshal(decodedProtocolData)
	if err != nil {
		log.Println("V10DeviceProtocolResolver json.Marshal(decodedProtocolData) error ", err, " decodedProtocolData:", decodedProtocolData)
		return "", err
	}
	return string(bs), err
}
