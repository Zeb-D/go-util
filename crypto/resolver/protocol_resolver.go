package resolver

import "errors"

var ProtocolResolvers map[string]IProtocolResolver = make(map[string]IProtocolResolver)

//	InitResolver 实现类自己主动注册进来
func InitResolver(r IProtocolResolver) error {
	if r == nil {
		return errors.New("IProtocolResolver is nil")
	}
	ProtocolResolvers[r.ProtocolVersion()] = r
	return nil
}

func GetDeviceProtocolResolver(protocolVersion string) IProtocolResolver {
	pr := ProtocolResolvers[protocolVersion]
	if pr == nil {
		pr = ProtocolResolvers[V_1_0]
	}
	return pr
}

type IProtocolResolver interface {
	ProtocolVersion() string
	Decode(protocolData interface{}, devId, localKey, protocolVersion string) (string, error)
}
