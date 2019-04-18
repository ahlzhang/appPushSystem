package template

import "jiaotou.com/appPushSystem/control/igetui/protobuf"

//import proto "code.google.com/p/goprotobuf/proto"

//type BaseTemplate struct {
//	AppKey              string
//	AppId               string
//	TransmissionType    int32
//	TransmissionContent string
//	PushType            string
//}

type ITemplate interface {
	GetTransparent() *protobuf.Transparent
	GetActionChains() []*protobuf.ActionChain
	GetPushInfo() *protobuf.PushInfo
	GetTransmissionContent() string
	GetPushType() string
}

//func (bt *BaseTemplate) GetTransparent() *protobuf.Transparent {
//	transparent := &protobuf.Transparent{
//		Id:          proto.String(""),
//		Action:      proto.String("pushMessage"),
//		TaskId:      proto.String(""),
//		AppKey:      proto.String(bt.AppKey),
//		AppId:       proto.String(bt.AppId),
//		MessageId:   proto.String(""),
//		PushInfo:    bt.GetPushInfo(),
//		ActionChain: bt.GetActionChains(),
//	}
//	return transparent

//}

//func (bt *BaseTemplate) GetActionChains() []*protobuf.ActionChain {
//	return []*protobuf.ActionChain{}

//}

//func (bt *BaseTemplate) GetPushInfo() *protobuf.PushInfo {
//	pushInfo := &protobuf.PushInfo{
//		Message:   proto.String(""),
//		ActionKey: proto.String(""),
//		Sound:     proto.String(""),
//		Badge:     proto.String(""),
//	}
//	return pushInfo
//}
