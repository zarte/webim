package Web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"websocket/Config"
	"websocket/conManage"
	"websocket/util"
)

func Userlist(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" {
		query := r.URL.Query()
		gettype :=query.Get("type")
		if gettype=="kefu" {
			//返回客服列表
			list := Config.ConMaster.GetKeFuList()
			result, err := json.Marshal(conManage.DataMessage{
				Code: "200",
				Msg:  "success",
				Data: list,
			})
			if err != nil {
				w.Write([]byte("api err"))
			} else {
				w.Write(result)
			}
		} else {
			token := query["token"]
			if len(token) <=0 {
				w.Write([]byte("no login"))
				return
			}
			ntoken, err := util.ParseToken(token[0])
			if err !=nil {
				w.Write([]byte("token err"))
				return
			}
			//判断uid是否已登录聊天室
			userInfo,gres :=Config.ConMaster.GetUserInfo(ntoken.UserId)
			if !gres {
				w.Write([]byte("no login2"))
				return
			}
			if !userInfo.(conManage.UserInfo).Iskefu {
				w.Write([]byte("no kefu"))
				return
			}
			//返回所有成员

			list := Config.ConMaster.GetUserList()
			result, err := json.Marshal(conManage.DataMessage{
				Code: "200",
				Msg:  "success",
				Data: list,
			})
			if err != nil {
				w.Write([]byte("api err"))
			} else {
				w.Write(result)
			}
		}
	}
}

func EmojiList(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" {
		query := r.URL.Query()
		path := query["path"]
		var list []conManage.EmojiInfo
		if len(path)>0 {
			switch path[0] {
			case "emoji":
				rd, err := ioutil.ReadDir(Config.Gconfig.CurExePath+"/imgs/emoji/")
				if err!=nil {
					fmt.Println(err)
				} else {
					for _, fi := range rd {
						if !fi.IsDir() {
							tmp :=strings.Split(fi.Name(),".")
							if len(tmp)==2 {
								list = append(list,conManage.EmojiInfo{
									Name: fi.Name(),
									Code: tmp[0],
									Path: "/imgs/emoji/"+fi.Name(),
								})
							}
						}
					}
				}
				break
			default:
			}
		}


		result, err := json.Marshal(conManage.DataMessage{
			Code: "200",
			Msg:  "success",
			Data: list,
		})
		if err != nil {
			w.Write([]byte("api err"))
		} else {
			w.Write(result)
		}
	}
}