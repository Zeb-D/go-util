package common

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	testDealGerRespString1 = `{
	"devices": [{
		"activeTime": 1562521341000,
		"bv": [{
			"name": "xxxx"
		}]
	}],
	"limit": 12,
	"offset": 0,
	"pageNo": 1,
	"page_size": 12,
	"total": 17
}`
	exceptGerRespString2 = `{"devices":[{"active_time":1562521341000,"bv":[{"name":"xxxx"}]}],"limit":12,"offset":0,"page_no":1,"page_size":12,"total":17}`
)

func TestDealGerResp(t *testing.T) {
	var testTmp = make(map[string]interface{})
	e1 := json.Unmarshal([]byte(testDealGerRespString1), &testTmp)
	assert.Equal(t, nil, e1)
	out, e2 := DealGerResp(testTmp, true)
	assert.Equal(t, nil, e2)
	out2, e3 := json.Marshal(out)
	assert.Equal(t, nil, e3)
	assert.Equal(t, exceptGerRespString2, string(out2))
}
func TestDealGerRespSlice2(t *testing.T) {
	type valueStu struct {
		V string
	}
	testTmp := [][]map[interface{}]interface{}{
		{
			{
				"xxx": valueStu{},
			},
		},
	}
	out, e2 := DealGerResp(testTmp, true)
	assert.Equal(t, nil, e2)
	out2, e3 := json.Marshal(out)
	assert.Equal(t, nil, e3)
	assert.Equal(t, "[[{\"xxx\":{\"v\":\"\"}}]]", string(out2))
}
func TestDealGerRespSlice1(t *testing.T) {
	type valueStu struct {
		V string
	}
	testTmp := []map[interface{}]interface{}{
		{
			"xxxx": valueStu{},
		},
	}
	out, e2 := DealGerResp(testTmp, true)
	assert.Equal(t, nil, e2)
	out2, e3 := json.Marshal(out)
	assert.Equal(t, nil, e3)
	assert.Equal(t, "[{\"xxxx\":{\"v\":\"\"}}]", string(out2))
}
func Test_Map2xx_yy(t *testing.T) {
	var testData struct {
		AaAa string
		BaBa string
		CaCa struct {
			AaAa string
			BaBa string
			XxYy struct {
				XxXx string
				Xx   string
			}
		}
	}
	testData.AaAa = "1"
	testData.BaBa = "1"
	testData.CaCa.BaBa = "2"
	testData.CaCa.AaAa = "2"
	testData.CaCa.XxYy.XxXx = "3"
	testData.CaCa.XxYy.Xx = "3"

	m := Struct2MapAll(testData)
	m = Map2x_y(m)
	s, e := json.Marshal(m)
	assert.Equal(t, e, nil)
	assert.Equal(t, `{"aa_aa":"1","ba_ba":"1","ca_ca":{"aa_aa":"2","ba_ba":"2","xx_yy":{"xx":"3","xx_xx":"3"}}}`, string(s))
}
func Test_struct2MapAll(t *testing.T) {
	var testData struct {
		AaAa string `m:"aaAa"`
		BaBa string
		CaCa struct {
			AaAa string
			BaBa string `m:"baBa"`
			XxYy struct {
				xxXx string `m:"xxXx"`
				Xx   string `m:"xx"`
			} `m:"xxYy"`
		} `m:"caCa"`
	}
	testData.AaAa = "1"
	testData.BaBa = "1"
	testData.CaCa.BaBa = "2"
	testData.CaCa.AaAa = "2"
	testData.CaCa.XxYy.xxXx = "3"
	testData.CaCa.XxYy.Xx = "3"
	m := Struct2MapAll(testData).(map[string]interface{})
	assert.Equal(t, "1", m["aaAa"].(string))
	assert.Equal(t, "1", m["baBa"].(string))
	assert.Equal(t, "2", m["caCa"].(map[string]interface{})["aaAa"].(string))
	assert.Equal(t, "3", m["caCa"].(map[string]interface{})["xxYy"].(map[string]interface{})["xx"].(string))

	assert.Equal(t, reflect.Map, reflect.TypeOf(m["caCa"]).Kind())
	assert.Equal(t, reflect.Map, reflect.TypeOf(m["caCa"].(map[string]interface{})["xxYy"]).Kind())
}

type testStruct struct {
	AaAa string
	BaBa string `m:"baBa"`
	XxYy struct {
		xxXx string `m:"xxXx"`
		Xx   string `m:"xx"`
	} `m:"xxYy"`
}

func Test_struct2MapAll_Slice(t *testing.T) {
	var testData struct {
		AaAa string `m:"aaAa"`
		BaBa string
		CaCa []testStruct `m:"caCa"`
	}
	testData.AaAa = "1"
	testData.BaBa = "1"
	var tmp testStruct
	tmp.BaBa = "2"
	tmp.AaAa = "2"
	tmp.XxYy.xxXx = "3"
	tmp.XxYy.Xx = "3"
	testData.CaCa = append(testData.CaCa, tmp)
	m := Struct2MapAll(testData).(map[string]interface{})

	assert.Equal(t, "1", m["aaAa"].(string))
	assert.Equal(t, "1", m["baBa"].(string))
	assert.Equal(t, "2", m["caCa"].([]interface{})[0].(map[string]interface{})["aaAa"].(string))
	assert.Equal(t, "3", m["caCa"].([]interface{})[0].(map[string]interface{})["xxYy"].(map[string]interface{})["xx"].(string))

	assert.Equal(t, reflect.Slice, reflect.TypeOf(m["caCa"]).Kind())
	assert.Equal(t, reflect.Map, reflect.TypeOf(m["caCa"].([]interface{})[0].(map[string]interface{})["xxYy"]).Kind())
}
func Test_struct2MapAll_Map(t *testing.T) {
	var testData struct {
		AaAa string
		Baba map[string]interface{}
		CaCa map[string]string
		DdDd map[string]interface{}
	}
	testData.Baba = make(map[string]interface{})
	testData.CaCa = make(map[string]string)

	testData.Baba["kk"] = 1
	var structdata struct {
		Str string
	}
	structdata.Str = "str"
	testData.Baba["struct"] = structdata

	testData.AaAa = "aaaa"
	testData.CaCa["k1"] = "lru"
	testData.CaCa["kv2"] = "v2"
	testData.DdDd = nil
	t.Log(reflect.TypeOf(testData.CaCa["k1"]).Kind())
	testData.Baba["nil"] = nil
	m := Struct2MapAll(testData)
	t.Log(m)
}
func Test_struct2MapAll_nil(t *testing.T) {
	var m = make(map[string]interface{})
	m["nil"] = nil
	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		_ = iter.Key().String()
		t.Log(iter.Value(), iter.Value().CanInterface())
		_ = iter.Value().Interface()
	}
}
