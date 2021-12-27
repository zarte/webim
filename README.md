# WebIm 
Supply Customer Service 
## 界面预览

## 快速上手
拷贝install目录下文件 
### 本地部署模式
配置修改  
client.js&kefu.js     
var webhost = 'web路径';    
var wshost = 'websocket路径';    
client.html&kefu.html  
对应js与css路径  
  
启动  ./webim.exe  
游客访问:  
http://127.0.0.1:5555/  
客服访问：  
http://127.0.0.1:5555/kefu  
预置账号  
test22   密码112233  

### 插件模式
参考plugtest.html页面  
依赖jq库需先引入  
```js
<script src="http://libs.baidu.com/jquery/2.0.0/jquery.min.js"></script>
```
引入js与css文件,同时修改plug.js里wshost参数
```js
<link rel="stylesheet" href="/public/css/plug.css"/>
<script src="/public/js/plug.js"></script>
```
## 配置
config.ini  
KeFuList:客服账号列表|分隔账号&分割账号与密码
Port:监听端口

## 二次开发备注
用户： 客服，游客  
客服登录：  参数：用户名，密码  获取所有用户列表  
游客登录： 参数：自动生成唯一用户名（不可与客服重复） 获取在线客服列表  
客服调用api需带上token参数  
默认60秒无消息传输自动断开链接  
登录错误过多自动banip  

## code

1001: 登录成功  
1002：登录失败   
1003:  断开连接（仅通知接口）  
2001: 上线通知  
2002: 下线通知  
5001: 用户列表   

## 更新日志
v1.0 211225  
v1.1 211227  
插件模式  
通过ws获取客服列表解决跨域问题  
## todo
插件版表情支持
