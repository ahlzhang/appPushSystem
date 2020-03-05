package template

import "github.com/ahlzhang/appPushSystem/app/android/igetui/protobuf"
import proto "github.com/golang/protobuf/proto"

type TransmissionTemplate struct {
	AppId               string
	AppKey              string
	TransmissionType    int32
	TransmissionContent string
	PushType            string
}

func NewTransmissionTemplate(appid string, appkey string, transmissiontype int32,
	transmissionconntent string) *TransmissionTemplate {
	return &TransmissionTemplate{
		AppId:               appid,
		AppKey:              appkey,
		TransmissionType:    transmissiontype,
		TransmissionContent: transmissionconntent,
		PushType:            "TransmissionMsg",
	}
}

func (t *TransmissionTemplate) GetTransmissionContent() string {
	return t.TransmissionContent
}

func (t *TransmissionTemplate) GetPushType() string {
	return t.PushType
}

func (t *TransmissionTemplate) GetTransparent() *protobuf.Transparent {
	transparent := &protobuf.Transparent{
		Id:          proto.String(""),
		Action:      proto.String("pushmessage"),
		TaskId:      proto.String(""),
		AppKey:      proto.String(t.AppKey),
		AppId:       proto.String(t.AppId),
		MessageId:   proto.String(""),
		PushInfo:    t.GetPushInfo(),
		ActionChain: t.GetActionChains(),
	}
	return transparent
}

func (t *TransmissionTemplate) GetPushInfo() *protobuf.PushInfo {
	pushInfo := &protobuf.PushInfo{
		Message:   proto.String(""),
		ActionKey: proto.String(""),
		Sound:     proto.String(""),
		Badge:     proto.String(""),
	}
	return pushInfo
}

func (t *TransmissionTemplate) GetActionChains() []*protobuf.ActionChain {

	//set actionChain
	actionChain1 := &protobuf.ActionChain{
		ActionId: proto.Int32(1),
		Type:     protobuf.ActionChain_Goto.Enum(),
		Next:     proto.Int32(10030),
	}

	//appStartUp
	appStartUp := &protobuf.AppStartUp{
		Android: proto.String(""),
		Symbia:  proto.String(""),
		Ios:     proto.String(""),
	}

	//start up app
	actionChain2 := &protobuf.ActionChain{
		ActionId:     proto.Int32(10030),
		Type:         protobuf.ActionChain_startapp.Enum(),
		Appid:        proto.String(""),
		Autostart:    proto.Bool(t.TransmissionType == 1),
		Appstartupid: appStartUp,
		FailedAction: proto.Int32(100),
		Next:         proto.Int32(100),
	}

	//end
	actionChain3 := &protobuf.ActionChain{
		ActionId: proto.Int32(100),
		Type:     protobuf.ActionChain_eoa.Enum(),
	}

	actionChains := []*protobuf.ActionChain{actionChain1, actionChain2, actionChain3}
	return actionChains
}
