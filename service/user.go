// addUser
// checkUser
package service

import (
	"fmt"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model"

	"github.com/olivere/elastic/v7"
)

//1. Read ES based on username and compare given password
//2.

func CheckUser(username, password string) (bool, error) {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("username", username))
	query.Must(elastic.NewTermQuery("password", password))
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	var utype model.User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(model.User)
		if u.Password == password {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}
	}
	return false, nil
}

// bool: 用户名是否已经存在
func AddUser(user *model.User) (bool, error) {
	//从post里拿
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	if searchResult.TotalHits() > 0 {
		return false, nil
	}

	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}
