package igetui

type Target struct {
	AppId    string
	ClientId string
}

func NewTarget(appid, clientid string) *Target {
	return &Target{
		AppId:    appid,
		ClientId: clientid,
	}
}
