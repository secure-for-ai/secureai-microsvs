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

	user, err = findUser(context.Background(), dbClient, &query)

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

	user, err = findUser(context.Background(), dbClient, &query)

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

	err := updateUser(context.Background(), dbClient, &query, &update)
	if err != nil {
		return err
	}

	return nil
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

	err := deleteUser(context.Background(), dbClient, &query)
	if err != nil {
		return err
	}
	return nil
}

/* Database Transaction */
func DbTxInsertUser(client *db.MongoDBClient, userInfo *UserInfo) error {
	query := bson.M{
		"username": userInfo.Username,
	}

	_, err := client.WithTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := findUser(sessCtx, client, &query)
		if err != nil {
			insertedID, err := insertUser(sessCtx, client, &userInfo)
			return insertedID, err
		}
		return nil, constant.ErrAccountExist
	})

	return err
}

/* Database Operation: Insert, Deletion, Update, Select */
func findUser(ctx context.Context, client *db.MongoDBClient, filter interface{}) (user *UserInfo, err error) {
	user = &UserInfo{}
	err = client.FindOne(ctx, constant.TableUser, filter, user)
	return user, err
}

func insertUser(ctx context.Context, client *db.MongoDBClient, userInfo interface{}) (interface{}, error) {
	result, err := client.InsertOne(ctx, constant.TableUser, userInfo)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, err
}

func updateUser(ctx context.Context, client *db.MongoDBClient,
	filter interface{}, userInfo interface{}) error {
	_, err := client.UpdateOne(ctx, constant.TableUser, filter, userInfo)
	if err != nil {
		return err
	}
	return err
}

func deleteUser(ctx context.Context, client *db.MongoDBClient, filter interface{}) error {
	_, err := client.DeleteOne(ctx, constant.TableUser, filter)
	if err != nil {
		return err
	}
	return err
}
