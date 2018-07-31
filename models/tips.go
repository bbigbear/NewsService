package models

import (
	"time"
)

type Tips struct {
	Tid        int64 `orm:"pk"`
	Source     string
	SourcePic  string
	Title      string
	Content    string
	Tag        string
	ArticlePic string
	ReadNum    int64
	Keyword    string
	Url        string
	CreatTime  time.Time
}
