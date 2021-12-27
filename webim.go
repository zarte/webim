package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"websocket/Config"
	"websocket/Web"
	"websocket/Ws"
)

func initserver()  {
	Config.SetConfig()
}
func main() {
	initserver()
	//接受websocket的路由地址
	http.Handle("/websocket", websocket.Handler(Ws.Echo))
	http.Handle("/islet", websocket.Handler(Ws.Islet))

	//html页面
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", Web.Web)
	http.HandleFunc("/kefu", Web.Kefu)
	http.HandleFunc("/userlist", Web.Userlist)
	http.HandleFunc("/emojilist", Web.EmojiList)
	http.Handle("/imgs/",   http.StripPrefix("/imgs/", http.FileServer(http.Dir("imgs"))))

	fmt.Println("start listen:"+Config.Gconfig.WebPort)
	if err := http.ListenAndServe(":"+Config.Gconfig.WebPort, nil); err != nil {
		fmt.Println(err)
	}

}