package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	// 新用户到来，通过该 channel 进行登记
	enteringChannel = make(chan *User)
	// 用户离开，通过该 channel 进行登记
	leavingChannel = make(chan *User)
	// 广播专用的用户普通消息 channel，缓冲是尽可能避免出现异常情况堵塞，这里简单给了 8，具体值根据情况调整
	messageChannel = make(chan Message, 8)
)

type Message struct {
	OwnerID int
	Content string
}

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}
	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			//新用户进入
			users[user] = struct{}{}
		case user := <-leavingChannel:
			// 用户离开
			delete(users, user)
			// 关掉用户的channel，避免goroutine泄漏
			close(user.MessageChannel)
		case msg := <-messageChannel:
			// 给所有在线用户发送消息
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}
}

type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u *User) String() string {
	return "用户:" + u.Addr
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// 新用户进来，构建该用户的示例
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}
	// 当前在一个新的goroutine中，用来进行读操作，因此需要开一个goroutine用于写操作
	go sendMessage(conn, user.MessageChannel)

	// 给当前用户发送欢迎信息，给所有用户告知新用户到来
	user.MessageChannel <- "Welcome, " + user.String()
	messageChannel <- Message{
		Content: "user: " + strconv.Itoa(user.ID) + "` has enter",
	}

	// 将该记录到全局的用户列表中，避免用锁
	enteringChannel <- user

	// 控制超时将用户踢出
	var userActive = make(chan struct{})
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 循环读取用户的输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- Message{
			Content: strconv.Itoa(user.ID) + ":" + input.Text(),
		}

		// 用户活跃
		userActive <- struct{}{}
	}
	if err := input.Err(); err != nil {
		log.Println("读取错误: ", err)
	}

	// 用户离开
	leavingChannel <- user
	messageChannel <- Message{
		Content: "user: `" + strconv.Itoa(user.ID) + "` has left",
	}
}

//func GenUserID() string {
//	uid,err := uuid.NewV1()
//	if err != nil {
//		log.Fatal(err)
//	}
//	return uid.String()
//}

func GenUserID() int {
	return 1
}
