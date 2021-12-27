$(document).ready(function () {
    var idpre='hbywebim_';
    var title = 'HBY-WebIm';
    //websocket地址
    var wshost = 'ws://127.0.0.1:5555';
    //api域名
    var webhost ='';



    var ws;
    var wsstatus = false;
    var loginflag = false;
    var weTimer = null;
    var CurChatUser =undefined

    //创建右下角按钮
    $('body').append('<div id="'+idpre+'btn">客</div>')
    $('body').append('<div id="'+idpre+'chat">' +
        '<h1 class="title">'+title+'</h1>' +
        '<div><span>当前客服：</span><span id="kefuname"></span></div>' +
        '<div class="chat-main"><div class="contentdiv"></div>' +
        '<div class="operadiv">' +
        '<div><textarea rows="3" id="sendcontent" placeholder="消息"></textarea></div>\n' +
        '<div class="chatoperadiv" style="margin-top:5px;"><button type="button" class="hidechat">隐藏</button><button type="button" class="sendmsg">发送</button></div>' +
        '</div>' +
        '</div></div>')


    //登录
    function checklogin() {
        if (!loginflag) {
            //自动登录
            ws.send(JSON.stringify({
                User: Date.parse(new Date()) + "",
                Passwd: ""
            }));
        }
    }
    function aflogin() {
        loginflag = true
        weTimer = setInterval(function () {
            ws.send("11")
        },30000)
    }
    function getuserlist() {
        //采用ws方式获取解决跨域问题
        ws.send("kefulist")
    }

    function getFullTime() {
        let date = new Date(),//时间戳为10位需*1000，时间戳为13位的话不需乘1000
            Y = date.getFullYear() + '',
            M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1),
            D = (date.getDate() < 10 ? '0'+(date.getDate()) : date.getDate()),
            h = (date.getHours() < 10 ? '0'+(date.getHours()) : date.getHours()),
            m = (date.getMinutes() < 10 ? '0'+(date.getMinutes()) : date.getMinutes()),
            s = (date.getSeconds() < 10 ? '0'+(date.getSeconds()) : date.getSeconds());
        return h +':' + m +':' + s
    }
    function apprecmsg(fuser,msg,timestr) {
        var tmparr = timestr.split(' ')
        var str ='<div class="omsgdiv">\n' +
            '<p><span>'+tmparr[1]+'</span><span class="username">'+fuser+'</span></p>\n' +
            '<p class="content"><pre>'+msg+'</pre></p>\n' +
            '</div>';
        $('#'+idpre+'chat .contentdiv').append(str);
        $('#'+idpre+'chat .contentdiv').scrollTop( $('#'+idpre+'chat .contentdiv')[0].scrollHeight)
    }
    function appsendmsg(msg) {
        var str ='<div class="somsgdiv">' +
            '<p><span>'+getFullTime()+'</span><span class="username">我</span></p>' +
            '<p class="content"><pre>'+msg+'</pre></p>' +
            '</div>';
        $('#'+idpre+'chat .contentdiv').append(str);
        $('#'+idpre+'chat .contentdiv').scrollTop( $('#'+idpre+'chat .contentdiv')[0].scrollHeight)
    }


    function WebSocketTest()
    {
        if ("WebSocket" in window)
        {
            console.log("您的浏览器支持 WebSocket!");
            // 打开一个 web socket
            ws = new WebSocket(wshost +"/websocket");
            ws.onopen = function()
            {
                wsstatus = true
                checklogin()
                console.log("开启成功...");
            };
            ws.onmessage = function (evt)
            {
                var received_msg = evt.data;
                var recmsg = eval('(' + received_msg + ')');
                if(recmsg.code){
                    switch (recmsg.code) {
                        case "1001":
                            aflogin()
                            getuserlist()
                            break;
                        case "1002":
                            alert(recmsg.msg)
                            break;
                        case "1003":
                            alartmsg(recmsg.msg)
                            break;
                        case "2001":
                        case "2002":
                            //上下线
                            getuserlist()
                            break
                        case "200":
                            apprecmsg(recmsg.from,recmsg.msg,recmsg.time)
                            break;
                        case "4":
                            alartmsg(recmsg.msg)
                            break;
                        case "5001":
                            dealUserListMsg(recmsg);
                            break;
                        default:
                            alert('未定义消息')
                            break;
                    }
                }
            };

            ws.onclose = function()
            {
                wsstatus = false
                // 关闭 websocket
                console.log("连接已关闭...");
                alert('连接已关闭')
            };
        }

        else
        {
            // 浏览器不支持 WebSocket
            alert("您的浏览器不支持 WebSocket!");
        }
    }

    function dealUserListMsg(res) {
        var str = '';
        if(res.data!="null"){
            for(var i in res.data){
                str =res.data[i].name
                CurChatUser=res.data[i].id
                break
            }
        }
        if (str=="") {
            str = '无客服在线'
        }
        console.log(str)
        console.log($('#kefuname'))
        $('#kefuname').html(str);
    }
    function sendmsg() {
        if(!wsstatus){
            WebSocketTest()
        }else{
            // 获取消息接收人
            var sendto = getNowChat()
            if(!sendto){
                alert("请点击用户列表选择聊天对象!");
                return
            }
            var content =     $('#'+idpre+'chat #sendcontent').val();
            if(content){
                $('#'+idpre+'chat #sendcontent').val('');
                appsendmsg(content);
                ws.send(JSON.stringify({
                    code:"200",
                    from:"test",
                    msg: content,
                    to:sendto
                }));
            }
        }
    }

    function  getNowChat() {
        return CurChatUser;
    }
    $('#'+idpre+'btn').click(function () {
        $(this).hide();
        $('#'+idpre+'chat').show();
    });
    $('#'+idpre+'chat').on('click','.sendmsg',function () {
        console.log(11);
        sendmsg();
    })
    $('#'+idpre+'chat').on('click','.hidechat',function () {
        $('#'+idpre+'chat').hide();
        $('#'+idpre+'btn').show();
    })

    WebSocketTest();

})