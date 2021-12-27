package Ws

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"html"
	"regexp"
	"time"
	"websocket/Config"
	"websocket/conManage"
)

/**
处理消息发送
*/
func transferMsg(ws *websocket.Conn, reply string,curUserInfo conManage.UserInfo){
	var sendmsg conManage.Message
	recmsg := new(conManage.Message)
	err := json.Unmarshal([]byte(reply), recmsg)
	if err!=nil{
		fmt.Println("json failed:", err)
		return
	}
	//消息处理
	sendmsg.Code = "200"
	sendmsg.From = curUserInfo.Name
	sendmsg.Msg = msgAclDeal(recmsg.Msg)
	sendmsg.Timestr  = time.Now().Format("2006-01-02 15:04:05")
	if recmsg.To == "" {
		if !curUserInfo.Iskefu{
			//非客服直接丢弃消息
			return
		}
		//广播
		Config.ConMaster.Broadcast(ws,&sendmsg)
	}else {
		//私聊
		sendmsg.To = recmsg.To
		//非客服不可发送消息给任意用户
		if !curUserInfo.Iskefu {
			// 判断消息接收对象是否为客服
			toinfo,flag := Config.ConMaster.GetUserInfo(recmsg.To)
			if !flag{
				Config.Gconfig.GLoger.InfoLog("send user getfail:" + sendmsg.From +"-"+ws.Request().RemoteAddr,"acl_")
				return
			}
			if !toinfo.(conManage.UserInfo).Iskefu {
				Config.Gconfig.GLoger.InfoLog("send to no kefu:" + sendmsg.From +"-"+ws.Request().RemoteAddr,"acl_")
				return
			}
		}
		//这里是发送消息
		res,errmsg :=Config.ConMaster.Send(recmsg.To,&sendmsg)
		if !res {
			sendmsg.Code = "4"
			sendmsg.Msg = errmsg
			if err := websocket.JSON.Send(ws, sendmsg); err != nil {
				fmt.Println("Send msg error: ", err)
			}
		}
	}
}
func msgAclDeal(msg string) string  {
	msg = html.EscapeString(msg)
	re2 := regexp.MustCompile(`\[\:([a-zA-Z0-9\-]*?)\:\]`)
	msg = re2.ReplaceAllString(msg, `<img src="/imgs/emoji/${1}.png"/>`)
	return msg
}
