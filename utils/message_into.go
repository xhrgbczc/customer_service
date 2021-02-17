package utils

import (
	"encoding/base64"
	"kf_server/models"
	"kf_server/services"
	"time"

	"github.com/astaxie/beego/orm"
)

// MessageInto push message
func MessageInto(message models.Message) {

	// 不处理的类型
	if message.BizType == "contacts" || message.BizType == "handshake" || message.BizType == "pong" || message.BizType == "welcome" || message.BizType == "into" || message.BizType == "search_knowledge" {
		return
	}

	// MessageRepository instance
	messageRepository := services.GetMessageRepositoryInstance()

	// message create time
	message.Timestamp = time.Now().Unix()
	message.Sequence = time.Now().UnixNano()

	// content内容转base64
	message.Payload = base64.StdEncoding.EncodeToString([]byte(message.Payload))

	// 接收者是客服
	admin := services.GetAdminRepositoryInstance().GetAdmin(message.ToAccount)

	// 接收者是客服
	if admin != nil {

		// 默认已读消息
		message.Read = 0
		if admin.Online == 0 || admin.CurrentConUser != message.FromAccount {
			message.Read = 1
		}

	} else {

		// UserRepository instance
		userRepository := services.GetUserRepositoryInstance()
		user := userRepository.GetUser(message.ToAccount)
		if user != nil {
			message.Read = 0
			if user.IsWindow == 0 {
				message.Read = 1
			}

			// 处理是否已回复
			admin := services.GetAdminRepositoryInstance().GetAdmin(message.FromAccount)
			payload, _ := base64.StdEncoding.DecodeString(message.Payload)
			if admin != nil && admin.AutoReply != string(payload) {
				services.GetStatisticalRepositoryInstance().CheckIsReplyAndSetReply(user.ID, message.FromAccount, user.Platform)
			}

		}

	}

	// message.BizType == "end" is not read
	if message.BizType == "end" || message.BizType == "timeout" {
		message.Read = 0
	}

	// 消息入库
	_, _ = messageRepository.Add(&message)

	// RobotRepository instance
	robotRepository := services.GetRobotRepositoryInstance()

	// 判断是否机器人对话（不处理聊天列表）
	if rbts, _ := robotRepository.GetRobotWithInIds(message.ToAccount, message.FromAccount); len(rbts) > 0 {
		return
	}

	// ContactRepository instance
	contactRepository := services.GetContactRepositoryInstance()

	// 处理客服聊天列表
	if contact, err := contactRepository.GetContactWithIds(message.ToAccount, message.FromAccount); err != nil {
		newContact := models.Contact{}
		newContact.ToAccount = message.ToAccount
		newContact.FromAccount = message.FromAccount
		newContact.LastMessageType = message.BizType
		newContact.CreateAt = time.Now().Unix()
		newContact.LastAccount = message.FromAccount
		newContact.LastMessage = message.Payload
		_, _ = contactRepository.Add(&newContact)
	} else {
		isSessionEnd := 0
		if message.BizType == "end" || message.BizType == "timeout" {
			isSessionEnd = 1
		}
		_, _ = contactRepository.Update(contact.ID, orm.Params{
			"LastMessageType": message.BizType,
			"CreateAt":        time.Now().Unix(),
			"LastMessage":     message.Payload,
			"LastAccount":     message.FromAccount,
			"IsSessionEnd":    isSessionEnd,
			"Delete":          0,
		})
	}

	// 接收者是客服就推送
	if admin != nil {
		go func() {
			time.Sleep(1 * time.Second)
			PushNewContacts(message.ToAccount)
		}()
	}

}
