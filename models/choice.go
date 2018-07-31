package models

import (
	"time"
)

//import (
//	"github.com/astaxie/beego/orm"
//)

type Choice struct {
	Cid        int64 `orm:"pk"`
	Source     string
	SourcePic  string
	Title      string
	Brief      string
	Content    string
	Tag        string
	ArticlePic string
	Keyword    string
	Url        string
	ReadNum    int64
	CreatTime  time.Time
}
