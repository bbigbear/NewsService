package controllers

import (
	"NewsService/models"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	//	"strings"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//	"github.com/go-redis/redis"
)

type ChoiceController struct {
	BaseController
}

//添加精选
func (this *ChoiceController) Add() {
	fmt.Println("添加精选方案")
	o := orm.NewOrm()
	list := make(map[string]interface{})
	var choice models.Choice
	json.Unmarshal(this.Ctx.Input.RequestBody, &choice)
	fmt.Println("获取得到提交的数据,choice_info", &choice)

	//redis
	//	client := redis.NewClient(&redis.Options{
	//		Addr:     "127.0.0.1:6379",
	//		Password: "",
	//		DB:       0,
	//	})
	//	pong, err := client.Ping().Result()
	//	fmt.Println(pong, err)

	//	val, err := client.Get("cid").Result()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println("cid", val)

	//string to int
	//	cid, err := strconv.ParseInt(val, 10, 64)
	//	if err != nil {
	//		fmt.Println("cid to int err!", err)
	//	}

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
	//c:=this.GetRandomString(6)
	//	fmt.Println("choice", choice.Title)
	if choice.Title == "" {
		fmt.Println("post err!")
		this.ajaxMsg("post err", MSG_ERR_Param)
	}

	cid, err := strconv.ParseInt(this.GetRandomString(6), 10, 64)
	if err != nil {
		fmt.Println("cid string to int64 err")
	}
	choice.Cid = cid
	//	if choice.ArticlePic != "" {
	//		articlePicList := strings.Split(choice.ArticlePic, ":")
	//		choice.ArticlePic = beego.AppConfig.String("img_upload_url") + articlePicList[2]
	//	}

	//	if choice.SourcePic != "" {
	//		sourcePicList := strings.Split(choice.SourcePic, ":")
	//		choice.SourcePic = beego.AppConfig.String("img_upload_url") + sourcePicList[2]
	//	}
	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("新建失败")
	}
	choice.CreatTime = nowtime
	choice.Url = beego.AppConfig.String("vue_url") + choice.Url + strconv.FormatInt(choice.Cid, 10)
	//insert
	_, err1 := o.Insert(&choice)
	if err1 != nil {
		fmt.Println("inset err!", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}

	list["cid"] = choice.Cid
	this.ajaxList("添加精选成功", MSG_OK, 1, list)

}

//编辑精选
func (this *ChoiceController) Edit() {
	fmt.Println("编辑精选")
	o := orm.NewOrm()
	var choice models.Choice
	//获取编辑的精选
	json.Unmarshal(this.Ctx.Input.RequestBody, &choice)
	fmt.Printf("choice_info:", &choice)

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
	cid := choice.Cid
	if cid == 0 {
		fmt.Printf("cid 不能为空")
		this.ajaxMsg("cid 不能为空", MSG_ERR_Param)
	}
	choice.ReadNum = choice.ReadNum + 1
	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("新建失败")
	}
	choice.CreatTime = nowtime
	//	if choice.ArticlePic != "" {
	//		articlePicList := strings.Split(choice.ArticlePic, ":")
	//		choice.ArticlePic = beego.AppConfig.String("img_upload_url") + articlePicList[2]
	//	}

	//	if choice.SourcePic != "" {
	//		sourcePicList := strings.Split(choice.SourcePic, ":")
	//		choice.SourcePic = beego.AppConfig.String("img_upload_url") + sourcePicList[2]
	//	}
	choice.Url = beego.AppConfig.String("vue_url") + "/choices/" + strconv.FormatInt(choice.Cid, 10)
	//更新
	num, err := o.Update(&choice)
	if err != nil {
		fmt.Println("更新精选信息失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	fmt.Println("updata choice num", num)

	this.ajaxMsg("更新成功", MSG_OK)

}

//删除精选
func (this *ChoiceController) Del() {
	fmt.Println("删除精选")
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

	//获取得到cid
	var choice_info models.Choice
	json.Unmarshal(this.Ctx.Input.RequestBody, &choice_info)
	//cid, err := this.GetInt("cid")
	//	if err != nil {
	//		fmt.Println("获取cid失败")
	//		this.ajaxMsg("cid 不能为空", MSG_ERR_Param)
	//	}
	cid := choice_info.Cid
	if cid == 0 {
		fmt.Println("获取cid失败")
		this.ajaxMsg("cid 不能为空", MSG_ERR_Param)
	}
	choice := new(models.Choice)
	num, err := o.QueryTable(choice).Filter("Cid", cid).Delete()
	if err != nil {
		fmt.Println("删除精选失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	if num == 0 {
		fmt.Println("删除精选失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	this.ajaxMsg("删除成功", MSG_OK)

}

//获取精选
func (this *ChoiceController) GetData() {
	fmt.Println("获取精选内容")
	o := orm.NewOrm()
	choice := new(models.Choice)
	query := o.QueryTable(choice)

	//获取token
	//	var token models.Token
	//	json.Unmarshal(this.Ctx.Input.RequestBody, &token)
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

	//type
	//	typ, err := this.GetInt("typ")
	//	if err != nil {
	//		fmt.Println("获取type失败")
	//		this.ajaxMsg("type不能为空", MSG_ERR_Param)
	//	}
	//	if typ == 1 {
	//		fmt.Println("获取单个精选文章内容", typ)
	//		//cid
	//		cid, err := this.GetInt("cid")
	//		if err != nil {
	//			fmt.Println("cid 错误")
	//			this.ajaxMsg("参数cid 错误", MSG_ERR_Param)
	//		}
	//		if cid == 0 {
	//			fmt.Println("cid 不能为空")
	//			this.ajaxMsg("cid 不能为空", MSG_ERR_Param)
	//		}
	//		query = query.Filter("Cid", cid)
	//	} else if typ == 0 {
	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index下标错误")
		//this.ajaxMsg("参数index错误", MSG_ERR_Param)
	}
	//	if index == 0 {
	//		fmt.Println("获取index下标为空")
	//		this.ajaxMsg("参数index不能为空", MSG_ERR_Param)
	//	}
	//pagemax  一页多少
	pagemax, err := this.GetInt("pagemax")
	if err != nil {
		fmt.Println("获取每页数量为空")
		//this.ajaxMsg("参数pagemax错误", MSG_ERR_Param)
	}
	//	if pagemax == 0 {
	//		fmt.Println("获取每页数量为空")
	//		this.ajaxMsg("参数pagemax不能为空", MSG_ERR_Param)
	//	}
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
	//	} else {
	//		this.ajaxMsg("typ 未定义", MSG_ERR_Param)
	//	}
	_, err1 := query.OrderBy("-CreatTime").Values(&maps, "Cid", "Title", "Brief", "Content", "Tag", "ArticlePic", "Url", "ReadNum", "Source", "SourcePic")
	if err1 != nil {
		fmt.Println("获取消息失败")
		this.ajaxMsg("获取消息失败", MSG_ERR_Resources)
	}
	this.ajaxList("获取精选消息成功", MSG_OK, count, maps)

}

//获取精选详情
func (this *ChoiceController) GetDetail() {
	fmt.Println("获取精选详情")
	o := orm.NewOrm()
	choice := new(models.Choice)
	query := o.QueryTable(choice)

	//获取token
	//	token := this.Input().Get("token")

	//	if token == "" {
	//		fmt.Println("token 为空")
	//		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	//	}

	//	name, err := this.Token_auth(token, "gree")
	//	if err != nil {
	//		fmt.Println("token err", err)
	//		this.ajaxMsg("token err!", MSG_ERR_Verified)
	//	}
	//	fmt.Println("当前访问用户为:", name)

	var maps []orm.Params

	//cid
	cid, err := this.GetInt("cid")
	if err != nil {
		fmt.Println("cid 错误")
		this.ajaxMsg("参数cid 错误", MSG_ERR_Param)
	}
	if cid == 0 {
		fmt.Println("cid 不能为空")
		this.ajaxMsg("cid 不能为空", MSG_ERR_Param)
	}

	num, err := query.Filter("Cid", cid).Values(&maps, "Cid", "Title", "Brief", "Content", "Tag", "ArticlePic", "Url", "ReadNum", "Source", "SourcePic")
	if err != nil {
		fmt.Println("获取消息失败")
		this.ajaxMsg("不存在该文章", MSG_ERR_Resources)
	}
	for _, m := range maps {
		m["ReadNum"] = m["ReadNum"].(int64) + 1
		//更新阅读量
		num, err := o.QueryTable(choice).Filter("Cid", cid).Update(orm.Params{
			"ReadNum": m["ReadNum"],
		})
		if err != nil {
			fmt.Println("err!")
		}
		fmt.Println("update readnum", num)
	}
	this.ajaxList("获取精选消息成功", MSG_OK, num, maps)

}

//api web
//获取精选
func (this *ChoiceController) GetArticle() {
	fmt.Println("获取精选内容")
	o := orm.NewOrm()
	choice := new(models.Choice)
	query := o.QueryTable(choice)

	var maps []orm.Params

	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index下标为空")
		//this.ajaxMsg("参数index错误", MSG_ERR_Param)
	}

	//pagemax  一页多少
	pagemax, err := this.GetInt("pagemax")
	if err != nil {
		fmt.Println("获取每页数量为空")
		//this.ajaxMsg("参数pagemax错误", MSG_ERR_Param)
	}

	//count
	count, err := query.Count()
	if err != nil {
		fmt.Println("获取数据总数错误")
		this.ajaxMsg("服务未知错误", MSG_ERR)
	}
	pagenum := int(math.Ceil(float64(count) / float64(pagemax)))
	if pagemax != 0 {
		if index > pagenum {
			//index = pagenum
			this.ajaxMsg("无法翻页了", MSG_ERR_Param)
		}
		fmt.Println("index&pagemax&pagenum", index, pagemax, pagenum)
	}
	query = query.Limit(pagemax, (index-1)*pagemax)
	//	} else {
	//		this.ajaxMsg("typ 未定义", MSG_ERR_Param)
	//	}
	num, err := query.OrderBy("-CreatTime").Values(&maps, "Cid", "Title", "Brief", "Tag", "ArticlePic", "Url", "ReadNum", "CreatTime", "Source", "SourcePic")
	if err != nil {
		fmt.Println("获取精选内容失败")
		this.ajaxMsg("获取精选内容失败", MSG_ERR_Param)
	}
	fmt.Println("当前获取的数量", num)

	for _, m := range maps {
		m["CreatTime"] = m["CreatTime"].(time.Time).Format("2006-01-02 15:04:05")
	}

	this.ajaxList("获取精选内容成功", MSG_OK, count, maps)

}

//func (this *AdController) GetCollect() {

//	//获取Cid []string
//	cid := this.GetInt()
//	o:=orm.NewOrm()
//	choice:=new(models.Choice)
//	o.QueryTable(choice).Filter("Cid_")

//}
