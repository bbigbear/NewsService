package main

import (
	"NewsService/controllers"
	"fmt"
)

type Data struct {
	controllers.BaseController
}

func main() {
	data := new(Data)
	appid, appscrect := data.AppKeyScrect()
	fmt.Println("appid", appid)
	fmt.Println("appscrent", appscrect)
}
