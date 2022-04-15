package crypto

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"reflect"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type ProtocolVO struct {
	Protocol int    `json:"protocol"`
	Data     []byte `json:"Data"`
	GwId     string `json:"gwId"`
	Pv       string `json:"pv"`
	T        int64  `json:"t"`
	Sign     string `json:"sign"`
	S        uint32 `json:"s,omitempty"`
	From     uint32 `json:"-"`
}

func NewProtocolVO22() ProtocolVO {
	return ProtocolVO{Pv: "2.2"}
}

func (pv *ProtocolVO) EnCrypt(data interface{}, localKey []byte) (decryptData []byte, err error) {
	if pv.Pv == "2.0" {
		return pv.EnCrypt20(data, localKey)
	} else if pv.Pv == "2.1" {
		return pv.EnCrypt21(data, localKey)
	} else if pv.Pv == "2.2" {
		return pv.EnCrypt22(data, localKey, pv.S, pv.From)
	}
	err = errors.New("pv not support")
	return
}
func DeCrypt(pv string, payload, localKey []byte) (decryptData []byte, err error) {
	if pv == "2.0" {
		return DeCrypt20(payload, localKey)
	} else if pv == "2.1" {
		return DeCrypt21(payload, localKey)
	} else if pv == "2.2" {
		return DeCrypt22(payload, localKey)
	}
	return
}

func DeCrypt22(payload, localKey []byte) (decryptData []byte, err error) {
	if string(payload[:3]) != "2.2" {
		err = errors.New("proto version not match")
		return
	}
	sign := binary.BigEndian.Uint32(payload[3:7])

	svrSign := CheckSum(payload[7:])
	if sign != svrSign {
		err = fmt.Errorf("sign not match %v %v", sign, svrSign)
		return
	}
	decryptData, err = AesECBDecrypt(payload[15:], localKey)
	return
}

func (pv *ProtocolVO) EnCrypt22(data interface{}, localKey []byte, seq uint32, from uint32) (encryptData []byte, err error) {
	msgPayload, err := json.Marshal(data)
	if err != nil {
		return
	}

	encryptPayload, err := AesECBEncrypt(msgPayload, localKey)
	if err != nil {
		return
	}

	bytesLen := 8 + len(encryptPayload)
	signInput := make([]byte, bytesLen)
	copy(signInput[0:4], IntToBytes(seq))
	copy(signInput[4:8], IntToBytes(from))
	copy(signInput[8:], encryptPayload)

	sign := CheckSum(signInput)
	bytesLen += 7
	encryptData = make([]byte, bytesLen)
	copy(encryptData[:3], []byte("2.2"))
	copy(encryptData[3:7], IntToBytes(sign))
	copy(encryptData[7:], signInput)
	return
}

func CheckSum(signInput []byte) uint32 {
	// crc32q := crc32.MakeTable(0xffffffff)
	return crc32.ChecksumIEEE(signInput)
}

func DeCrypt21(payload, localKey []byte) (decryptData []byte, err error) {
	if string(payload[:3]) != "2.1" {
		err = errors.New("proto version not match")
		return
	}
	sign := payload[3:19]
	signInput := fmt.Sprintf("Data=%s||pv=%s||%s", string(payload[19:]), "2.1", string(localKey))
	signCompute := MakeSign(signInput, 8, 24)
	if string(sign) != signCompute {
		err = errors.New("sign error")
		return
	}
	encryptData, err := base64.StdEncoding.DecodeString(string(payload[19:]))
	if err != nil {
		return
	}
	decryptData, err = AesECBDecrypt(encryptData, localKey)
	return
}

func (pv *ProtocolVO) EnCrypt21(data interface{}, localKey []byte) (encryptData []byte, err error) {
	pv.Data, err = json.Marshal(data)
	if err != nil {
		return
	}
	pv.Data, err = AesECBEncrypt(pv.Data, localKey)
	if err != nil {
		return
	}
	Base64Str := base64.StdEncoding.EncodeToString(pv.Data)
	signInput := fmt.Sprintf("Data=%s||pv=%s||%s", string(Base64Str), "2.1", string(localKey))
	sign := MakeSign(signInput, 8, 24)
	encryptData = []byte("2.1" + sign + string(Base64Str))
	return
}

func DeCrypt20(Payload []byte, localKey []byte) (decryptData []byte, err error) {
	pv := ProtocolVO{}
	err = json.Unmarshal(Payload, &pv)
	if err != nil {
		return
	}
	mapData := map[string]interface{}{
		"t":        pv.T,
		"pv":       pv.Pv,
		"protocol": pv.Protocol,
		"Data":     string(pv.Data),
	}
	signFormat := formatSign20(mapData, localKey)
	if pv.Sign != signFormat {
		err = fmt.Errorf("sign not match %v %v", pv.Sign, signFormat)
		return
	}
	var encryptData interface{}
	var ok bool
	if encryptData, ok = mapData["Data"]; !ok {
		err = errors.New("payload no have Data")
		return
	}
	decryptData, err = AesECBDecrypt([]byte(encryptData.(string)), localKey)
	return
}

func (pv *ProtocolVO) EnCrypt20(data interface{}, localKey []byte) (encryptData []byte, err error) {
	pv.Data, err = json.Marshal(data)
	if err != nil {
		return
	}
	pv.Data, err = AesECBEncrypt(pv.Data, localKey)
	if err != nil {
		return
	}
	pv.Sign, err = pv.getSign20(localKey)
	if err != nil {
		return
	}
	encryptData, err = json.Marshal(pv)

	return
}

func (pv *ProtocolVO) getSign20(localKey []byte) (sign string, err error) {
	mapData := map[string]interface{}{}
	if pv.Protocol == 0 {
		err = errors.New("protocol is 0")
		return
	}
	mapData["protocol"] = pv.Protocol
	if pv.T == 0 {
		pv.T = time.Now().Unix()
	}
	mapData["t"] = pv.T
	if pv.Pv == "" || len(pv.Pv) == 0 {
		pv.Pv = "2.0"
	}
	mapData["pv"] = pv.Pv
	if pv.GwId != "" || len(pv.GwId) != 0 {
		mapData["gwId"] = pv.GwId
	}

	if pv.Data != nil {
		mapData["Data"] = string(pv.Data)
	}
	sign = formatSign20(mapData, localKey)
	return
}
func formatSign20(mapData map[string]interface{}, localKey []byte) (sign string) {
	keys := make([]string, 0)
	for k := range mapData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	formatStr := "%s"
	formatVals := []interface{}{string(localKey)}
	for _, k := range keys {
		if k != "sign" {
			if mapData[k] != nil || len(k) != 0 {
				formatStr = k + "=%v||" + formatStr
				formatVals = append(formatVals, mapData[k])
			}
		}
	}
	ReverseSlice(formatVals)
	sign = Md5(fmt.Sprintf(formatStr, formatVals...))

	return
}

// MakeSign 可根据 secKey 生成 mqtt password
func MakeSign(signInput string, Start int, End int) (sign string) {
	sign = Md5(signInput)
	return sign[Start:End]
}

func Md5(source string) (md5str string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(source))
	cipherStr := md5Ctx.Sum(nil)
	md5str = hex.EncodeToString(cipherStr)
	return
}

func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func IntToBytes(n uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, &n)
	return bytesBuffer.Bytes()
}

//func SHA265Sign(msg []byte, secret []byte) string {
//	h := hmac.New(sha256.New, secret)
//	h.Write(msg)
//	return hex.EncodeToString(h.Sum(nil))
//}
