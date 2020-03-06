/** 
 * @api post core.
 *
 * User: zhangbob 
 * Date: 2018-11-01 
 * Time: 18:29 
 */
package core

type IMessageHandle interface {
	PushSingle(message IMessage, callback IHandleMessageCallback)
}

type IMessage interface {
	/**
	 * 转成单条消息
	 *
	 * @return pushId int64 推送的id
	 * @return clientId string 推送手机的唯一标实
	 * @return content string 推送的内容
	 **/
	ToMessage() (pushId int64, clientId string, content string)

	/**
	 * 获取手机类型
	 *
	 * 1:android 2:iOS
	 **/
	GetSystemType() int
}

//消息处理后的回调
type IHandleMessageCallback interface {
	//消息处理中
	Sending(message IMessage)
	//消息处理成功
	Success(message IMessage)
	//消息处理失败
	Fail(message IMessage, err error)
}
