package controllers

import (
	"NewsService/models"
	"encoding/json"
	"fmt"
	//	"strconv"
	"math"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	BaseController
}

func (this *LoginController) Login() {
	fmt.Println("点击登录")
	o := orm.NewOrm()
	l := new(models.Login)
	var user_info models.Login
	json.Unmarshal(this.Ctx.Input.RequestBody, &user_info)
	fmt.Println("user_info:", &user_info)
	n := user_info.Name
	p := user_info.Pwd
	query := o.QueryTable(l)
	exist := query.Filter("Name", n).Exist()
	if !exist {
		fmt.Println("用户不存在")
		this.ajaxMsg("用户不存在", MSG_ERR_Param)
	}
	err := query.Filter("Name", n).One(&user_info)
	if err != nil {
		this.ajaxMsg("登录失败", MSG_ERR_Resources)
	}
	pwd := user_info.Pwd
	if p != pwd {
		this.ajaxMsg("密码错误", MSG_ERR_Verified)
	}
	fmt.Println("用户名、密码正确")
	//授权
	//	token := this.GenToken(n)
	//	if token == "" {
	//		this.ajaxMsg("服务器内部错误", MSG_ERR)
	//	}
	token, i := this.Create_token(n, "gree")
	fmt.Println("token&time", token, i)

	//	success, err := this.Token_auth(token, "gree1")
	//	fmt.Println("验证:", success, err)

	list := make(map[string]interface{})
	list["name"] = n
	list["token"] = token
	list["pid"] = user_info.Pid
	this.ajaxList("登录成功", MSG_OK, 1, list)
}

func (this *LoginController) Add() {
	fmt.Println("用户添加")
	o := orm.NewOrm()
	list := make(map[string]interface{})
	var user models.Login
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	fmt.Println("user info", &user)

	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "gree")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)
	if user.Name == "" {
		fmt.Println("用户名不能为空")
		this.ajaxMsg("用户名不能为空", MSG_ERR_Param)
	}
	if user.Pwd == "" {
		fmt.Println("密码不能为空")
		this.ajaxMsg("密码不能为空", MSG_ERR_Param)
	}
	if user.Pid == 0 {
		fmt.Println("权限不能为空")
		this.ajaxMsg("权限不能为空", MSG_ERR_Param)
	}

	//	aid, err := strconv.ParseInt(this.GetRandomString(6), 10, 64)
	//	if err != nil {
	//		fmt.Println("cid string to int64 err")
	//	}
	//	ad.Aid = aid

	//insert
	_, err1 := o.Insert(&user)
	if err1 != nil {
		fmt.Println("inset err!", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}

	list["id"] = user.Id
	this.ajaxList("添加用户成功", MSG_OK, 1, list)
}

//编辑精选
func (this *LoginController) Edit() {
	fmt.Println("编辑用户")
	o := orm.NewOrm()
	var user models.Login
	//获取用户
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	fmt.Printf("user_info:", &user)

	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "gree")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	//cid
	id := user.Id
	if id == 0 {
		fmt.Printf("id 不能为空")
		this.ajaxMsg("id 不能为空", MSG_ERR_Param)
	}

	num, err := o.Update(&user)
	if err != nil {
		fmt.Println("更新用户信息失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	fmt.Println("updata user num", num)

	this.ajaxMsg("更新成功", MSG_OK)

}

//删除精选
func (this *LoginController) Del() {
	fmt.Println("删除用户")
	o := orm.NewOrm()
	//获取token
	var token models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	if token.Token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token.Token, "gree")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	//获取得到id
	var user_info models.Login
	json.Unmarshal(this.Ctx.Input.RequestBody, &user_info)
	id := user_info.Id
	if id == 1 {
		fmt.Println("管理员不能被删除")
		this.ajaxMsg("管理员不能被删除", MSG_ERR_Param)
	}
	if id == 0 {
		fmt.Println("获取id失败")
		this.ajaxMsg("id 不能为空", MSG_ERR_Param)
	}

	user := new(models.Login)
	num, err := o.QueryTable(user).Filter("Id", id).Delete()
	if err != nil {
		fmt.Println("删除用户失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	if num == 0 {
		fmt.Println("删除用户失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	this.ajaxMsg("删除成功", MSG_OK)

}

//获取精选
func (this *LoginController) GetData() {
	fmt.Println("获取用户")
	o := orm.NewOrm()
	user := new(models.Login)
	query := o.QueryTable(user)

	//获取token
	token := this.Input().Get("token")

	if token == "" {
		fmt.Println("token 为空")
		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	}

	name, err := this.Token_auth(token, "gree")
	if err != nil {
		fmt.Println("token err", err)
		this.ajaxMsg("token err!", MSG_ERR_Verified)
	}
	fmt.Println("当前访问用户为:", name)

	var maps []orm.Params

	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index为空")
	}
	//pagemax  一页多少
	pagemax, err := this.GetInt("pagemax")
	if err != nil {
		fmt.Println("获取每页数量为空")
	}
	//count
	count, err := query.Count()
	if err != nil {
		fmt.Println("获取数据总数为空")
		this.ajaxMsg("服务未知错误", MSG_ERR)
	}
	if pagemax != 0 {
		pagenum := int(math.Ceil(float64(count) / float64(pagemax)))

		if index > pagenum {
			//index = pagenum
			this.ajaxMsg("无法翻页了", MSG_ERR_Param)
		}
		fmt.Println("index&pagemax&pagenum", index, pagemax, pagenum)
	}
	query = query.Limit(pagemax, (index-1)*pagemax)

	_, err1 := query.OrderBy("-CreatTime").Values(&maps)
	if err1 != nil {
		fmt.Println("获取消息失败")
		this.ajaxMsg("获取消息失败", MSG_ERR_Resources)
	}
	this.ajaxList("获取用户消息成功", MSG_OK, count, maps)

}
