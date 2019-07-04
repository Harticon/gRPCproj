package order

import (
	"context"
	"encoding/base64"
	"errors"
	user "github.com/Harticon/gRPCproj/user/proto"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
	"golang.org/x/crypto/scrypt"
	"sync"
)

type RouteGuideServer struct {
	access     Accesser
	mu         sync.Mutex
	routeNotes *user.User
}

func (s *RouteGuideServer) SignUp(ctx context.Context, usr *user.User) (*user.User, error) {

	if usr.Email == "" || usr.Password == "" {
		return &user.User{}, errors.New("non valid user")
	}

	_, err := govalidator.ValidateStruct(usr)
	if err != nil {
		return &user.User{}, err
	}

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return &user.User{}, err
	}

	//_, err = s.access.CreateUser(*usr)
	//if err != nil {
	//	return &user.User{}, err
	//}

	return usr, nil
}

func (s *RouteGuideServer) SignIn(ctx context.Context, usr *user.User) (*user.Token, error) {

	//return login token

	_, err := govalidator.ValidateStruct(usr)
	if err != nil {
		return &user.Token{}, err
	}

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return &user.Token{}, err
	}

	//_, err = s.access.GetUser(usr.Email, usr.Password)
	//if err != nil {
	//	fmt.Println("error: ", err)
	//	return &user.Token{}, err
	//}

	return &user.Token{Token: "validtoken"}, err
}

func hashPassword(raw string) (string, error) {

	dk, err := scrypt.Key([]byte(raw), []byte(viper.GetString("hashSecret")), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}
