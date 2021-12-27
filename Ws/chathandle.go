package Ws

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	idworker "github.com/gitstliu/go-id-worker"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
	"websocket/Config"
	"websocket/conManage"
	"websocket/util"
)

func doLogin(recmsg string,userinfo *conManage.UserInfo,ip string) (bool,string)  {
	var err error
	logininfo := new(conManage.LoginMsg)
	err = json.Unmarshal([]byte(recmsg), logininfo)
	if err!=nil{
		fmt.Println("json failed:", err)
		return false,"登录失败:10002"
	}
	//正则验证
	if logininfo.User==""{
		return false,"昵称不能为空:10002"
	}
	match,err:=regexp.MatchString("^[0-9a-zA-Z\u4e00-\u9fa5]{6,20}$",logininfo.User)
	if !match {
		return false,"昵称格式不符要求:10002"
	}
	if Config.ConMaster.CheckNameUnique(logininfo.User) {
		return false,"已存在相同用户:10002"
	}
	Config.Gconfig.GLoger.InfoLog("login:" + logininfo.User+",ip:"+ip,"")
	if logininfo.Passwd!=""{
		//客服
		fmt.Println(Config.Gconfig.KeFuAcc)
		for i,k :=range Config.Gconfig.KeFuAcc {
			fmt.Println(i)
			if i == logininfo.User {
				if logininfo.Passwd!="" && k==logininfo.Passwd {
					userinfo.Name = logininfo.User
					userinfo.Iskefu = true
					return true,""
				}
			}
		}
		return false,"登录失败:10003"
	}else {
		//客户
		//判断名称与客服名称是否重复
		for i,_ :=range Config.Gconfig.KeFuAcc {
			if i ==logininfo.User {
				return false,"非法操作:10002"
			}
		}
	}
	userinfo.Name = logininfo.User
	return true,""
}


/**
生成token
 */
func getJwt(userid string, userInfo *conManage.UserInfo) string  {
	customClaims :=&util.CustomClaims{
		UserId: userid,//用户id
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(Config.Gconfig.JwtExpire)*time.Second).Unix(), // 过期时间，必须设置
			Issuer:userInfo.Name,   // 非必须，也可以填充用户名，
		},
	}
	//采用HMAC SHA256加密算法
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString,err:= token.SignedString([]byte(Config.Gconfig.SecretKey))
	if err!=nil {
		fmt.Println("get token err",err)
		return ""
	}
	return tokenString
}

/**
握手初始化
 */
func handshake(ws *websocket.Conn,curUserInfo *conManage.UserInfo) bool {
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
				if Config.Gconfig.Sysini.LoginV {
					res,errmsg := doLogin(reply,curUserInfo,ws.Request().RemoteAddr)
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
				}
				//生成jwt
				token := getJwt(curUserInfo.Id, curUserInfo)
				if token =="" {
					return false
				}
				Config.ConMaster.Connected(curUserInfo.Id,ws,*curUserInfo)
				// 发送登录成功信息

				Config.ConMaster.Send(curUserInfo.Id,&conManage.Message{
					Code: "1001",
					Msg:  token,
					From: "",
					To:   "",
				})
				close(loginclosech)
				initFlag = true
				initmsg.Msg ="用户" +curUserInfo.Id+ "上线,当前用户数：" +strconv.Itoa(int(*Config.ConMaster.Online))
				initmsg.Code = "2001"
				initmsg.From = "system"
				Config.ConMaster.Broadcast(ws,&initmsg)
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
func Echo(ws *websocket.Conn) {
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
	curUserInfo := conManage.UserInfo{
		Id:     strconv.Itoa(int(newId)),
		Name:   "test",
		Avatar: "avatar",
		Iskefu: false,
		Ip: ws.Request().RemoteAddr,
	}

	if !handshake(ws,&curUserInfo){
		return
	}
	go func() {
		updateMaxOnline()
	}()
	defer Config.ConMaster.DisConnected(curUserInfo.Id)
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
		if reply=="11" {
			//心跳
			continue
		}
		if reply=="kefulist" {
			//客服列表
			dealKefuList(curUserInfo.Id)
			continue
		}
		Config.Gconfig.GLoger.InfoLog("reveived from client: " + reply,"")
		fmt.Println("reveived from client: " + reply)
		transferMsg(ws,reply,curUserInfo)
	}
}
func dealKefuList(userid string)  {
	list := Config.ConMaster.GetKeFuList()
	Config.ConMaster.Send(userid,&conManage.Message{
		Code: "5001",
		Msg:  "success",
		Data: list,
	})
}
func updateMaxOnline()  {
	var f    *os.File
	var err   error

	var filename = "maxonline.txt"

	var content []byte
	if checkFileIsExist(Config.Gconfig.CurExePath+filename) {  //如果文件存在
		content, err = ioutil.ReadFile(Config.Gconfig.CurExePath+filename)
		if err != nil {
			fmt.Println(err)
			content = nil
		}
		f, err = os.OpenFile(Config.Gconfig.CurExePath+filename,  os.O_RDWR, os.ModePerm)  //打开文件
		//fmt.Println("文件存在");
	}else {
		f, err = os.Create(Config.Gconfig.CurExePath+filename)  //创建文件
		//fmt.Println("文件不存在");
	}
	check(err)
	defer f.Close()
	if string(content) !="" {
		onum, err := strconv.Atoi(string(content))
		if err !=nil {
			fmt.Println(err)
			return
		}
		if onum>=int(*Config.ConMaster.Online){
			return
		}
	}

	_,err = io.WriteString(f, strconv.Itoa(int(*Config.ConMaster.Online))) //写入文件(字符串)
	check(err)
	return
}
func addBanIp(ip string)  {
	if _, ok := Config.Gconfig.BanIp[ip]; !ok {
		Config.Gconfig.BanIp[ip] ++
	}else {
		Config.Gconfig.BanIp[ip] = 1
	}
}
func removeBanIp(ip string)  {
	delete(Config.Gconfig.BanIp, ip)
}
func clearBanIp()  {
	Config.Gconfig.BanIp = nil
}
func checkFileIsExist(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fmt.Println(err)
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}