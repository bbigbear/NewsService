package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	//	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
)

const (
	MSG_OK            = 200 //返回成功
	MSG_ERR_Param     = 400 //用户发出的请求，服务器没有进行新建或修改操作
	MSG_ERR_Verified  = 401 //用户没有权限，用户密码错误
	MSG_ERR_Authority = 403 //用户得到授权（与401错误相对），但是访问是被禁止的
	MSG_ERR_Resources = 404 //请求不存在
	MSG_ERR           = 500 //服务器发生错误，用户将无法判断发出的请求是否成功
)

var (
	key []byte = []byte("greeplusapp")
)

type BaseController struct {
	beego.Controller
}

//ajax返回
func (this *BaseController) ajaxMsg(msg interface{}, msgno int) {
	out := make(map[string]interface{})
	out["code"] = msgno
	out["message"] = msg
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

//ajax返回 列表
func (this *BaseController) ajaxList(msg interface{}, msgno int, count int64, data interface{}) {
	out := make(map[string]interface{})
	out["code"] = msgno
	out["message"] = msg
	out["count"] = count
	out["data"] = data
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

// 通过两重循环过滤重复元素
func (this *BaseController) RemoveRepBySlice(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 文件上传接口
func (this *BaseController) PutFile() {
	fmt.Println("文件上传")
	f, h, err := this.GetFile("file")
	fmt.Println("文件名称", h.Filename)
	fmt.Println("文件大小", h.Size)
	if err != nil {
		log.Fatal("getfile err ", err)
		this.ajaxMsg(h.Filename+"文件上传失败", MSG_ERR_Resources)
	}
	defer f.Close()
	//imgUrl:=beego.AppConfig.String("img_url")
	b := httplib.Post(beego.AppConfig.String("img_url"))
	b.PostFile("file", h.Filename)
	str, err := b.String()
	if err != nil {
		fmt.Println("post err", err)
	}
	fmt.Println("str", str)

	fid := gjson.Get(str, "fid")
	fmt.Println("fid", fid)

	picurl := gjson.Get(str, "publicUrl")
	fmt.Println("picurl", picurl)

	url := fmt.Sprintf("%s%s%s", picurl, "/", fid)
	fmt.Println("url", url)

	path := "static/upload/" + h.Filename
	this.SaveToFile("file", path) // 保存位置在 static/upload, 没有文件夹要先创建
	//上传
	c := httplib.Post(fmt.Sprintf("%s%s", "http://", url))
	c.Param("key", "log2")
	c.PostFile("file", path)
	str1, err := c.String()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("final file", str1)
	list := make(map[string]interface{})
	picurlList := strings.Split(picurl.String(), ":")
	list["url"] = fmt.Sprintf("%s%s%s", beego.AppConfig.String("img_upload_url")+picurlList[1], "/", fid)
	list["name"] = h.Filename
	list["size"] = h.Size
	this.ajaxList("文件上传成功", MSG_OK, 1, list)
}

//将时间化为秒
func (this *BaseController) GetSecs(ordertime string) int64 {
	var s int64
	t, err := time.ParseInLocation("2006-01-02 15:04:05", ordertime, time.Local)
	if err == nil {
		s = t.Unix()
		return s
	} else {
		return -1
	}
}

//获取相差时间
func (this *BaseController) GetMinuteDiffer(server_time, mqtime string) int64 {
	var minute int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", server_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", mqtime, time.Local)
	if err == nil {
		diff := t1.Unix() - t2.Unix()
		minute = diff / 60
		return minute
	} else {
		return -1
	}
}

//生成随机数
func (this *BaseController) randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

//随机字符
func (this *BaseController) GetRandomString(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

type Claims struct {
	Appid string `json:"appid"`
	// recommended having
	jwt.StandardClaims
}

//jwt
func (this *BaseController) Create_token(appid string, secret string) (string, int64) {
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	claims := Claims{
		appid,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    appid,
		},
	}

	// Create the token using your claims
	c_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := c_token.SignedString([]byte(secret))

	return signedToken, expireToken
}

func (this *BaseController) Token_auth(signedToken, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Username, claims.StandardClaims.ExpiresAt)
		//fmt.Println(reflect.TypeOf(claims.StandardClaims.ExpiresAt))
		//return claims.Appid, err
		return claims.Appid, err
	}
	return "", err
}

//md5
func (this *BaseController) Md5(msg string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(msg))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

//uuid
func (this *BaseController) Uuid() string {
	u1 := uuid.NewV4().String()
	//fmt.Println("uuid", u1)
	return u1
}

//appid&scrent
func (this *BaseController) AppKeyScrect() (string, string) {
	appKey := this.Uuid()
	appScrect := this.Md5("key" + appKey)
	return appKey, appScrect
}
