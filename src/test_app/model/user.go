package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"template2/lib/db"
	"template2/lib/util"
	"template2/test_app/config"
	"template2/test_app/constant"
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
//DefaultSelector = bson.M{}
)

/* API used Graph QL */
func CreateUser(user *UserInfo) error {
	/*query := bson.M{
		"username": user.Username,
	}*/

	userInfo := UserInfo{
		UID:        primitive.NewObjectID(),
		Nickname:   user.Nickname,
		Username:   user.Username,
		Email:      user.Email,
		CreateTime: util.GetNowTimestamp(),
		UpdateTime: util.GetNowTimestamp(),
	}
	dbClient := config.MongoDBClient

	return DbTxInsertUser(dbClient, &userInfo)
	/*user, err := DbOpFindUser(dbClient, &query)
	if err != nil {
		insertedID, err := DbOpInsertUser(dbClient, &userInfo)
		log.Println("Inserted a single document: ", insertedID)
		return err
	}

	return constant.ErrAccountExist*/
}

func GetUser(username string) (user *UserInfo, err error) {
	query := bson.M{
		"username": username,
	}
	dbClient := config.MongoDBClient

	user, err = DbOpFindUser(context.Background(), dbClient, &query)

	if err != nil {
		return user, constant.ErrAccountNotExist
	}
	return user, err
}

/* Database Transaction */
func DbTxInsertUser(client *db.MongoDBClient, userInfo *UserInfo) error {
	query := bson.M{
		"username": userInfo.Username,
	}

	_, err := client.WithTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := DbOpFindUser(sessCtx, client, &query)
		if err != nil {
			insertedID, err := DbOpInsertUser(sessCtx, client, &userInfo)
			log.Println("Inserted a single document: ", insertedID)
			log.Println(err)
			log.Println("==========")
			//result = insertedID
			return insertedID, err
		}
		return nil, constant.ErrAccountExist
	})

	return err
}

/* Database Operation: Insert, Deletion, Update, Select */
func DbOpFindUser(ctx context.Context, client *db.MongoDBClient, filter interface{}) (user *UserInfo, err error) {
	user = &UserInfo{}
	err = client.FindOne(ctx, constant.TableUser, filter, user)
	return user, err
}

func DbOpInsertUser(ctx context.Context, client *db.MongoDBClient, userInfo interface{}) (uid interface{}, err error) {
	result, err := client.InsertOne(ctx, constant.TableUser, userInfo)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}
