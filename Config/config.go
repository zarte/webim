package Config

import (
	"fmt"
	"github.com/zarte/comutil/Zloger"
	"github.com/zarte/comutil/Goconfig"
	"os"
	"path/filepath"
	"strings"
	"websocket/conManage"
)

var Gconfig = new(Sysconfig)

type SysIni struct {
	LoginV bool
	ConnectTimeOut int
	LoginTimeOut int
}
var ConMaster *conManage.ConnManager
type Sysconfig struct {
	 Sysini SysIni
	 Tkey int
	 CurExePath string
	 GLoger *Zloger.Loger
	 WebPort string
	 BanIp map[string]int
	 SecretKey string
	 JwtExpire int
	 KeFuAcc map[string]string
}

func SetConfig() {
	Gconfig = &Sysconfig{
		Tkey:0,
	}
	Gconfig.Sysini =SysIni{
		LoginV:true,
		ConnectTimeOut: 60,
		LoginTimeOut: 30,
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		Gconfig.CurExePath =  "./"
	}else{
		Gconfig.CurExePath =  dir+"/"
	}

	fmt.Println("curpath:"+Gconfig.CurExePath)
	config, errr := Goconfig.LoadConfigFile(Gconfig.CurExePath + "config.ini")
	if errr!=nil{
		fmt.Println(errr)
		os.Exit(1)
	}

	kefulistsrt, _ := config.GetValue(Goconfig.DEFAULT_SECTION, "KeFuList")
	tmp := strings.Split(kefulistsrt,"|")
	Gconfig.KeFuAcc = make(map[string]string)
	for _,v:= range tmp {
		tmp2 := strings.Split(v,"&")
		if len(tmp2)==2 {
			Gconfig.KeFuAcc[tmp2[0]]=tmp2[1]
		}
	}
	fmt.Println("KeFu List:")
	fmt.Println(Gconfig.KeFuAcc)
	port, _ := config.GetValue(Goconfig.DEFAULT_SECTION, "Port")
	Gconfig.WebPort = "5555"
	if port!="" {
		Gconfig.WebPort = port
	}

	Gconfig.BanIp = make(map[string]int)
	Gconfig.SecretKey = "cptbtptp"
	Gconfig.JwtExpire = 60*60*24


	Gconfig.GLoger = Zloger.NewLog(Gconfig.CurExePath +"logs")


	ConMaster = conManage.NewManager()
}