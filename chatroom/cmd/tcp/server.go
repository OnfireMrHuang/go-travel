package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"time"
	"fmt"
)

var enteringChannel chan *User
var leavingChannel chan *User
var messageChannel chan string

func main()  {
	listener,err := net.Listen("tcp",":2020")
	if err != nil {
		panic(err)
	}
	go broadcaster()

	for {
		conn,err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster()  {
	users := make(map[*User]struct{})

	for {
		select {
			case user := <-enteringChannel:
				//新用户进入
				users[user] = struct{}{}
			case user:= <-leavingChannel:
				// 用户离开
				delete(users,user)
				// 关掉用户的channel，避免goroutine泄漏
				close(user.MessageChannel)
			case msg <- messageChannel:
				for user := range users {

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

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn)  {
	defer conn.Close()

	// 新用户进来，构建该用户的示例
	user := &User{
		ID: GenUserID(),
		Addr: conn.RemoteAddr().String(),
		EnterAt: time.Now(),
		MessageChannel: make(chan string,8),
	}
	// 当前在一个新的goroutine中，用来进行读操作，因此需要开一个goroutine用于写操作
	go sendMessage(conn,user.MessageChannel)

	// 给当前用户发送欢迎信息，给所有用户告知新用户到来
	user.MessageChannel <- "Welcome, " + user.String()
	messageChannel <- "user: " + strconv.Itoa(user.ID) + "` has enter"

	// 将该记录到全局的用户列表中，避免用锁
	enteringChannel <- user

	// 循环读取用户的输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- strconv.Itoa(user.ID) + ":" + input.Text()
	}
	if err := input.Err(); err !=  nil {
		log.Println("读取错误: ",err)
	}

	// 用户离开
	leavingChannel <- user
	messageChannel <- "user: `" + strconv.Itoa(user.ID) + "` has left"
}
