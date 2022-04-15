package resolver

import (
	"github.com/Zeb-D/go-util/common"
	"github.com/Zeb-D/go-util/crypto"
	"github.com/Zeb-D/go-util/crypto/decode"
	"log"
	"strings"
)

func init() {
	InitResolver(&V21DeviceProtocolResolver{})
}

type V21DeviceProtocolResolver struct {
}

func (r *V21DeviceProtocolResolver) ProtocolVersion() string {
	return V_2_1
}

func (r *V21DeviceProtocolResolver) Decode(protocolData interface{}, devId, localKey, protocolVersion string) (string, error) {
	decodedProtocolData, ok := protocolData.(string)
	if !ok {
		return "", UnsupportedType
	}
	if IsJSONMessage(decodedProtocolData) {
		//设备固件升级消息没有进行加密处理,故直接返回
		return decodedProtocolData, nil
	}
	signInput := decode.SignInput_2_1(decodedProtocolData, localKey)
	sign := crypto.Get16BitMD5(signInput)
	comparedSign := common.SubStringRange(decodedProtocolData, decode.SIGN_BEGIN_INDEX, decode.SIGN_END_INDEX)

	if !strings.EqualFold(comparedSign, sign) {
		log.Println("sign not equals,compared=", comparedSign,
			",computed=", sign,
			",signInput=", signInput,
			",msg=", decodedProtocolData)
		return "", SignFailed
	}

	decryptMessageData := crypto.AESDecryptByECB(
		common.SubString(decodedProtocolData, decode.DATA_BEGIN_INDEX), localKey)

	return decryptMessageData, nil
}
