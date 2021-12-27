package Ws

import (
	"encoding/json"
	"fmt"
	idworker "github.com/gitstliu/go-id-worker"
	"golang.org/x/net/websocket"
	"strconv"
	"strings"
	"time"
	"websocket/Config"
	"websocket/Game"
	"websocket/conManage"
	"websocket/util"
)

func doIsletLogin(recmsg string,userinfo *conManage.GUserInfo,ip string) (bool,string)  {
	var err error
	logininfo := new(conManage.JwtLoginMsg)
	err = json.Unmarshal([]byte(recmsg), logininfo)
	if err!=nil{
		fmt.Println("json failed:", err)
		return false,"登录失败:10002"
	}
	//正则验证
	if logininfo.Jwt==""{
		return false,"验证信息缺少:10002"
	}
	token, err := util.ParseToken(logininfo.Jwt)
	if err !=nil {
		return false,err.Error()+":10002"
	}
	//判断uid是否已登录聊天室
	userInfo,gres :=Config.ConMaster.GetUserInfo(token.UserId)
	if !gres {
		// return false,"未登录:10002"
	}
	Config.Gconfig.GLoger.InfoLog("login:" + token.UserId+",ip:"+ip,"")
	fmt.Println(userInfo)
	//userinfo.Name = userInfo.(conManage.UserInfo).
	userinfo.Name = "test222"
	return true,token.UserId +"-"+userinfo.Name
}
func islethandshake(ws *websocket.Conn,curUserInfo *conManage.GUserInfo) bool {
	Config.Gconfig.GLoger.InfoLog("try connect:"+ws.Request().RemoteAddr,"")
	var err error
	var reply string
	initFlag := false
	loginch := make(chan bool)
	loginclosech :=  make(chan bool)
	go func() {
		for {
			_, ok := <-loginclosech
			if !ok {
				break
			}else{
				if err = websocket.Message.Receive(ws, &reply); err != nil {
					if err.Error() !="EOF"{
						Config.Gconfig.GLoger.InfoLog("receive failed:" + err.Error(),"")
					}
					loginch <- false
				}
				loginch <- true
			}
		}
		close(loginch)
	}()
	for {
		//开始等待登录
		loginclosech <- true
		select {
		case res := <-loginch:
			if !res {
				continue
			} else {
				//初次连接
				var initmsg conManage.Message
				res,errmsg := doIsletLogin(reply,curUserInfo,ws.Request().RemoteAddr)
				if !res {
					initmsg.Code = "1002"
					initmsg.Msg =errmsg
					msg, err := json.Marshal(initmsg)
					if err != nil {
						fmt.Println("json.marshal failed, err:", err)
						return false
					}
					if err = websocket.Message.Send(ws, string(msg)); err != nil {
						fmt.Println("send failed:", err)
						return false
					}
					//添加ip到ban列表
					addBanIp(ws.Request().RemoteAddr)
					continue
				}
				tmp := strings.Split(errmsg,"-")
				curUserInfo.Id = tmp[0]
				curUserInfo.Name = tmp[1]
				//用户信息完善
				curUserInfo.Data = conManage.UserData{
					Hp:    100,
					Mp:    100,
					Vp:    100,
					Level: 1,
					Atn:   1,
					Int:   1,
					Spd:   1,
					Con:   1,
					Ler:   1,
					Def:   1,
					Res:   1,
					X:     0,
					Y:     0,
					Z:     0,
				}
				Config.ConMaster.GConnected(tmp[0],ws,*curUserInfo)
				// 发送登录成功信息
				Config.ConMaster.GSend(curUserInfo.Id,&conManage.Message{
					Code: "1001",
					Msg:  "success",
					From: "",
					Data: curUserInfo,
					To:   "",
				})
				close(loginclosech)
				initFlag = true
				initmsg.Msg ="用户" +curUserInfo.Id+ "上线,当前用户数：" +strconv.Itoa(int(*Config.ConMaster.Online))
				initmsg.Code = "2001"
				initmsg.From = "system"
				Config.ConMaster.GBroadcast(ws,&initmsg)
			}
		case <-time.After(time.Second * time.Duration(Config.Gconfig.Sysini.LoginTimeOut)):
			if err = websocket.JSON.Send(ws, &conManage.Message{
				Code: "1003",
				Msg:  "未及时登录(1003)",
				From: "system",
				To:   "",
			}); err != nil {
				fmt.Println("send failed:", err)
				return false
			}
			close(loginclosech)
			return false
		}
		if initFlag {
			//已登录
			return true
		}
	}
	return false
}
/**
孤岛求生
 */
func Islet(ws *websocket.Conn) {
	//生成唯一id
	currWoker := &idworker.IdWorker{}
	currWoker.InitIdWorker(1, 1)
	newId , newIdErr:= currWoker.NextId()
	if newIdErr != nil {
		fmt.Println(newIdErr)
		Config.Gconfig.GLoger.ErrorLog(newIdErr.Error())
		return
	}
	var err error
	var reply string
	curUserInfo := conManage.GUserInfo{
		Id:     strconv.Itoa(int(newId)),
		Name:   "test",
		Avatar: "avatar",
		Ip: ws.Request().RemoteAddr,
	}

	if !islethandshake(ws,&curUserInfo){
		return
	}
	defer Config.ConMaster.GDisConnected(curUserInfo.Id)
	removeBanIp(ws.Request().RemoteAddr)

	for {
		//websocket接受信息
		ch := make(chan bool)
		go func() {
			if err = websocket.Message.Receive(ws, &reply); err != nil {
				if err.Error() !="EOF"{
					Config.Gconfig.GLoger.InfoLog("receive failed:" + err.Error(),"")
				}
				ch <- false
			}
			ch <- true
		}()

		select {
		case res := <-ch:
			if !res {
				return
			}
		case <-time.After(time.Second * time.Duration(Config.Gconfig.Sysini.ConnectTimeOut)):
			Config.ConMaster.Send(curUserInfo.Id,&conManage.Message{
				Code: "1003",
				Msg:  "你已断开连接(1003)",
				From: "system",
				To:   "",
			})
			return
		}
		fmt.Println("reveived from client: " + reply)
		if reply=="11" {
			//心跳
			continue
		}
		Game.IsletHandle(ws,&curUserInfo,reply)
	}
}