package models

type Answer struct {
	Id      int64
	Title   string
	Content string
	Keyword string
}

type Ask struct {
	Id      int64
	Content string
}

type Keyword struct {
	Id   int64
	Aid  int64
	Word string
}

type AnswerDel struct {
	Id int64 `json:"qid"`
}
