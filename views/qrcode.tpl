<!DOCTYPE html>
<html lang="en" style="height: 100%">
<head>
<!--    <link rel="stylesheet" href="css/main.css" type="text/css"/>-->
    <meta name="viewport" content="height=device-height,width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no">
    <meta charset="UTF-8">
    <title>格力+</title>
<!--    <script type="text/javascript" src="scripts/jquery.js"></script>-->
<!--<script src="/static/js/jquery-1.10.1.min.js"></script>-->
<script src="https://apps.bdimg.com/libs/jquery/2.1.4/jquery.min.js">
</script>
</head>

<body style="height: 100%;background:-webkit-gradient(linear, 0 0, 0 100%, from(#8fc0d6), to(#d2ebf6));">
<div id="header">
    <div id="img">
        <img src="img/gp_cn_icon.png" width=45% height=65%>
    </div>
</div>
<div id="pic">
    <div id="content_google">
        <p id="text_google"></p>
        <img src="img/cn_android.png" width=60% height=50% onclick="isWetChat_Google()">
    </div>
    <div id="contont_ios">
        <p id="text_ios"></p>
        <img src="img/cn_ios.png" width=60% height=50% onclick="isWetChat_ios()">
    </div >
    <div>
        <a href="https://itunes.apple.com/app/id1234475684?mt=8" id="openApp" style="display: none">ios</a>
        <a href="http://grih.gree.com/app/GetAppLastVersion?name=com.gree.greesmarthome" id="openAndroid" style="display: none">android</a>
    </div>
    <div>
        <div id="cover">
            <div id="wechat">
                <img src="img/wechat.png" width=80% height=60% >
            </div>
        </div>
    </div>
</div>
<script>
	//自动加载
	$(function(){
		//$.get("http://grih.gree.com/app/AppLastVersion?name=com.gree.greesmarthome&lang=cn", function(result){
		//	console.log(result)
		//  });
		$.getJSON("http://grih.gree.com/app/AppLastVersion?name=com.gree.greesmarthome&lang=cn&callback=?",function(json){ 
		//要求远程请求页面的数据格式为： ?(json_data) //例如： //?([{"_name":"湖南省","_regionId":134},{"_name":"北京市","_regionId":143}]) alert(json[0]._name); 		
			console.log(json)
		});		
	});
	
	function m1(data) {
     console.log(data)
 	}
</script>
</body>
</html>