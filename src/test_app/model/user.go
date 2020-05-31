package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

/* API used Graph QL */
func CreateUser(user *UserInfo) error {
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

func GetUserById(id string) (user *UserInfo, err error) {
	var uid primitive.ObjectID
	uid, ok := primitive.ObjectIDFromHex(id)
	if ok != nil {
		return user, constant.ErrParamIDFormatWrong
	}
	query := bson.M{
		"_id": uid,
	}
	dbClient := config.MongoDBClient

	user, err = DbOpFindUser(context.Background(), dbClient, &query)

	if err != nil {
		return user, constant.ErrAccountNotExist
	}

	return user, err
}

func UpdateUser(user *UserInfo) error {
	query := bson.M{
		"_id": user.UID,
	}
	update := bson.M{
		"$set": user,
	}
	dbClient := config.MongoDBClient

	success, err := DbOpUpdateUser(context.Background(), dbClient, &query, &update)

	if err != nil {
		return constant.ErrDatabase
	}
	if success {
		return nil
	}

	return constant.ErrAccountNotExist
}

func DeleteUser(id string) error {

	var uid primitive.ObjectID
	uid, ok := primitive.ObjectIDFromHex(id)
	if ok != nil {
		return constant.ErrParamIDFormatWrong
	}

	query := bson.M{
		"_id": uid,
	}
	dbClient := config.MongoDBClient

	success, err := DbOpDeleteUser(context.Background(), dbClient, &query)

	if err != nil {
		return constant.ErrDatabase
	}
	if success {
		return nil
	}

	return constant.ErrAccountNotExist
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

func DbOpUpdateUser(ctx context.Context, client *db.MongoDBClient,
	filter interface{}, userInfo interface{}) (success bool, err error) {
	result, err := client.UpdateOne(ctx, constant.TableUser, filter, userInfo)

	if err != nil {
		return false, err
	}

	return result.MatchedCount > 0, err
}

func DbOpDeleteUser(ctx context.Context, client *db.MongoDBClient, filter interface{}) (success bool, err error) {
	result, err := client.DeleteOne(ctx, constant.TableUser, filter)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, err
}
