package decode

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Zeb-D/go-util/common"
	"github.com/Zeb-D/go-util/struct/Map"
	"strconv"
	"strings"
)

const NOSignKey = "sign"
const ProtocolKey = "protocol"
const TKey = "t"
const PVKey = "pv"
const DATAKey = "data"

const PV_BEGIN_INDEX = 0
const PV_END_INDEX = 3
const SIGN_BEGIN_INDEX = 3
const SIGN_END_INDEX = 19
const DATA_BEGIN_INDEX = 19

var NodeNilError = errors.New("node is nil")

func CreateSignInput(jsonObject map[string]string, localKey string) (string, error) {
	treeMap := Map.NewBSTMap()
	//	放到treeMap
	for key, value := range jsonObject {
		if len(value) > 0 && strings.EqualFold(key, NOSignKey) {
			treeMap.Add(key, value)
		}
	}
	var buffer bytes.Buffer
	for node := treeMap.RemoveMin(); !treeMap.IsEmpty(); {
		if node == nil {
			fmt.Println(jsonObject, " -->", treeMap, " ->node is nil")
			return "", NodeNilError
		} else {
			buffer.WriteString(node.Key().(string))
			buffer.WriteString("=")
			buffer.WriteString(node.Value().(string))
			buffer.WriteString("||")
		}
	}
	buffer.WriteString(localKey)
	return buffer.String(), nil
}

func CreateCloudSignInput_2_1(protocol, data interface{}, pv, sign, localKey string, t int64) (string, error) {
	params := Map.NewBSTMap()
	//	放到treeMap
	p, err := common.ToString(protocol)
	if err != nil {
		fmt.Println(protocol, " err->", err)
		return "", err
	}
	params.Add(ProtocolKey, p)
	params.Add(TKey, strconv.FormatInt(t, 10))
	params.Add(PVKey, pv)
	if data != nil {
		d, err := common.ToString(data)
		if err != nil {
			fmt.Println(data, " err->", err)
			return "", err
		}
		params.Add(DATAKey, d)
	}

	var buffer bytes.Buffer
	for node := params.RemoveMin(); !params.IsEmpty(); {
		if node == nil {
			fmt.Println(params, " ->node is nil")
			return "", NodeNilError
		} else {
			buffer.WriteString(node.Key().(string))
			buffer.WriteString("=")
			buffer.WriteString(node.Value().(string))
			buffer.WriteString("||")
		}
	}
	buffer.WriteString(localKey)
	return buffer.String(), nil
}

func SignInput_2_1(message, localKey string) string {
	pv := common.SubStringRange(message, PV_BEGIN_INDEX, PV_END_INDEX)
	data := common.SubString(message, DATA_BEGIN_INDEX)
	return fmt.Sprintf("data=%s||pv=%s||%s", data, pv, localKey)
}
