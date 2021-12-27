var loginflag = false;
var ws;
var webhost = '';
var wshost = 'ws://127.0.0.1:5555';
var wsstatus = false;

var jwt ='';
var jwtkey = 'hbyloginkey'
var weTimer = null;

var CurChatUser =undefined
var Guserlist =[]
function alartmsg(str) {
    alert(str);
}
function getFullTime() {
    let date = new Date(),//时间戳为10位需*1000，时间戳为13位的话不需乘1000
        Y = date.getFullYear() + '',
        M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1),
        D = (date.getDate() < 10 ? '0'+(date.getDate()) : date.getDate()),
        h = (date.getHours() < 10 ? '0'+(date.getHours()) : date.getHours()),
        m = (date.getMinutes() < 10 ? '0'+(date.getMinutes()) : date.getMinutes()),
        s = (date.getSeconds() < 10 ? '0'+(date.getSeconds()) : date.getSeconds());
    return Y +'-' +M +'-' + D +' ' +h +':' + m +':' + s
}
function getuserlist() {
    $.ajax({
        url: webhost +'/userlist',
        type:'get',
        data:{
            'type':'kefu'
        },
        dataType: 'json',
        success:function (res) {
            if(res.code=="200"){
                var str = '';
                var onlinenum = 0;
                if(res.data!="null"){
                    for(var i in res.data){
                        onlinenum++;
                        Guserlist[res.data[i].name] =res.data[i].id
                        str +=' <li data-id='+res.data[i].id+' id="userli'+res.data[i].id+'">\n' +
                            '<img src="https://s1.ax1x.com/2020/05/25/tC039K.jpg">\n' +
                            '<span class="badge" style="display: none"></span><span>'+res.data[i].name+'</span>\n' +
                            '</li>';
                    }
                }
                if (str=="") {
                    str = '无客服在线'
                }
                $('#userlist ul').html(str);
                $('#onlinenum').text(onlinenum);
            }else{
                alartmsg(res.msg)
            }
        },
        error:function (e) {
            console.log(e)
            alartmsg('获取用户列表失败')
        }
    })
}
function getemojilist() {
    $.ajax({
        url: webhost +'/emojilist?path=emoji',
        type:'get',
        dataType: 'json',
        success:function (res) {
            if(res.code=="200"){
                var str = '';
                if(res.data!="null"){
                    for(var i in res.data){
                        str +=' <li data-code='+res.data[i].code+'>\n' +
                            '<img src="'+res.data[i].path+'" title="'+res.data[i].name+'">\n' +
                            '</li>';
                    }
                }
                $('#emojidiv ul').html(str);
            }else{
                alartmsg(res.msg)
            }
        },
        error:function (e) {
            console.log(e)
            alartmsg('获取表情列表失败')
        }
    })
}
function apprecmsg(fuser,msg,timestr) {
    var chatid = Guserlist[fuser]
    if(!chatid){
        return
    }
    if(chatid!=CurChatUser){
        //红点提示
        $('#userli'+chatid).children('.badge').show();
    }
    if($("#contentdiv"+chatid).length==0) {
        //新增
        $('#chatmain').prepend("<div class=\"col-sm-9  col-xs-12 contentdiv\" id=\"contentdiv" + chatid + "\" style='display: none'></div>");
    }
    var str ='<div class="omsgdiv">\n' +
        '<p><span>'+timestr+'</span><span class="username">'+fuser+'</span></p>\n' +
        '<p class="content"><pre>'+msg+'</pre></p>\n' +
        '</div>';
    $('#contentdiv'+chatid).append(str);
    $('#contentdiv'+chatid).scrollTop($('#contentdiv'+chatid)[0].scrollHeight)
}
function appsendmsg(msg) {
    var str ='<div class="somsgdiv">' +
        '<p><span>'+getFullTime()+'</span><span class="username">我</span></p>' +
        '<p class="content"><pre>'+msg+'</pre></p>' +
        '</div>';
    $('#contentdiv'+CurChatUser).append(str);
    $('#contentdiv'+CurChatUser).scrollTop($('#contentdiv'+CurChatUser)[0].scrollHeight)
}

function aflogin() {
    loginflag = true
    weTimer = setInterval(function () {
        ws.send("11")
    },30000)
}
function afloginout() {
    loginflag = false
    if(weTimer){
        clearInterval(weTimer)
    }
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
            console.log(received_msg);
            var recmsg = eval('(' + received_msg + ')');
            if(recmsg.code){
                switch (recmsg.code) {
                    case "1001":
                        $('#loginModel').modal('hide')
                        $('#loginbtn').button('reset');
                        jwt = recmsg.msg
                        localStorage.setItem(jwtkey,jwt)
                        aflogin()
                        getuserlist()
                        break;
                    case "1002":
                        $('#loginbtn').button('reset');
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
            afloginout()

        };
    }

    else
    {
        // 浏览器不支持 WebSocket
        alert("您的浏览器不支持 WebSocket!");
    }
}
WebSocketTest()
getemojilist()
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
        var content = $("#sendcontent").val();
        if(content){
            $("#sendcontent").val('');
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
function checklogin() {
    if(!loginflag){
        //自动登录
        ws.send(JSON.stringify({
            User: Date.parse(new Date()) +"",
            Passwd:""
        }));
    }
}

$('#emojidiv').on('click','li',function () {
    var code = $(this).attr('data-code');
    var content = $('#sendcontent').val();
    $('#sendcontent').val(content+'[:'+code+':]');
    $('#sendcontent').focus();
});
$('#loginbtn').click(function(){
    if(!$('#name').val()){
        $('#name').focus()
        return
    }
    if(!wsstatus){
        WebSocketTest()
        return
    }
    $(this).button('loading');
    ws.send(JSON.stringify({
        User: $('#name').val(),
        Passwd:""
    }));
})

$('#userlist').on('click','li',function () {
    $('#userlist li').removeClass("on");
    $(this).addClass('on');
    CurChatUser = $(this).attr('data-id') +"";
    $(this).children('.badge').hide();
    $('#kefuname').text( $('#userlist li span:last').text())
    changeChatMain(CurChatUser)
})

function changeChatMain(id) {
    $(".contentdiv").hide();
    if($("#contentdiv"+id).length!=0){
        $("#contentdiv"+id).show();
    }else{
        //新增
        $('#chatmain').prepend("<div class=\"col-sm-9  col-xs-12 contentdiv\" id=\"contentdiv"+id+"\" ></div>");
    }
}

function showkefu() {
    if($('#userlist').css('display') == 'none'){
        $('#userlist').show();
    }else{
        $('#userlist').hide();
    }

}
