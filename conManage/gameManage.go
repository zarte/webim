package conManage

import (
	"fmt"
	"golang.org/x/net/websocket"
	"sync/atomic"
)
type GameData struct {
	Code        string       `json:"code"`
	Type        string       `json:"type"`
	Data        interface{}       `json:"data"`
	Time        string       `json:"time"`
}
type GItemData struct {
	Id        string       `json:"id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Attr 	GAttrData       `json:"attr"`
}
type GAttrData struct {
	Level        int       `json:"level"`
	//力量
	Atn        int       `json:"atn"`
	//法共
	Int       int	`json:"int"`
	//敏捷
	Spd        int       `json:"spd"`
	//体制
	Con        int       `json:"con"`
	//智力
	Ler       int	`json:"ler"`
	//防御
	Def       int	`json:"def"`
	//魔防御
	Res       int	`json:"res"`
}
type GAnimalData struct {
	Name        string       `json:"name"`
	Hp        int       `json:"hp"`
	//力量
	Atn        int       `json:"atn"`
	//防御
	Def       int	`json:"def"`
	Type        string       `json:"type"`
}
type MoveAction struct {
	X        int       `json:"x"`
	Y        int       `json:"y"`
	Z        int      `json:"z"`
}
type GUserInfo struct {
	Id        string      `json:"id"`
	Name        string       `json:"name"`
	Avatar        string       `json:"avatar"`
	Ip        string      `json:"-"`
	RandChan     chan bool    `json:"-"`
	Data        UserData      `json:"data"`
}
type UserData struct {
	Hp        int       `json:"hp"`
	Mp        int       `json:"mp"`
	Vp        int       `json:"vp"`
	Level        int       `json:"level"`
	//力量
	Atn        int       `json:"atn"`
	//法共
	Int       int	`json:"int"`
	//敏捷
	Spd        int       `json:"spd"`
	//体制
	Con        int       `json:"con"`
	//智力
	Ler       int	`json:"ler"`
	//防御
	Def       int	`json:"def"`
	//魔防御
	Res       int	`json:"res"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
	Z        int      `json:"z"`
}


func (m *ConnManager) GConnected(k, v interface{},u interface{}) {
	m.gconnections.Store(k, v)
	m.guserlist.Store(k, u)
	atomic.AddInt32(m.GOnline, 1)
}

func (m *ConnManager) GBroadcast(conn *websocket.Conn, msg *Message) {
	m.Foreach(func(k, v interface{}) {
		if c, ok := v.(*websocket.Conn); ok && c != conn {
			if err := websocket.JSON.Send(c, msg); err != nil {
				fmt.Println("Send msg error: ", err)
			}
		}
	})
}
func (m *ConnManager) GDisConnected(k interface{}) {
	fmt.Println("disconnect" )
	m.gconnections.Delete(k)
	m.guserlist.Delete(k)
	atomic.AddInt32(m.GOnline, -1)
	gsendLogout(m,k)
}

func (m *ConnManager) GetGUserList() []GUserInfo{
	var list []GUserInfo
	m.guserlist.Range(func(k, v interface{}) bool {
		list = append(list,v.(GUserInfo) )
		return true
	})
	return list
}
func gsendLogout(m *ConnManager,k interface{})  {
	user :=k.(string)
	msg :="用户:"+user+"下线"
	m.GBroadcast(nil, &Message{
		Code:"2002",
		Msg: msg,
		From:"system",
		To:user,
	})
}
func (m *ConnManager) GGet(k interface{}) (v interface{}, ok bool) {
	return m.gconnections.Load(k)
}

func (m *ConnManager) GSend(k string, msg *Message) (bool,string){
	v, ok := m.GGet(k)
	if ok {
		if conn, ok := v.(*websocket.Conn); ok {
			if err := websocket.JSON.Send(conn, msg); err != nil {
				fmt.Println("Send msg error: ", err)
				return false,""
			}
			return true,""
		} else {
			fmt.Println("invalid type, expect *websocket.Conn")
			return false,""
		}
	} else {
		fmt.Println("connection not exist2",msg)
		return false,"connection not exist"
	}
}