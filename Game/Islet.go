package Game

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"math/rand"
	"reflect"
	"time"
	"websocket/conManage"
)

func IsletHandle(ws *websocket.Conn,userInfo *conManage.GUserInfo,recmsg string)  {
	//解析
	recData := new(conManage.GameData)
	err := json.Unmarshal([]byte(recmsg), recData)
	if err!=nil{
		fmt.Println("json failed:", err)
		websocket.JSON.Send(ws,&conManage.GameData{
			Code: "-1001",
			Type: "1",
			Data: err.Error(),
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	fmt.Println("type:",recData.Type)
	switch recData.Type {
		case "1":
			domove(ws,userInfo,recData.Data)
			break
		case "2":
			handlrandomAct(ws,userInfo,recData.Data)
			break
		default:
			websocket.JSON.Send(ws,&conManage.GameData{
				Code: "-1001",
				Type: "1",
				Data: "未定义操作",
				Time: time.Now().Format("2006-01-02 15:04:05"),
			})
			break
	}
}

func domove(ws *websocket.Conn,userInfo *conManage.GUserInfo,data interface{})  {
	fmt.Println(reflect.TypeOf(data))
	move,ok := data.(conManage.MoveAction)
	if !ok {
		websocket.JSON.Send(ws,&conManage.GameData{
			Code: "-1001",
			Type: "1",
			Data: "消息格式不正确",
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}
	//获取用户信息
	//Config.ConMaster.GetUserList()
	userInfo.Data.X +=move.X
	move.X = userInfo.Data.X
	userInfo.Data.Y +=move.Y
	move.Y = userInfo.Data.Y
	userInfo.Data.Z +=move.Z
	move.Z = userInfo.Data.Z
	websocket.JSON.Send(ws,&conManage.GameData{
		Code: "200",
		Type: "1",
		Data: move,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	})
}
func handlrandomAct(ws *websocket.Conn,userInfo *conManage.GUserInfo,data interface{})  {
	op,ok := data.(string)
	if !ok {
		websocket.JSON.Send(ws,&conManage.GameData{
			Code: "-1001",
			Type: "1",
			Data: "消息格式不正确",
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	if op == "start" {
		//开启定时随机
		if userInfo.RandChan == nil {
			userInfo.RandChan = make(chan bool)
			go func() {
				for {
					select {
					case <-userInfo.RandChan:
						fmt.Println("acch")
						//退出
						return
						break
					case <-time.After(5*time.Second):
						//随机事件
						fmt.Println("randact")
						randact(ws,userInfo)
						break
					}
				}
			}()
		}
	}else {
		//关闭定时随机
		if userInfo.RandChan != nil {
			userInfo.RandChan <- true
		}

	}
}

func randact(ws *websocket.Conn,userInfo *conManage.GUserInfo)  {
	rand.Seed(time.Now().UnixNano())
	randval := rand.Intn(100)   //生成0-99随机整数
	if randval <10 {
		//获得道具
		websocket.JSON.Send(ws,&conManage.GameData{
			Code: "200",
			Type: "2",
			Data: conManage.GItemData{
				Id: "1",
				Name: "随机物品",
				Type: "type",
				Attr: conManage.GAttrData{
					Level: rand.Intn(10),
					Atn:   rand.Intn(10),
					Int:   rand.Intn(10),
					Spd:   rand.Intn(10),
					Con:   rand.Intn(10),
					Ler:   rand.Intn(10),
					Def:   rand.Intn(10),
					Res:   rand.Intn(10),
				},
			},
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
	}else if randval < 20{
		//普通怪
		websocket.JSON.Send(ws,&conManage.GameData{
			Code: "200",
			Type: "2",
			Data: conManage.GAnimalData{
				Name : "test",
				Hp: rand.Intn(10),
				//力量
				Atn: rand.Intn(10),
				//防御
				Def : rand.Intn(10),
				Type : "type",
			},
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
	}else if randval < 25{
		//中等
		websocket.JSON.Send(ws,&conManage.GameData{
			Code: "200",
			Type: "2",
			Data: conManage.GAnimalData{
				Name : "test",
				Hp: rand.Intn(20),
				//力量
				Atn: rand.Intn(20),
				//防御
				Def : rand.Intn(20),
				Type : "type",
			},
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
	}else if randval < 27{
		//金鹰
		websocket.JSON.Send(ws,&conManage.GameData{
			Code: "200",
			Type: "2",
			Data: conManage.GAnimalData{
				Name : "test",
				Hp: rand.Intn(30),
				//力量
				Atn: rand.Intn(30),
				//防御
				Def : rand.Intn(30),
				Type : "type",
			},
			Time: time.Now().Format("2006-01-02 15:04:05"),
		})
	}
}