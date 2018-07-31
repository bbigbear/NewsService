<!DOCTYPE html>
<head>
  <title>Beego</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <link href="/static/css/bootstrap.min.css" rel="stylesheet" />
  <script src="/static/js/jquery-1.10.1.min.js"></script>
</head>
<body style="padding:20px;">
  	<div class="container">    
	    <form class="form-inline">
	        <div class="col-md-6 form-group">
	            <input id="sendbox" type="text" class="form-control" onkeydown="if(event.keyCode==13)return false;" required>
	        </div>
	        <button id="sendbtn" type="button" class="btn btn-default">发送</button>
	    </form>
	</div>
	<div class="container">	    
	    <ul id="chatbox">	        
	    </ul>
	</div>
	<script>
	 	var list =[];
		$('#sendbtn').on('click',function(){
			//alert("发送")
			$('#chatbox').append("<li>提问:"+$("#sendbox").val()+"</li>")
			var data	
			//发布
			var p = /[0-9]/;
			var b = p.test($("#sendbox").val())
			if (!b){
				list=[]
			}
			//console.log(b)
			if (list.length==0){
				data={
					'keyword':$("#sendbox").val()		
				};
			}else{
				//var num=$("#sendbox").val().replace(/[^0-9]/ig,"");
				//alert(num)
				data={
					'keyword':list[parseInt($("#sendbox").val())-1]
				}				
			}
				$.ajax({
					type:"POST",
					url:"/v1/getAnswer",
					data:data,
					error:function(request){ 
						alert("post error")						
					},
					success:function(res){
						if(res.code==200){
							//alert("保存成功")
							console.log(res)
							var num=res.count
							var i
							if(num==1){
								$('#chatbox').append("<li>"+res.data[0].Content+"</li>")
							}else{
								for(i=0;i<num;i++){
									var q=i+1
									$('#chatbox').append("<li>"+q+"、"+res.data[i].Title+"</li>")
									list.push(res.data[i].Title)
								}
							}
							//console.log(list)																					
						}else{
							//console.log(res)
							$('#chatbox').append("<li>我好像不明白，请重新输入关键词</li>")
						}						
					}
			  	});			
		});
	</script>
</body>
</html>
