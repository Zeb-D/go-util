package processor

import (
	"errors"
	"regexp"
	"strings"
)

const DeviceTypeDefault = "Default"

var DpReportProcessors map[string]IDpReportProcessor = make(map[string]IDpReportProcessor)

//	InitProcessor 实现类自己主动注册进来
func InitProcessor(r IDpReportProcessor) error {
	if r == nil {
		return errors.New("IDpReportProcessor is nil")
	}
	DpReportProcessors[r.DeviceType()] = r
	return nil
}

func GetDpReportProcessor(protocolVersion string) IDpReportProcessor {
	pr := DpReportProcessors[protocolVersion]
	if pr == nil {
		pr = DpReportProcessors[DeviceTypeDefault]
	}
	return pr
}

type IDpReportProcessor interface {
	DeviceType() string
	Process(topicDevId string, Data map[string]string) (bool, error)
}

func MatchDeviceType(devId string) (string, error) {
	if len(devId) == 0 {
		return "", errors.New("len(devId) == 0")
	}

	if len(devId) == 44 {
		return DEVICETYPE_ANDROID, nil
	} else if isUuid(devId) {
		return DEVICETYPE_IOS, nil
	} else if strings.Index(devId, "server_") > -1 {
		return DEVICETYPE_SERVER, nil
	} else if strings.HasPrefix(devId, "me") {
		return DEVICETYPE_MESH, nil
	} else if len(devId) >= 19 && len(devId) <= 24 { //24位~28位硬件
		return DEVICETYPE_GATEWAY, nil
	}
	return DEVICETYPE_OTHER, nil
}

func isUuid(devId string) bool {
	return MatchString(
		"^[0-9a-zA-F]{8}-[0-9a-zA-F]{4}-[0-9a-zA-F]{4}-[0-9a-zA-F]{4}-[0-9a-zA-F]{12}$",
		devId)
}

func MatchString(pattern string, s string) bool {
	ret, err := regexp.MatchString(pattern, s)
	if !ret && err != nil {
		return false
	}
	return ret
}
