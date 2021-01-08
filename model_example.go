package main

import (
	"./collection"
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
	admin := "admin"
	idx := maleUsers.FindIndex(func(u interface{}) bool {
		return strings.Contains(u.(*UserModel).UserName, admin)
	})
	fmt.Println(idx)

	testData := []string{"1", "2", "34"}
	c := collection.ConvStrList2GoSlice(testData)
	index := c.FindIndex(func(s interface{}) bool {
		return s == "34"
	})
	fmt.Println(index, c[1:], len(c))

}
