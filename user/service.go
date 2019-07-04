package order

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	user "github.com/Harticon/gRPCproj/user/proto"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
	"golang.org/x/crypto/scrypt"
	"google.golang.org/grpc/status"
	"net/http"
)

type RouteGuideServer struct {
	access Accesser
}

func NewRouteGuideServer(a Accesser) *RouteGuideServer {
	return &RouteGuideServer{
		access: a,
	}
}

func (s *RouteGuideServer) SignUp(ctx context.Context, usr *user.User) (*user.User, error) {

	if usr.Email == "" || usr.Password == "" {
		return &user.User{}, status.Error(http.StatusBadRequest, errors.New("not valid user").Error())
	}

	_, err := govalidator.ValidateStruct(usr)
	if err != nil {
		return &user.User{}, status.Error(http.StatusBadRequest, err.Error())
	}

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return &user.User{}, status.Error(http.StatusBadRequest, err.Error())
	}

	_, err = s.access.CreateUser(*usr)
	if err != nil {
		return &user.User{}, err
	}

	return usr, nil
}

func (s *RouteGuideServer) SignIn(ctx context.Context, usr *user.User) (*user.Token, error) {

	//return login token

	_, err := govalidator.ValidateStruct(usr)
	if err != nil {
		return &user.Token{}, status.Error(http.StatusBadRequest, err.Error())
	}

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return &user.Token{}, status.Error(http.StatusBadRequest, err.Error())
	}

	query, err := s.access.GetUser(usr.Email, usr.Password)
	if err != nil {
		fmt.Println("error: couldnt find user", err)
		return &user.Token{}, err
	}

	fmt.Printf("user logged in %s with valid token\n", query.Email)
	return &user.Token{Token: "validtoken"}, nil
}

func hashPassword(raw string) (string, error) {

	dk, err := scrypt.Key([]byte(raw), []byte(viper.GetString("hashSecret")), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}
