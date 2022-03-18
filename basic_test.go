package gogo

import (
	"fmt"
	"github.com/wangleilei2010/gogo/collection"
	"strings"
	"testing"

	"github.com/wangleilei2010/gogo/db"
)

// UserModel @Table:`ci`.`t_user`
type UserModel struct {
	db.Table
	//TableName   string `db:"t_user"`
	UserName    string `db:"username"`
	DisplayName string `db:"display_name"`
}

type Person struct {
	Name        string
	ChineseName string
}

func TestBasicModelGenerics(t *testing.T) {
	var pool *db.ConnPool
	var err error
	pool, err = db.OpenPool("root:root@tcp(ip:3306)/ci", "basic_test.go")
	defer pool.Close()
	if err != nil {
		fmt.Println("open db err:", err)
		//return
	}

	if users, err := db.FetchAll[UserModel](pool, "display_name LIKE CONCAT(?,'%') AND sex=?", "李", 1); err != nil {
		fmt.Println(err)
	} else {
		users.Find(func(u UserModel) bool {
			return strings.Contains(u.DisplayName, "李")
		}).Foreach(func(u UserModel) {
			fmt.Println(u.DisplayName)
		})

		collection.Map[UserModel, Person](users, func(u UserModel) Person {
			return Person{
				Name:        u.UserName,
				ChineseName: u.DisplayName,
			}
		}).Foreach(func(p Person) {
			fmt.Println(p.ChineseName, p.Name, "##")
		})
	}

	if user, err := db.FetchOne[UserModel](pool, "sex=?", 1); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(user.UserName, user.IsDataNotNull())
		m := make(collection.GenericMap[string, UserModel])
		m[user.UserName] = user
		for k, v := range m {
			fmt.Println(k, v)
		}
	}

}
