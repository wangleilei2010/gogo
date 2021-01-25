package gogo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wangleilei2010/gogo/collection"

	"github.com/wangleilei2010/gogo/db"
)

type UserModel struct {
	db.Table
	UserName    string `db:"username"`
	DisplayName string `db:"display_name"`
}

type Person struct {
	Name string
	Age  int
}

func TestBasicModel(t *testing.T) {
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

}

func TestBasicCollection(t *testing.T) {
	testData := []string{"1", "2", "34"}
	c := collection.NewGoSlice(testData)
	index := c.FindIndex(func(s interface{}) bool {
		return s == "34"
	})
	fmt.Println(index, c[1:], len(c))

	testStrs := collection.GoSlice{"1", "2", "3"}
	testStrs.ForEach(func(p interface{}) {
		fmt.Println(p.(string) + "0")
	})

	var people = []Person{
		{"leilei", 18},
		{"doudou", 4},
	}

	personSlice := collection.NewGoSlice(people)
	personSlice.ForEach(func(p interface{}) {
		fmt.Println(p.(Person).Name)
	})
	personSlice.Map(func(p interface{}) interface{} {
		return Person{"wang" + p.(Person).Name, p.(Person).Age + 1}
	}).ForEach(func(p interface{}) {
		fmt.Println(p.(Person).Name, p.(Person).Age)
	})

}
