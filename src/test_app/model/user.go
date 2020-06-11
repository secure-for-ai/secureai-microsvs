package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

/* API used by Graph QL */
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

	return dbTxInsertUser(dbClient, &userInfo)
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

func ListUser(username string, page, perPage int64) (int64, *[]UserInfo, error) {
	query := bson.M{}

	if username != "" {
		query["username"] = bson.M{"$regex": username}
	}
	dbClient := config.MongoDBClient

	count, err := listUserCount(context.Background(), dbClient, &query)
	if err != nil {
		return 0, nil, constant.ErrDatabase
	}

	users, err := listUser(context.Background(), dbClient, &query, page, perPage)
	if err != nil {
		return 0, nil, constant.ErrDatabase
	}

	return count, users, err
}

/* Database Transaction */
func dbTxInsertUser(client *db.MongoDBClient, userInfo *UserInfo) error {
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

func listUserCount(ctx context.Context, client *db.MongoDBClient, filter interface{}) (count int64, err error) {
	count, err = client.GetTable(constant.TableUser).CountDocuments(ctx, filter)
	return count, err
}

func listUser(ctx context.Context, client *db.MongoDBClient, filter interface{}, page, perPage int64) (data *[]UserInfo, err error) {
	data = &[]UserInfo{}
	findOptions := options.Find()
	// Sort by `updateTime` field descending
	findOptions.SetSort(bson.D{{"updateTime", -1}})
	// Skip (page-1) pages
	findOptions.SetSkip((page - 1) * perPage)
	// only return perPage records
	findOptions.SetLimit(perPage)
	cursor, err := client.GetTable(constant.TableUser).Find(ctx, filter, findOptions)
	if err != nil {
		return data, err
	}
	err = cursor.All(ctx, data)

	return data, err
}
