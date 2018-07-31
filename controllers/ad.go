package controllers

import (
	"NewsService/models"
	//	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	//	"github.com/go-redis/redis"
)

type AdController struct {
	BaseController
}

func (this *AdController) Add() {
	fmt.Println("添加广告")
	o := orm.NewOrm()
	list := make(map[string]interface{})
	var ad models.Ad
	json.Unmarshal(this.Ctx.Input.RequestBody, &ad)
	fmt.Println("获取广告添加的数据,ad", &ad)
	//获取token
	//	var token models.Token
	//	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	//	if token.Token == "" {
	//		fmt.Println("token 为空")
	//		this.ajaxMsg("token is not nil", MSG_ERR_Param)
	//	}

	//	name, err := this.Token_auth(token.Token, "gree")
	//	if err != nil {
	//		fmt.Println("token err", err)
	//		this.ajaxMsg("token err!", MSG_ERR_Verified)
	//	}
	//	fmt.Println("当前访问用户为:", name)

	if ad.Mid == "" {
		fmt.Println("mid 不能为空")
		this.ajaxMsg("mid不能为空", MSG_ERR_Param)
	}

	aid, err := strconv.ParseInt(this.GetRandomString(6), 10, 64)
	if err != nil {
		fmt.Println("cid string to int64 err")
	}
	ad.Aid = aid
	//	if ad.PicUrl != "" {
	//		PicList := strings.Split(ad.PicUrl, ":")
	//		ad.PicUrl = beego.AppConfig.String("img_upload_url") + PicList[2]
	//	}
	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("新建失败")
	}
	ad.CreatTime = nowtime
	//state 0
	ad.State = 0
	//insert
	_, err1 := o.Insert(&ad)
	if err1 != nil {
		fmt.Println("inset err!", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}

	list["aid"] = ad.Aid
	this.ajaxList("添加广告成功", MSG_OK, 1, list)

}

//编辑广告
func (this *AdController) Edit() {
	fmt.Println("编辑广告")
	o := orm.NewOrm()
	var ad models.Ad
	//获取编辑的广告
	json.Unmarshal(this.Ctx.Input.RequestBody, &ad)
	fmt.Printf("ad_info:", &ad)

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

	//aid
	aid := ad.Aid
	if aid == 0 {
		fmt.Printf("aid 不能为空")
		this.ajaxMsg("aid 不能为空", MSG_ERR_Param)
	}
	//time
	nowtime, err := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Println("新建失败")
	}
	ad.CreatTime = nowtime
	//	if ad.PicUrl != "" {
	//		PicList := strings.Split(ad.PicUrl, ":")
	//		ad.PicUrl = beego.AppConfig.String("img_upload_url") + PicList[2]
	//	}
	//更新
	num, err := o.Update(&ad)
	if err != nil {
		fmt.Println("更新广告信息失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	fmt.Println("updata ad num", num)

	this.ajaxMsg("更新成功", MSG_OK)

}

//删除广告
func (this *AdController) Del() {
	fmt.Println("删除广告")
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

	//获取得到aid
	var ad_info models.Ad
	json.Unmarshal(this.Ctx.Input.RequestBody, &ad_info)
	aid := ad_info.Aid
	if aid == 0 {
		fmt.Println("获取aid失败")
		this.ajaxMsg("aid 不能为空", MSG_ERR_Param)
	}
	ad := new(models.Ad)
	num, err := o.QueryTable(ad).Filter("Aid", aid).Delete()
	if err != nil {
		fmt.Println("删除广告失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	if num == 0 {
		fmt.Println("删除广告失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	this.ajaxMsg("删除成功", MSG_OK)

}

//获取广告
func (this *AdController) GetData() {
	fmt.Println("获取广告内容")
	o := orm.NewOrm()
	ad := new(models.Ad)
	query := o.QueryTable(ad)
	var maps []orm.Params

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

	//index
	index, err := this.GetInt("index")
	if err != nil {
		fmt.Println("获取index下标错误")
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

	_, err1 := query.OrderBy("-CreatTime").Values(&maps, "Aid", "Mid", "PicUrl", "Url", "Title", "Category", "Model", "State", "Display")
	if err1 != nil {
		fmt.Println("获取广告失败")
		this.ajaxMsg("获取广告失败", MSG_ERR_Param)
	}
	this.ajaxList("获取广告消息成功", MSG_OK, count, maps)

}

//获取品类
func (this *AdController) GetProductStyle() {
	fmt.Println("获取产品品类")
	o := orm.NewOrm()
	product := new(models.Product)
	query := o.QueryTable(product)
	var maps, mid_maps []orm.Params
	//list:=make(map[string]interface{})
	var result []map[string]interface{}

	num, err := query.Distinct().Values(&maps, "Style")
	if err != nil {
		fmt.Println("获取产品品类失败")
		this.ajaxMsg("获取产品品类失败", MSG_ERR_Param)
	}
	for _, m := range maps {
		list := make(map[string]interface{})
		var childrenList []map[string]interface{}
		num, err := query.Filter("Style", m["Style"].(string)).Values(&mid_maps, "Mid")
		if err != nil {
			fmt.Println("获取产品mid失败")
			this.ajaxMsg("获取产品mid失败", MSG_ERR_Param)
		}
		fmt.Println("mid num", num)
		//fmt.Println("mid_maps", mid_maps)
		for _, n := range mid_maps {
			list1 := make(map[string]interface{})
			list1["value"] = n["Mid"].(string)
			list1["label"] = n["Mid"].(string)
			childrenList = append(childrenList, list1)
		}
		list["value"] = m["Style"].(string)
		list["label"] = m["Style"].(string)
		list["children"] = childrenList
		result = append(result, list)
	}
	//fmt.Println("product_maps", maps)
	this.ajaxList("获取产品品类成功", MSG_OK, num, result)

}

//获取mid
func (this *AdController) GetMid() {
	fmt.Println("获取产品mid")
	o := orm.NewOrm()
	product := new(models.Product)
	query := o.QueryTable(product)
	var maps orm.ParamsList

	p := this.GetString("product")
	if p == "" {
		fmt.Println("product 为空")
		this.ajaxMsg("类型不能为空", MSG_ERR_Param)
	}

	num, err := query.Filter("Style", p).ValuesFlat(&maps, "Mid")
	if err != nil {
		fmt.Println("获取产品mid失败")
		this.ajaxMsg("获取产品mid失败", MSG_ERR_Param)
	}
	//fmt.Println("product_maps", maps)
	if num == 0 {
		maps = orm.ParamsList{}
		this.ajaxList("没有该类型的mid", MSG_OK, num, maps)
	}
	this.ajaxList("获取产品mid成功", MSG_OK, num, maps)
}

//新品咨询推送
func (this *AdController) NewAdPush() {
	fmt.Println("点击新品推荐")
	//获取数据
	o := orm.NewOrm()
	//	list := make(map[string]interface{})
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
	//获取Aid
	var ad_info models.Ad
	json.Unmarshal(this.Ctx.Input.RequestBody, &ad_info)
	fmt.Println("获取广告数据,ad", &ad_info)

	if ad_info.State == 1 {
		fmt.Println("该广告已推送过,不进行推送")
		this.ajaxMsg("该广告已推送过,不进行推送", MSG_ERR_Param)
	}

	if ad_info.Aid == 0 {
		fmt.Println("aid 不能为空")
		this.ajaxMsg("aid 不能为空", MSG_ERR_Param)
	}

	aid := ad_info.Aid
	ad := new(models.Ad)
	err1 := o.QueryTable(ad).Filter("Aid", aid).One(&ad_info)
	if err1 != nil {
		fmt.Println("get ad info err!", err)
		this.ajaxMsg("推送失败", MSG_ERR_Param)
	}
	//push
	msg := make(map[string]interface{})
	extras := make(map[string]interface{})
	notice := make(map[string]interface{})
	data := make(map[string]interface{})
	data["aid"] = aid
	data["imageurl"] = ad_info.PicUrl
	data["url"] = ad_info.Url
	data["t"] = "adpush"
	data["time"] = time.Now().UTC().Format("2006-01-02 15:04:05")

	extras["ext"] = data
	extras["msg"] = ad_info.Content
	extras["title"] = ad_info.Title
	extras["type"] = 0

	notice["alert"] = ""
	notice["extras"] = extras
	notice["loc-args"] = []string{}
	notice["loc-key"] = ""
	notice["title"] = ""

	msg["notice"] = notice
	msg["tag"] = []string{"all"}
	fmt.Println("push msg", msg)
	byteData, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("msg err")
	}
	push_url := beego.AppConfig.String("push_url")

	appid := beego.AppConfig.String("appid")
	appkey := beego.AppConfig.String("appkey")

	req, err := http.NewRequest("POST", push_url, strings.NewReader(string(byteData)))
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("appid", appid)
	tm := time.Now().UTC().Format("2006-01-02 15:04:05")
	req.Header.Set("t", tm)
	vc := fmt.Sprintf("%s_%s_%s", appid, appkey, tm)
	req.Header.Set("vc", this.Md5(vc))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		//noinspection GoPlaceholderCount
		fmt.Println("http client sending an HTTP request error: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		//noinspection GoPlaceholderCount
		fmt.Println("the response's status code from an HTTP request is: %d", res.StatusCode)
	} else {
		d, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("push_svr returned: ", string(d))

		//update state
		_, err1 := o.QueryTable(ad).Filter("Aid", aid).Update(orm.Params{
			"State": 1,
		})
		if err1 != nil {
			fmt.Println("update ad push state err!", err1.Error())
		}
	}

	this.ajaxMsg("推送成功", MSG_OK)

}

//api web
//获取广告
func (this *AdController) GetAd() {
	fmt.Println("获取广告内容")
	o := orm.NewOrm()
	ad := new(models.Ad)
	query := o.QueryTable(ad).Filter("Display", 1)
	var maps []orm.Params

	num, err := query.OrderBy("-CreatTime").Values(&maps, "Aid", "Mid", "PicUrl", "Url")
	if err != nil {
		fmt.Println("获取广告失败")
		this.ajaxMsg("获取广告失败", MSG_ERR_Param)
	}
	this.ajaxList("获取广告消息成功", MSG_OK, num, maps)

}
