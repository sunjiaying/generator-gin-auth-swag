package utils

import (
	"errors"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/store"

	"gopkg.in/mgo.v2/bson"
)

var TokenManager *manage.Manager

func GenerateToken(clientID string, clientSecret string) (ti oauth2.TokenInfo, err error) {
	theUser, err := GetUserTokenSecret(clientID, clientSecret)
	if (err != nil || len(theUser.EmployeeNo) == 0) && (clientID != "00001" && clientID != "99999") {
		err = errors.New("Account does not exist")
		return
	}
	//生成新的令牌前，作废此前的令牌 todo
	gt := oauth2.GrantType("client_credentials")
	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        "admin",
	}

	clientStore := store.NewClientStore()

	clientStore.Set(clientID, &models.Client{
		ID:     clientID,
		Secret: clientSecret,
	})

	TokenManager.MapClientStorage(clientStore)
	ti, err = TokenManager.GenerateAccessToken(gt, tgr)
	// token = ti.GetAccess()

	return
}

type TokenUser struct {
	Id         bson.ObjectId `bson:"_id"`
	EmployeeNo string        `bson:"ID"`
}

func GetUserTokenSecret(clientID string, clientSecret string) (theOne *TokenUser, err error) {
	// session, err := mgo.Dial(MONGODB)
	// defer session.Close()
	// session.SetMode(mgo.Monotonic, true)
	// c := session.DB("HR").C("EmployeeInfo")
	// query := c.Find(bson.M{"ID":clientID,"Info.ePhone": clientSecret, "eWorkState": "1"})
	// err = query.One(&theOne)

	// 请在这里完善获取令牌的身份验证，可参与以上注释代码
	theOne = new(TokenUser)
	theOne.EmployeeNo = "00209"

	return
}
