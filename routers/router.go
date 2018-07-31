package routers

import (
	"NewsService/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/v1/login", &controllers.LoginController{}, "post:Login")
	beego.Router("/v1/put_file", &controllers.BaseController{}, "*:PutFile")

	//智能客服
	beego.Router("/v1/chat/add", &controllers.MainController{}, "post:Add")
	beego.Router("/v1/chat/edit", &controllers.MainController{}, "post:Edit")
	beego.Router("/v1/chat/del", &controllers.MainController{}, "post:Del")
	beego.Router("/v1/chat/getdata", &controllers.MainController{}, "*:GetData")
	beego.Router("/v1/getAnswer", &controllers.MainController{}, "*:GetAnswer")

	//精选
	beego.Router("/v1/choice/add", &controllers.ChoiceController{}, "post:Add")
	beego.Router("/v1/choice/edit", &controllers.ChoiceController{}, "post:Edit")
	beego.Router("/v1/choice/del", &controllers.ChoiceController{}, "post:Del")
	beego.Router("/v1/choice/getdata", &controllers.ChoiceController{}, "*:GetData")
	beego.Router("/v1/choice/getdetail", &controllers.ChoiceController{}, "*:GetDetail")
	//小窍门
	beego.Router("/v1/tips/add", &controllers.TipsController{}, "post:Add")
	beego.Router("/v1/tips/edit", &controllers.TipsController{}, "post:Edit")
	beego.Router("/v1/tips/del", &controllers.TipsController{}, "post:Del")
	beego.Router("/v1/tips/getdata", &controllers.TipsController{}, "*:GetData")
	beego.Router("/v1/tips/getdetail", &controllers.TipsController{}, "*:GetDetail")
	//广告
	beego.Router("/v1/ad/add", &controllers.AdController{}, "post:Add")
	beego.Router("/v1/ad/edit", &controllers.AdController{}, "post:Edit")
	beego.Router("/v1/ad/del", &controllers.AdController{}, "post:Del")
	beego.Router("/v1/ad/getdata", &controllers.AdController{}, "*:GetData")
	beego.Router("/v1/ad/getproducts", &controllers.AdController{}, "*:GetProductStyle")
	beego.Router("/v1/ad/push", &controllers.AdController{}, "post:NewAdPush")
	//用户
	beego.Router("/v1/user/add", &controllers.LoginController{}, "post:Add")
	beego.Router("/v1/user/edit", &controllers.LoginController{}, "post:Edit")
	beego.Router("/v1/user/del", &controllers.LoginController{}, "post:Del")
	beego.Router("/v1/user/getdata", &controllers.LoginController{}, "*:GetData")

	//webapi
	beego.Router("/api/choice/getdata", &controllers.ChoiceController{}, "*:GetArticle")
	beego.Router("/api/tips/getdata", &controllers.TipsController{}, "*:GetTips")
	beego.Router("/api/ad/getdata", &controllers.AdController{}, "*:GetAd")

}
