package Web

import (
	"html/template"
	"net/http"
	"websocket/Config"
)
func Web(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端
		t, _ := template.ParseFiles(Config.Gconfig.CurExePath +"client.html")
		t.Execute(w, nil)
	}
}
func Kefu(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端
		t, _ := template.ParseFiles(Config.Gconfig.CurExePath +"kefu.html")
		t.Execute(w, nil)
	}
}
