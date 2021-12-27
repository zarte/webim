package conManage

import (
	"fmt"
	"github.com/zarte/comutil/Zloger"
	"golang.org/x/net/websocket"
	"sync"
	"sync/atomic"
)


type Message struct {
	Code        string      `json:"code"`
	Msg        string       `json:"msg"`
	From        string       `json:"from"`
	To        string       `json:"to"`
	Data       interface{}   `json:"data"`
	Timestr        string       `json:"time"`
}
type DataMessage struct {
	Code        string      `json:"code"`
	Msg        string       `json:"msg"`
	Data       interface{}   `json:"data"`
}
type LoginMsg struct {
	User        string      `json:"user"`
	Passwd        string       `json:"passwd"`
}
type JwtLoginMsg struct {
	Jwt        string      `json:"jwt"`
}
type UserInfo struct {
	Id        string      `json:"id"`
	Name        string       `json:"name"`
	Iskefu        bool       `json:"iskefu"`
	Avatar        string       `json:"avatar"`
	Ip        string      `json:"-"`
}


type EmojiInfo struct {
	Name        string       `json:"name"`
	Path        string       `json:"path"`
	Code        string       `json:"code"`
}

type ConnManager struct {
	// websocket connection number
	Online *int32
	GOnline *int32
	// websocket connection
	connections *sync.Map
	userlist *sync.Map
	gconnections *sync.Map
	guserlist *sync.Map
}

var GLoger *Zloger.Loger
/**
初始化
 */
func NewManager() *ConnManager {
	//返回日志对象
	manage := &ConnManager{}
	manage.Online = new (int32)
	manage.connections = new (sync.Map)
	manage.userlist = new (sync.Map)

	manage.GOnline = new (int32)
	manage.gconnections = new (sync.Map)
	manage.guserlist = new (sync.Map)
	return manage
}

func (m *ConnManager) Connected(k, v interface{},u interface{}) {
	m.connections.Store(k, v)
	m.userlist.Store(k, u)
	atomic.AddInt32(m.Online, 1)
}
// remove websocket connection by key
// online number - 1
func (m *ConnManager) DisConnected(k interface{}) {
	fmt.Println("disconnect:"+k.(string) )
	m.connections.Delete(k)
	m.userlist.Delete(k)
	atomic.AddInt32(m.Online, -1)
	sendLogout(m,k)
}

func sendLogout(m *ConnManager,k interface{})  {
	user :=k.(string)
	msg :="用户:"+user+"下线"
	m.Broadcast(nil, &Message{
		Code:"2002",
		Msg: msg,
		From:"system",
		To:user,
	})
}
// get websocket connection by key
func (m *ConnManager) Get(k interface{}) (v interface{}, ok bool) {
	return m.connections.Load(k)
}

// iter websocket connections
func (m *ConnManager) Foreach(f func(k, v interface{})) {
	m.connections.Range(func(k, v interface{}) bool {
		f(k, v)
		return true
	})
}

func (m *ConnManager) GetUserList() []UserInfo{
	var list []UserInfo
	m.userlist.Range(func(k, v interface{}) bool {
		list = append(list,v.(UserInfo) )
		return true
	})
	return list
}
func (m *ConnManager) GetKeFuList() []UserInfo{
	var list []UserInfo
	m.userlist.Range(func(k, v interface{}) bool {
		if v.(UserInfo).Iskefu {
			list = append(list,v.(UserInfo) )
		}
		return true
	})
	return list
}
func (m *ConnManager) GetUserInfo(k interface{}) (v interface{}, ok bool) {
	return m.userlist.Load(k)
}

func (m *ConnManager) CheckNameUnique(mname string) bool {
	var match bool
	match = false
	m.userlist.Range(func(k, v interface{}) bool {
		if v.(UserInfo).Name == mname{
			match = true
		}
		return true
	})
	return match
}

// send message to one websocket connection
func (m *ConnManager) Send(k string, msg *Message) (bool,string){
	v, ok := m.Get(k)
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

// send message to multi websocket connections
func (m *ConnManager) SendMulti(keys []string, msg interface{}) {
	for _, k := range keys {
		v, ok := m.Get(k)
		if ok {
			if conn, ok := v.(*websocket.Conn); ok {
				if err := websocket.JSON.Send(conn, msg); err != nil {
					fmt.Println("Send msg error: ", err)
				}
			} else {
				fmt.Println("invalid type, expect *websocket.Conn")
			}
		} else {
			fmt.Println("connection not exist1")
		}
	}
}

// broadcast message to all websocket connections otherwise own connection
func (m *ConnManager) Broadcast(conn *websocket.Conn, msg *Message) {
	m.Foreach(func(k, v interface{}) {
		if c, ok := v.(*websocket.Conn); ok && c != conn {
			if err := websocket.JSON.Send(c, msg); err != nil {
				fmt.Println("Send msg error: ", err)
			}
		}
	})
}
