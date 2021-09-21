package logic

import (
	"fmt"
)

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, MessageQueueLen),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

type broadcaster struct {
	// 所有聊天室用户
	users map[string]*User

	// 所有 channel 统一管理，可以避免外部乱用

	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判断该昵称用户是否可进入聊天室（重复与否）：true 能，false 不能
	checkUserChannel      chan string
	checkUserCanInChannel chan bool
}

func (b *broadcaster) Start() {
	for  {
		select {
			case user := <- b.enteringChannel:
				// 新用户进入,添加到用户列表，然后广播欢迎消息
				b.users[user.NickName] = user
				b.Broadcast(NewWelcomeMessage(user))
				OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
				// 用户离开,删除该用户，关闭该用户的channel，然后广播离开消息
				delete(b.users, user.NickName)
				// 避免 goroutine 泄露
				user.CloseMessageChannel()
				b.Broadcast(NewUserLeaveMessage(user))
			case msg := <- b.messageChannel:
				// 给所有在线用户发送消息
				for _, user := range b.users {
					if user.UID == msg.User.UID {
						continue
					}
					fmt.Printf("发送消息给%s,内容是[%s]\n",user.NickName,msg.Content)
					user.MessageChannel <- msg
				}
				OfflineProcessor.Save(msg)
			case nickname := <- b.checkUserChannel:
				if _, ok := b.users[nickname]; ok {
					b.checkUserCanInChannel <- false
				} else {
					b.checkUserCanInChannel <- true
				}
		}
	}
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) UserEntering(user *User)  {
	b.enteringChannel <- user
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(message *Message) {
	b.messageChannel <- message
}