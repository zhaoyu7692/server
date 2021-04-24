package service

import (
	"github.com/satori/go.uuid"
	"main/model"
	"main/mysql"
	"sync"
	"time"
)

type AuthStatus int

const (
	Authority AuthStatus = iota
	SessionOverdue
	UnAuthority
)

//type AuthorizeServiceUser struct {
//	model.User
//	//model.Timestamp
//}

func AuthLogin(Uid int64, token string) AuthStatus {
	var user model.User
	if err := mysql.DBConn.Get(&user, "SELECT TOKEN, LAST_MODIFIED FROM user WHERE UID = ?", Uid); err != nil {
		return UnAuthority
	}
	if time.Now().After(user.LastModified.Add(6 * time.Hour)) {
		return SessionOverdue
	}
	if token != user.Token {
		return UnAuthority
	}
	if _, err := mysql.DBConn.Exec("UPDATE user SET LAST_MODIFIED = ? WHERE UID = ?", time.Now(), Uid); err != nil {
		return UnAuthority
	}
	return Authority
}

type LoginStatus int

const (
	LoginFail LoginStatus = iota
	LoginSuccess
	LoginWrongPassword
	LoginUnRegister
)

func Login(username string, password string) (*model.User, LoginStatus) {
	var user model.User
	if err := mysql.DBConn.Get(&user, "SELECT UID, USERNAME, PASSWORD, TOKEN FROM user WHERE USERNAME = ?", username); err != nil {
		return nil, LoginUnRegister
	}
	if password != user.Password {
		return nil, LoginWrongPassword
	}
	user.Token = uuid.NewV4().String()
	if _, err := mysql.DBConn.Exec("UPDATE user SET TOKEN = ?, LAST_MODIFIED = ? WHERE UID = ?", user.Token, time.Now(), user.Uid); err != nil {
		return nil, LoginFail
	}
	user.Password = ""
	return &user, LoginSuccess
}

type RegisterStatus int

const (
	RegisterFail RegisterStatus = iota
	RegisterSuccess
	RegisterRepetitiveUsername
)

var RegisterMutex sync.Mutex

func Register(username string, password string) (*model.User, RegisterStatus) {
	RegisterMutex.Lock()
	defer RegisterMutex.Unlock()
	var usernameCount int64
	if err := mysql.DBConn.Get(&usernameCount, "SELECT COUNT(*) FROM user where USERNAME = ?", username); err != nil {
		return nil, RegisterFail
	}
	if usernameCount > 0 {
		return nil, RegisterRepetitiveUsername
	}
	user := model.User{
		Username: username,
		Token:    uuid.NewV4().String(),
	}
	if result, err := mysql.DBConn.Exec("INSERT INTO user (USERNAME, PASSWORD, TOKEN) VALUES (?,?,?,?,?)", user.Username, password, user.Token); err == nil {
		if user.Uid, err = result.LastInsertId(); err != nil {
			if err := mysql.DBConn.Get(&user.Uid, "SELECT UID FROM user WHERE USERNAME = ?", user.Username); err == nil {
				return &user, RegisterSuccess
			}
		} else {
			return &user, RegisterSuccess
		}
	}
	return nil, RegisterFail
}
