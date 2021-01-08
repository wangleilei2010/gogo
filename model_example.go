package main

import (
	"./db"
	"fmt"
	"strings"
)

type UserModel struct {
	db.Table
	UserName    string `db:"username"`
	DisplayName string `db:"display_name"`
}

func main() {
	db.SetGlobalConnStr("user:pwd@tcp(localhost:3306)/ci")
	userModel := db.NewDBModel(&UserModel{}, "t_user")
	maleUsers := userModel.FetchAll("sex=1")

	//for _, u := range maleUsers {
	//	user := u.(*UserModel)
	//	fmt.Println(user.DisplayName, user.UserName)
	//}

	maleUsers.ForEach(func(u interface{}) {
		user := u.(*UserModel)
		fmt.Println(user.DisplayName, user.UserName)
	})

	displayNames := maleUsers.MapToStrList(func(u interface{}) string {
		return u.(*UserModel).DisplayName
	})
	fmt.Println(displayNames)

	maleUsers.Filter(func(u interface{}) bool {
		return strings.Contains(u.(*UserModel).UserName, "wang")
	}).ForEach(func(u interface{}) {
		user := u.(*UserModel)
		fmt.Println(user.DisplayName, user.UserName)
	})
}
