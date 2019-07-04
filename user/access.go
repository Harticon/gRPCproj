package order

import (
	"fmt"
	"github.com/Harticon/gRPCproj/user/proto"
	"github.com/jinzhu/gorm"
)

type Access struct {
	db *gorm.DB
}

func NewAccess(db *gorm.DB) *Access {
	return &Access{
		db: db,
	}
}

type Accesser interface {
	CreateUser(usr user.User) (user.User, error)
	GetUser(email, password string) (user.User, error)
}

func (a *Access) CreateUser(usr user.User) (user.User, error) {

	err := a.db.Create(&usr).Error
	if err != nil {
		fmt.Println(err)
		return user.User{}, err
	}

	return usr, nil
}

func (a *Access) GetUser(email, password string) (user.User, error) {
	var query user.User

	err := a.db.Where("email = ? AND password = ?", email, password).Find(&query).Error
	if err != nil {
		//fmt.Println("err:", err)
		return user.User{}, err
	}
	return query, nil
}
