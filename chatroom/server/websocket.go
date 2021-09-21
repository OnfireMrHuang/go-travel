package server

import (
	"go-travel/chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request)  {

	conn,err := websocket.Accept(w,req,nil)
	if err != nil {
		log.Println("websocket accept error: ",err)
		return
	}

	// 1： 新用户进来，构建该用户的实例
	nickname := req.FormValue("nickname")
	if l := len(nickname); l < 2 || 1 > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度: 4-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已经存在: ", nickname)
		wsjson.Write(req.Context(),conn, logic.NewErrorMessage("该昵称已经存在!"))
		return
	}


	user := logic.NewUser(conn, nickname, req.RemoteAddr)

	// 2、开启给用户回执消息的goroutine
	go user.SendMessage(req.Context())


	// 3、通知其他人该用户的加入
	logic.Broadcaster.UserEntering(user)
	log.Println("user:" , nickname, "joins chat")

	// 5、接收用户消息,并广播到聊天室
	err = user.ReceiveMessage(req.Context())

	// 6、用户离开
	logic.Broadcaster.UserLeaving(user)
	log.Println("user:", nickname, "leaves chat")

	// 根据读取时的错误执行不同的 Close链接
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
