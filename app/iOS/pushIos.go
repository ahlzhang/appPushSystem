/**
 * iOS推送.
 *
 * User: zhangbob
 * Date: 2018/5/9
 * Time: 上午10:16
 */
package iOS

import (
	"appPushSystem"
	"errors"

	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

func NewIosPush(pem []byte, password, packageName string, isDev bool) (error, *messageHandle) {
	cert, err := certificate.FromPemBytes(pem, password)
	if err != nil {
		return err, nil
	}

	if packageName == "" {
		return errors.New("包名不能为空"), nil
	}

	m := &messageHandle{
		packageName: packageName,
	}

	if !isDev {
		m.apnsClient = apns.NewClient(cert).Production()
	} else {
		m.apnsClient = apns.NewClient(cert).Development()
	}

	return nil, m
}

type messageHandle struct {
	apnsClient  *apns.Client
	packageName string
}

func (t messageHandle) PushSingle(message appPushSystem.IMessage, callBack appPushSystem.IHandleMessageCallback) {
	_, clientId, content := message.ToMessage()
	notification := &apns.Notification{}
	notification.DeviceToken = clientId
	notification.Payload = []byte(content)
	notification.Topic = t.packageName

	res, err := t.apnsClient.Push(notification)
	if err != nil {
		callBack.Fail(message, err)
		return
	}

	// 200  Success
	// 400 Bad request
	// 403 There was an error with the certificate.
	// 405 The request used a bad :method value. Only POST requests are supported.
	// 410 The device token is no longer active for the topic.
	// 413 The notification payload was too large.
	// 429 The server received too many requests for the same device token.
	// 500 Internal server error
	// 503 The server is shutting down and unavailable.
	if res.Sent() {
		callBack.Success(message)
	} else {
		callBack.Fail(message, errors.New(res.Reason))
	}
}
