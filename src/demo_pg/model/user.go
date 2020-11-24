package model

import (
	"context"
	"encoding/gob"
	"strconv"
	"template2/demo_pg/config"
	"template2/demo_pg/constant"
	"template2/lib/util"
)

type UserInfo struct {
	UID        int64  `pg:"uid,omitempty" json:"uid,omitempty"`
	Username   string `pg:"username" json:"username"`      // username
	Nickname   string `pg:"nickname" json:"nickname"`      // nickname
	Email      string `pg:"email" json:"email"`            // email
	CreateTime int64  `pg:"create_time" json:"createTime"` // create time
	UpdateTime int64  `pg:"update_time" json:"updateTime"` // update time
}

// Helpers --------------------------------------------------------------------

func init() {
	gob.Register(UserInfo{})
}

/* API used by Graph QL */
func CreateUser(user *UserInfo) error {
	user.UID = config.SnowflakeNode.Generate().Int64()
	user.CreateTime = util.GetNowTimestamp()
	user.UpdateTime = user.CreateTime

	client := config.PGClient
	ctx := context.Background()
	conn, err := client.GetConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	count, err := conn.Insert(ctx,
		"INSERT INTO "+constant.TableUser+
			" (uid, username, nickname, email, create_time, update_time)"+
			" VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (username) DO NOTHING;",
		user.UID, user.Username, user.Nickname, user.Email, user.CreateTime, user.UpdateTime)
	if err != nil {
		return err
	}

	if count == 0 {
		return constant.ErrAccountExist
	}

	return nil
}

func GetUser(username string) (user *UserInfo, err error) {
	user = &UserInfo{}
	client := config.PGClient
	ctx := context.Background()
	conn, err := client.GetConn(ctx)

	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.FindOne(ctx, "SELECT * FROM "+constant.TableUser+" WHERE username=$1;", user, username)

	if err != nil {
		return nil, err
	}

	return user, err
}

func GetUserById(id string) (user *UserInfo, err error) {
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return user, constant.ErrParamIDFormatWrong
	}

	user = &UserInfo{}
	client := config.PGClient
	ctx := context.Background()
	conn, err := client.GetConn(ctx)

	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.FindOne(ctx, "SELECT * FROM "+constant.TableUser+" WHERE uid=$1;", user, uid)

	if err != nil {
		return nil, err
	}

	return user, err
}

func UpdateUser(user *UserInfo) error {
	client := config.PGClient
	ctx := context.Background()
	conn, err := client.GetConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Update(ctx,
		"UPDATE "+constant.TableUser+
			" SET username=$1, nickname=$2, email=$3, update_time=$4 WHERE uid=$5;",
		user.Username, user.Nickname, user.Email, user.UpdateTime, user.UID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id string) error {
	uid, ok := strconv.ParseInt(id, 10, 64)
	if ok != nil {
		return constant.ErrParamIDFormatWrong
	}

	client := config.PGClient
	ctx := context.Background()
	conn, err := client.GetConn(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Update(ctx, "DELETE FROM "+constant.TableUser+" WHERE uid=$1;", uid)
	if err != nil {
		return err
	}

	return nil
}

func ListUser(username string, page, perPage int64) (count int64, users *[]UserInfo, err error) {
	client := config.PGClient
	ctx := context.Background()
	tx, err := client.Begin(ctx)
	if err != nil {
		return 0, nil, constant.ErrDatabase
	}
	defer tx.RollBackDefer(ctx)
	users = &[]UserInfo{}
	err = tx.FindAll(ctx,
		"SELECT * FROM "+constant.TableUser+" WHERE username LIKE $1 ORDER BY uid ASC LIMIT $2 OFFSET $3;",
		users, "%"+username+"%", perPage, (page-1)*perPage)
	if err != nil {
		return 0, nil, err
	}

	count, err = tx.Count(ctx,
		"SELECT count(1) AS cc FROM "+constant.TableUser+" WHERE username LIKE $1 GROUP BY uid ORDER BY uid ASC;",
		"%"+username+"%")
	if err != nil {
		return 0, nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {

		return 0, nil, constant.ErrDatabase
	}

	return count, users, err
}

/* Database Transaction */
//func dbTxInsertUser(client *db.PGClient, userInfo *UserInfo) error {
//	query := bson.M{
//		"username": userInfo.Username,
//	}
//
//	_, err := client.WithTransaction(func(sessCtx mongo.SessionContext) (interface{}, error) {
//		_, err := findUser(sessCtx, client, &query)
//		if err != nil {
//			insertedID, err := insertUser(sessCtx, client, &userInfo)
//			return insertedID, err
//		}
//		return nil, constant.ErrAccountExist
//	})
//
//	return err
//}

/* Database Operation: Insert, Deletion, Update, Select */
//func findUser(ctx context.Context, client *db.PGClient, filter interface{}) (user *UserInfo, err error) {
//	user = &UserInfo{}
//	conn, err := client.GetConn(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Release()
//	err = conn.FindOne(ctx, "SELECT * FROM " + constant.TableUser + " WHERE username=$1;", user, filter)
//	if err != nil {
//		return nil, err
//	}
//	return user, err
//}

//func insertUser(ctx context.Context, client *db.PGClient, userInfo interface{}) (interface{}, error) {
//	result, err := client.InsertOne(ctx, constant.TableUser, userInfo)
//	if err != nil {
//		return nil, err
//	}
//	return result.InsertedID, err
//}

//func updateUser(ctx context.Context, client *db.PGClient,
//	filter interface{}, userInfo *UserInfo) error {
//	conn, err := client.GetConn(ctx)
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//
//	_, err = conn.Update(ctx,
//		"UPDATE user SET username=$1, nickname=$2, email=$3, update_time=$4 WHERE uid=$5;",
//		userInfo.Username, userInfo.Nickname, userInfo.Email, userInfo.UpdateTime, userInfo.UID)
//	if err != nil {
//		return err
//	}
//
//	return err
//}
//
//func deleteUser(ctx context.Context, client *db.PGClient, filter interface{}) error {
//	_, err := client.DeleteOne(ctx, constant.TableUser, filter)
//	if err != nil {
//		return err
//	}
//	return err
//}
//
//func listUserCount(ctx context.Context, client *db.PGClient, filter interface{}) (count int64, err error) {
//	count, err = client.GetTable(constant.TableUser).CountDocuments(ctx, filter)
//	return count, err
//}
//
//func listUser(ctx context.Context, client *db.PGClient, filter interface{}, page, perPage int64) (data *[]UserInfo, err error) {
//	data = &[]UserInfo{}
//	findOptions := options.Find()
//	// Sort by `updateTime` field descending
//	findOptions.SetSort(bson.D{{"updateTime", -1}})
//	// Skip (page-1) pages
//	findOptions.SetSkip((page - 1) * perPage)
//	// only return perPage records
//	findOptions.SetLimit(perPage)
//	cursor, err := client.GetTable(constant.TableUser).Find(ctx, filter, findOptions)
//	if err != nil {
//		return data, err
//	}
//	err = cursor.All(ctx, data)
//
//	return data, err
//}
