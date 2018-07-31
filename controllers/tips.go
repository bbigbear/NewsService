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

type TipsController struct {
	BaseController
}

//添加小窍门
func (this *TipsController) Add() {
	fmt.Println("添加小窍门方案")
	o := orm.NewOrm()
	list := make(map[string]interface{})
	var tips models.Tips
	json.Unmarshal(this.Ctx.Input.RequestBody, &tips)
	fmt.Println("获取得到提交的数据,choice_info", &tips)
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
	if tips.Title == "" {
		fmt.Println("post err!")
		this.ajaxMsg("post err", MSG_ERR_Param)
	}

	tid, err := strconv.ParseInt(this.GetRandomString(6), 10, 64)
	if err != nil {
		fmt.Println("tid string to int64 err")
	}
	tips.Tid = tid
	//	if tips.ArticlePic != "" {
	//		articlePicList := strings.Split(tips.ArticlePic, ":")
	//		tips.ArticlePic = beego.AppConfig.String("img_upload_url") + articlePicList[2]
	//	}
	//	if tips.SourcePic != "" {
	//		sourcePicList := strings.Split(tips.SourcePic, ":")
	//		tips.SourcePic = beego.AppConfig.String("img_upload_url") + sourcePicList[2]
	//	}
	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("新建失败")
	}
	tips.CreatTime = nowtime

	tips.Url = beego.AppConfig.String("vue_url") + tips.Url + strconv.FormatInt(tips.Tid, 10)
	//insert
	_, err1 := o.Insert(&tips)
	if err1 != nil {
		fmt.Println("inset err!", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}

	list["tid"] = tips.Tid
	this.ajaxList("添加小窍门成功", MSG_OK, 1, list)

}

//编辑小窍门
func (this *TipsController) Edit() {
	fmt.Println("编辑小窍门")
	o := orm.NewOrm()
	var tips models.Tips
	//获取编辑的小窍门
	json.Unmarshal(this.Ctx.Input.RequestBody, &tips)
	fmt.Printf("tip_info:", &tips)

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
	tid := tips.Tid
	if tid == 0 {
		fmt.Printf("tid 不能为空")
		this.ajaxMsg("tid 不能为空", MSG_ERR_Param)
	}
	tips.ReadNum = tips.ReadNum + 1
	//	if tips.ArticlePic != "" {
	//		articlePicList := strings.Split(tips.ArticlePic, ":")
	//		tips.ArticlePic = beego.AppConfig.String("img_upload_url") + articlePicList[2]
	//	}
	//	if tips.SourcePic != "" {
	//		sourcePicList := strings.Split(tips.SourcePic, ":")
	//		tips.SourcePic = beego.AppConfig.String("img_upload_url") + sourcePicList[2]
	//	}
	tips.Url = beego.AppConfig.String("vue_url") + "/littletips/" + strconv.FormatInt(tips.Tid, 10)
	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("新建失败")
	}
	tips.CreatTime = nowtime
	//更新
	num, err := o.Update(&tips)
	if err != nil {
		fmt.Println("更新小窍门信息失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	fmt.Println("updata choice num", num)

	this.ajaxMsg("更新成功", MSG_OK)

}

//删除小窍门
func (this *TipsController) Del() {
	fmt.Println("删除小窍门")
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

	//获取得到tid
	var tip models.Tips
	json.Unmarshal(this.Ctx.Input.RequestBody, &tip)
	tid := tip.Tid
	//	tid, err := this.GetInt("tid")
	//	if err != nil {
	//		fmt.Println("获取tid失败")
	//		this.ajaxMsg("tid 不能为空", MSG_ERR_Param)
	//	}
	if tid == 0 {
		fmt.Println("获取tid失败")
		this.ajaxMsg("tid 不能为空", MSG_ERR_Param)
	}
	tips := new(models.Tips)
	num, err := o.QueryTable(tips).Filter("Tid", tid).Delete()
	if err != nil {
		fmt.Println("删除小窍门失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	if num == 0 {
		fmt.Println("删除小窍门失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	this.ajaxMsg("删除成功", MSG_OK)

}

//获取小窍门
func (this *TipsController) GetData() {
	fmt.Println("获取小窍门内容")
	o := orm.NewOrm()
	tips := new(models.Tips)
	query := o.QueryTable(tips)

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

	//type
	//	typ, err := this.GetInt("typ")
	//	if err != nil {
	//		fmt.Println("获取type失败")
	//		this.ajaxMsg("type不能为空", MSG_ERR_Param)
	//	}
	//	if typ == 1 {
	//		fmt.Println("获取单个小窍门文章内容", typ)
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
		fmt.Println("获取每页数量错误")
		//this.ajaxMsg("参数pagemax错误", MSG_ERR_Param)
	}
	//	if pagemax == 0 {
	//		fmt.Println("获取每页数量为空")
	//		this.ajaxMsg("参数pagemax不能为空", MSG_ERR_Param)
	//	}
	//count
	count, err := query.Count()
	if err != nil {
		fmt.Println("获取数据总数错误")
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
	_, err1 := query.OrderBy("-CreatTime").Values(&maps, "Tid", "Title", "Content", "Tag", "ArticlePic", "Url", "ReadNum", "Source", "SourcePic")
	if err1 != nil {
		fmt.Println("获取消息失败")
		this.ajaxMsg("获取消息失败", MSG_ERR_Resources)
	}
	this.ajaxList("获取小窍门消息成功", MSG_OK, count, maps)

}

//获取小窍门详情
func (this *TipsController) GetDetail() {
	fmt.Println("获取小窍门详情")
	o := orm.NewOrm()
	tips := new(models.Tips)
	query := o.QueryTable(tips)

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

	//tid
	tid, err := this.GetInt("tid")
	if err != nil {
		fmt.Println("tid 错误")
		this.ajaxMsg("参数tid 错误", MSG_ERR_Param)
	}
	if tid == 0 {
		fmt.Println("tid 不能为空")
		this.ajaxMsg("tid 不能为空", MSG_ERR_Param)
	}

	num, err := query.Filter("Tid", tid).Values(&maps, "Tid", "Title", "Content", "Tag", "ArticlePic", "Url", "ReadNum", "Source", "SourcePic")
	if err != nil {
		fmt.Println("获取消息失败")
		this.ajaxMsg("不存在该文章", MSG_ERR_Resources)
	}

	for _, m := range maps {
		m["ReadNum"] = m["ReadNum"].(int64) + 1
		//更新阅读量
		num, err := o.QueryTable(tips).Filter("Tid", tid).Update(orm.Params{
			"ReadNum": m["ReadNum"],
		})
		if err != nil {
			fmt.Println("err!")
		}
		fmt.Println("update readnum", num)
	}

	this.ajaxList("获取小窍门消息成功", MSG_OK, num, maps)

}

//api web
//获取小窍门
func (this *TipsController) GetTips() {
	fmt.Println("获取小窍门内容")
	o := orm.NewOrm()
	tips := new(models.Tips)
	query := o.QueryTable(tips)

	var maps []orm.Params

	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index下标错误")
		//this.ajaxMsg("参数index错误", MSG_ERR_Param)
	}

	//pagemax  一页多少
	pagemax, err := this.GetInt("pagemax")
	if err != nil {
		fmt.Println("获取每页数量错误")
		//this.ajaxMsg("参数pagemax错误", MSG_ERR_Param)
	}

	//count
	count, err := query.Count()
	if err != nil {
		fmt.Println("获取数据总数错误")
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
	_, err1 := query.OrderBy("-CreatTime").Values(&maps, "Tid", "Title", "ArticlePic", "Url", "ReadNum", "CreatTime", "Source", "SourcePic")
	if err1 != nil {
		fmt.Println("获取小窍门内容失败")
		this.ajaxMsg("获取小窍门内容失败", MSG_ERR_Param)
	}
	for _, m := range maps {
		m["CreatTime"] = m["CreatTime"].(time.Time).Format("2006-01-02 15:04:05")
	}
	this.ajaxList("获取小窍门内容成功", MSG_OK, count, maps)

}
