package resolver

import (
	"fmt"
	"github.com/Zeb-D/go-util/crypto"
	"log"
)

//const	VERSION_LENGTH	= 3
//const	SIGN_LENGTH		= 4
//const	SEQ_LENGTH		= 4
//const	FROM_LENGTH		= 4

type V22HexDeviceProtocolResolver struct {
}

func (r *V22HexDeviceProtocolResolver) ProtocolVersion() string {
	return V_2_2
}

func (r *V22HexDeviceProtocolResolver) Decode(protocolData interface{}, devId, localKey, protocolVersion string) (string, error) {
	pd, ok := protocolData.(string)
	if !ok {
		return "", UnsupportedType
	}
	if IsJSONMessage(pd) {
		//设备固件升级消息没有进行加密处理,故直接返回
		return pd, nil
	}
	decodedProtocolData, err := DecodeString(pd)
	if err != nil {
		log.Println(pd, " hex.DecodeString(pd) has error:", err)
		return "", err
	}
	// 获取到sign数据
	comparedSign := make([]int8, SIGN_LENGTH)
	copy(comparedSign[0:SIGN_LENGTH], decodedProtocolData[VERSION_LENGTH:VERSION_LENGTH+SIGN_LENGTH])

	signInputLength := len(decodedProtocolData) - VERSION_LENGTH - SIGN_LENGTH
	signInput := make([]int8, signInputLength)
	copy(signInput,
		decodedProtocolData[VERSION_LENGTH+SIGN_LENGTH:VERSION_LENGTH+SIGN_LENGTH+signInputLength])

	sign := Checksum(signInput)
	if equals(sign, comparedSign) {
		fmt.Println(sign, " not equals comparedSign:", comparedSign)
		return "", SignFailed
	}

	dataLength := len(decodedProtocolData) - VERSION_LENGTH - SIGN_LENGTH - SEQ_LENGTH - FROM_LENGTH
	data := make([]int8, dataLength)
	copy(data, decodedProtocolData[VERSION_LENGTH+SIGN_LENGTH+SEQ_LENGTH+FROM_LENGTH:])

	decryptMessageData := crypto.AESDecryptByECB(newStringUtf8(data), localKey)

	return decryptMessageData, nil
}
