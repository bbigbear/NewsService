package controllers

import (
	"NewsService/models"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.TplName = "qrcode.tpl"
}

func (this *MainController) Add() {
	fmt.Println("新建客服问题")
	o := orm.NewOrm()
	var answer models.Answer
	json.Unmarshal(this.Ctx.Input.RequestBody, &answer)
	fmt.Println("获取新建问题")

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

	//keyword
	if answer.Keyword == "" {
		fmt.Println("keyword为空")
		this.ajaxMsg("keyword 不能为空", MSG_ERR_Param)
	}
	//title
	if answer.Title == "" {
		fmt.Println("title为空")
		this.ajaxMsg("title 不能为空", MSG_ERR_Param)
	}
	//content
	if answer.Content == "" {
		fmt.Println("content为空")
		this.ajaxMsg("content 不能为空", MSG_ERR_Param)
	}

	id, err := strconv.ParseInt(this.GetRandomString(4), 10, 64)
	if err != nil {
		fmt.Println("cid string to int64 err")
	}
	answer.Id = id
	_, err1 := o.Insert(&answer)
	if err1 != nil {
		fmt.Println("inset answer err!")
		this.ajaxMsg("param err", MSG_ERR_Param)
	}

	//获取keyword
	wordlist := strings.Split(answer.Keyword, ",")
	l := len(wordlist)

	for i := 0; i < l; i++ {
		var key models.Keyword
		key.Aid = id
		key.Word = wordlist[i]
		_, err := o.Insert(&key)
		if err != nil {
			fmt.Println("insert key err!")
			this.ajaxMsg("param err", MSG_ERR_Param)
		}

	}

	list := make(map[string]interface{})
	list["qid"] = id
	this.ajaxList("插入成功", MSG_OK, 1, list)

}

func (this *MainController) Edit() {
	fmt.Println("编辑客服问题")
	o := orm.NewOrm()
	var answer models.Answer

	json.Unmarshal(this.Ctx.Input.RequestBody, &answer)
	fmt.Println("answer_info", answer)

	//获取问题id
	if answer.Id == 0 {
		fmt.Println("qid 为空")
		this.ajaxMsg("qid 不能为空", MSG_ERR_Param)
	}

	num, err := o.Update(&answer)
	if err != nil {
		fmt.Println("更新answer失败")
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	fmt.Println("updata answer num", num)

	//获取keyword 更新关键词
	wordlist := strings.Split(answer.Keyword, ",")
	l := len(wordlist)

	for i := 0; i < l; i++ {
		var key models.Keyword
		key.Aid = answer.Id
		key.Word = wordlist[i]

		if created, id, err := o.ReadOrCreate(&key, "Aid", "Word"); err == nil {
			if created {
				fmt.Println("new insert a object,id", id)
			} else {
				fmt.Println("get an object ,id", id)
			}
		}

	}

	this.ajaxMsg("更新成功", MSG_OK)
}

func (this *MainController) Del() {
	fmt.Println("删除问题")
	o := orm.NewOrm()
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

	//获取得到qid
	var answerdel models.AnswerDel
	json.Unmarshal(this.Ctx.Input.RequestBody, &answerdel)
	qid := answerdel.Id
	if qid == 0 {
		fmt.Println("获取qid失败")
		this.ajaxMsg("qid 不能为空", MSG_ERR_Param)
	}
	answer := new(models.Answer)
	num, err := o.QueryTable(answer).Filter("Id", qid).Delete()
	if err != nil {
		fmt.Println("删除问题失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	if num == 0 {
		fmt.Println("删除问题失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	//删除keyword
	keyword := new(models.Keyword)
	num1, err := o.QueryTable(keyword).Filter("Aid", qid).Delete()
	if err != nil {
		fmt.Println("删除keyword失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}
	if num1 == 0 {
		fmt.Println("删除keyword失败", err)
		this.ajaxMsg("请求参数有误", MSG_ERR_Param)
	}

	this.ajaxMsg("删除成功", MSG_OK)

}

func (this *MainController) GetAnswer() {
	fmt.Println("获取回复的内容")
	var answer_key models.Answer
	json.Unmarshal(this.Ctx.Input.RequestBody, &answer_key)
	keyword := answer_key.Keyword
	if keyword == "" {
		fmt.Println("null")
	}
	o := orm.NewOrm()
	//ask 收集问题
	var ask models.Ask
	ask.Content = keyword
	ask_num, err := o.Insert(&ask)
	if err != nil {
		fmt.Println("insert ask err!", err.Error())
	}
	fmt.Println("insert ask num", ask_num)

	//keyword
	var id int64
	//var kw_maps []orm.Params
	kw := new(models.Keyword)
	var kw_info models.Keyword
	query := o.QueryTable(kw)

	//answer
	as := new(models.Answer)
	var as_info []orm.Params
	cond := orm.NewCondition()

	exist := o.QueryTable(as).Filter("Title", keyword).Exist()
	if exist {
		num, err := o.QueryTable(as).Filter("Title", keyword).Values(&as_info)
		if err != nil {
			fmt.Println("get word err!")
		}
		this.ajaxList("ok", MSG_OK, num, as_info)
	} else {
		exist := query.Filter("Word", keyword).Exist()
		if exist {
			count, err := query.Filter("Word", keyword).Count()
			if err != nil {
				fmt.Println("get kw count err!")
			}
			if count == 1 {
				err1 := query.Filter("Word", keyword).One(&kw_info)
				if err1 == orm.ErrNoRows {
					fmt.Println("没有找到该keyword")
					this.ajaxMsg("找不到该keyword", MSG_ERR_Resources)
				}
				id = kw_info.Aid
				cond = cond.And("Id", id)
			} else {
				var id_list orm.ParamsList
				num, err := query.Filter("Word", keyword).ValuesFlat(&id_list, "Aid")
				if err != nil {
					fmt.Println("get id list err!")
				}
				fmt.Println("id list", id_list)
				fmt.Println("num", num)
				for i := 0; i < len(id_list); i++ {
					cond = cond.Or("Id", id_list[i])
				}
			}
		} else {
			this.ajaxMsg("找不到该keyword", MSG_ERR_Resources)
		}
	}

	//answer
	//	as := new(models.Answer)
	//	var as_info []orm.Params
	as_num, err2 := o.QueryTable(as).SetCond(cond).Values(&as_info)
	if err2 != nil {
		fmt.Println("没有找到该回复")
		this.ajaxMsg("找不到该回复", MSG_ERR_Resources)
	}
	fmt.Println("answer num:", as_num)
	this.ajaxList("ok", MSG_OK, as_num, as_info)
	return

}

//获取问题
func (this *MainController) GetData() {
	fmt.Println("获取问题内容")
	o := orm.NewOrm()
	answer := new(models.Answer)
	query := o.QueryTable(answer)

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

	_, err1 := query.Values(&maps)
	if err1 != nil {
		fmt.Println("获取问题失败")
		this.ajaxMsg("获取问题失败", MSG_ERR_Resources)
	}
	this.ajaxList("获取问题消息成功", MSG_OK, count, maps)

}
