package processor

import (
	"fmt"
	"github.com/Zeb-D/go-util/crypto/resolver"
)

type MeshDpReportProcessor struct {
}

func (p *MeshDpReportProcessor) DeviceType() string {
	return DEVICETYPE_MESH
}

func (p *MeshDpReportProcessor) Process(topicDevId string, Data map[string]string) (bool, error) {
	timestamp := Data[resolver.KEY_TIME_STAMP]
	payload := Data[resolver.KEY_PAYLOAD]
	//localKey :=Data[resolver.Loa]
	fmt.Println(timestamp, payload)
	return true, nil
}
