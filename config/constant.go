/** 
 * 常量.
 *
 * User: zhangbob 
 * Date: 2018/5/22 
 * Time: 下午7:38 
 */
package config

const (
	_             = iota
	SystemAndroid  //操作系统 android
	SystemIos      //操作系统 ios
)

const (
	ContentTypeNotJump     = iota //不跳转
	ContentTypeJumpOutside        //跳转外部链接
	ContentTypeJumpApp            //跳转内部app
)

const (
	_               = iota
	PushRangeAll     // 推送所有人
	PushRangeCustom  // 推送指定人群
)

const (
	_                 = iota
	MessageTypeSystem  //系统消息
	MessageTypeHot     //热门消息
)

const (
	ReadStateNotRead     = iota //消息状态 未读
	ReadStateAlreadyRead        //消息状态 已读

)

const (
	PushStateWaiting = iota //消息状态 等待推送
	PushStateAlready        //消息状态 已推送
	PushStateFail           //消息状态 推送失败
	PushStateSending        //消息状态 推送中
	PushCanelSend           //消息状态 取消推送
)

const (
	DelStateNot     = iota //删除状态 未删除
	DelStateAlready        //删除状态 已删除
)
