package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"template2/test_app/config"
	"template2/test_app/constant"
	"time"
)

type UserInfo struct {
	UID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Username   string `bson:"username" json:"username"`     // username
	Nickname   string `bson:"nickname" json:"nickname"`     // nickname
	Email      string `bson:"email" json:"email"`           // email
	CreateTime int64  `bson:"createTime" json:"createTime"` // create time
	UpdateTime int64  `bson:"updateTime" json:"updateTime"` // update time
}

var (
	DefaultSelector = bson.M{}
)

/* API used Graph QL */
func GetUser(username string) (user UserInfo, err error) {
	query := bson.M{
		"username": username,
	}
	var userptr UserInfo
	userptr, err = dbTxFindUser(query, DefaultSelector)
	user = userptr

	if err != nil {
		return user, constant.ErrAccountNotExist
	}
	return user, err
}

func CreateUser(user *UserInfo) error {
	data := UserInfo{
		UID:        primitive.NewObjectID(),
		Nickname:   user.Nickname,
		Username:   user.Username,
		Email:      user.Email,
		CreateTime: time.Now().Unix(), //util.GetNowTimestamp(),
		UpdateTime: time.Now().Unix(), //util.GetNowTimestamp(),
	}
	log.Println(user.Nickname)
	dbClient := config.MongoDBClient
	table := dbClient.GetTable(constant.TableUser)

	insertResult, err := table.InsertOne(context.TODO(), data)

	log.Println("Inserted a single document: ", insertResult.InsertedID)
	//err := insertManagers(Manager)
	return err
}

/* Database Transaction */
func dbTxFindUser(query, selectField interface{}) (user UserInfo, err error) {
	//data := UserInfo{}
	err = nil

	//cntrl := db.NewCopyMgoDBCntlr()
	//defer cntrl.Close()

	user, err = dbOpFindUser(query, selectField)
	if err != nil {
		log.Println(err)
		return UserInfo{}, constant.ErrAccountNotExist
	}
	return user, err
}

/* Database Operation: Insert, Deletion, Update, Select */
func dbOpFindUser(query, selectField interface{}) (user UserInfo, err error) {
	user = UserInfo{}

	dbClient := config.MongoDBClient
	log.Println(dbClient)
	//dbClient, err := db.NewMongoDB()
	//log.Println(dbClient)
	//dbClient.UseDatabase("gtest")
	//
	table := dbClient.GetTable(constant.TableUser)

	/*ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//ctx := context.TODO()
	defer cancel()
	clientOption := options.Client().ApplyURI("mongodb://test:password@localhost:27017/admin")
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		log.Fatal(err)
	}*/

	/*err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}*/

	//db := client.Database("gtest")
	//table := db.Collection(constant.TableUser)

	yy := bson.M{
		"username": "test",
	}
	err = table.FindOne(context.TODO(), yy).Decode(&user)
	log.Print(user)
	log.Println("lll")

	return user, err
}
