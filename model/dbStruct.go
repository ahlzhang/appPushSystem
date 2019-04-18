/** 
 * model.
 *
 * User: zhangbob 
 * Date: 2018/5/9 
 * Time: 上午11:37 
 */
package model

type PushStruct struct {
	Id         int64  `xorm:"id"`
	UserId     int64  `xorm:"user_id"`
	MessageId  int64  `xorm:"message_id"`
	PushState  int    `xorm:"push_state"`
	Title      string `xorm:"title"`
	Content    string `xorm:"content"`
	CreateTime int64  `xorm:"create_time"`
	AppType    int    `xorm:"app_type"`
	PushTo     string `xorm:"cid"`
	TypeId     int    `xorm:"type_id"`
	PushRange  int    `xorm:"push_range"` //推送范围 1群体 2单体 3全部用户
}

type MessagesCenter struct {
	MessageId   int64  `xorm:"'message_id' autoincr"`
	TypeId      int    `xorm:"type_id"` //消息类型
	Title       string `xorm:"title"`
	Img         string `xorm:"img"`
	ContentType int    `xorm:"content_type"`
	Status      int    `xorm:"status"` //推送状态（全部用户推送才会有此字段）
	Url         string `xorm:"url"`
	CreateTime  int64  `xorm:"create_time"`
	Content     string `xorm:"content"`
	PushRange   int    `xorm:"push_range"` //推送范围 1群体 2单体 3全部用户
	Channel     int    `xorm:"channel"`    //消息发送渠道 1app 2运营后台
}

func (MessagesCenter) TableName() string {
	return "t_app_message_center"
}

type MessageUser struct {
	Id         int64 `xorm:"id"`
	UserId     int64 `xorm:"user_id"`
	MessageId  int64 `xorm:"message_id"`
	ReadState  int   `xorm:"read_state"` //用户是否已读 0未读取 1读取
	ReadTime   int64 `xorm:"read_time"`
	CreateTime int64 `xorm:"create_time"`
	PushState  int   `xorm:"push_state"` //推送状态 0未推送 1已推送 2推送失败 3推送中
	PushTime   int64 `xorm:"push_time"`
	IsDel      int   `xorm:"is_del"` //是否删除 0未删除 1已删除
	DelTime    int64 `xorm:"del_time"`
}

func (MessageUser) TableName() string {
	return "t_app_message_user_state"
}

type UserInfo struct {
	UserId int64  `xorm:"user_id"`
	Cid    string `xorm:"cid"`                   //推送ID
	Type   int    `xorm:"operating_system_type"` //设备类型 1.安卓 2.iOS
}

func (UserInfo) TableName() string {
	return "t_app_user"
}

type MessageType struct {
	TypeId     int    `json:"typeId"`
	TypeName   string `json:"typeName"`
	OrderId    int    `json:"orderId"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
	TypeImg    string `json:"typeImg"`
}

func (MessageType) TableName() string {
	return "t_app_message_type"
}
