package todo

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

func Yaml2Struct(confFile string, in interface{}) (err error) {
	if path.Ext(confFile) != ".yml" {
		return err
	}
	bs, err := ioutil.ReadFile(confFile)
	fmt.Errorf("Yaml2Struct Marshal path:%s,err:%s \n", confFile, err)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bs, in)
	fmt.Errorf("Yaml2Struct bs:%s,err:%s \n", string(bs), err)
	if err != nil {
		return err
	}
	return err
}
