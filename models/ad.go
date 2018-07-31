package models

import (
	"time"
)

type Ad struct {
	Aid       int64 `orm:"pk"`
	Title     string
	Category  string
	Model     string
	PicUrl    string
	Url       string
	Mid       string
	Content   string
	CreatTime time.Time
	Display   int
	Type      int
	State     int
}
