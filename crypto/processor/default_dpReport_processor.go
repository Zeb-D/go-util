package processor

import (
	"github.com/Zeb-D/go-util/common"
	"github.com/Zeb-D/go-util/log"
)

func init() {
	InitProcessor(&DefaultDpReportProcessor{})
}

type DefaultDpReportProcessor struct {
}

func (p *DefaultDpReportProcessor) DeviceType() string {
	return DeviceTypeDefault
}

func (p *DefaultDpReportProcessor) Process(topicDevId string, Data map[string]string) (bool, error) {
	data, _ := common.ToString(Data)
	log.Warn("illegal gateway id=" + topicDevId + ",msg=" + data)
	return true, nil
}
