package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type MongoDBConf struct {
	Host        string `json:"Host"`
	Port        string `json:"Port"`
	DBName      string `json:"DBName"`
	User        string `json:"User"`
	PW          string `json:"PW"`
	AdminDBName string `json:"AdminDBName"`
}

func (c MongoDBConf) GetMongoURI() string {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?retryWrites=false&replSet=rs0&connect=direct",
		c.User, c.PW, c.Host, c.Port, c.AdminDBName)
	return mongoURI
}

type MongoDBClient struct {
	client   *mongo.Client
	database *mongo.Database
}

/*type MongoDBSession struct {
	session *mongo.Session
}*/

func NewMongoDB(conf MongoDBConf) (client *MongoDBClient, err error) {
	client = &MongoDBClient{}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	clientOption := options.Client().ApplyURI(conf.GetMongoURI()) //"mongodb://test:password@127.0.0.1:27017/admin")
	client.client, err = mongo.Connect(ctx, clientOption)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println("Connected to MongoDB!")

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println("Ping MongoDB Success!")

	client.UseDatabase(conf.DBName)
	log.Printf("Use Database: \"%s\"\n", conf.DBName)

	return client, err
}

func (c *MongoDBClient) Copy(ref *MongoDBClient) {
	c.client = ref.client
	c.database = ref.database
}

func (c *MongoDBClient) Disconnect() {
	err := c.client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDBClient closed.")
}

func (c *MongoDBClient) GetClient() *mongo.Client {
	return c.client
}

func (c *MongoDBClient) UseDatabase(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	c.database = c.client.Database(name, opts...)
	return c.database
}

func (c *MongoDBClient) GetCurDatabase() *mongo.Database {
	return c.database
}

func (c *MongoDBClient) GetDatabaseByName(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return c.client.Database(name, opts...)
}

func (c *MongoDBClient) GetTable(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return c.database.Collection(name, opts...)
}

func (c *MongoDBClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return c.client.Ping(ctx, rp)
}

func (c *MongoDBClient) StartSession() (mongo.Session, error) {
	return c.client.StartSession()
}

func (c *MongoDBClient) WithTransaction(
	fn func(sessCtx mongo.SessionContext) (interface{}, error)) (result interface{}, err error) {
	var session mongo.Session
	var ctx = context.Background()

	if session, err = c.StartSession(); err != nil {
		return nil, err
	}

	if err = session.StartTransaction(); err != nil {
		return nil, err
	}

	// auto call AbortTransaction() if failed
	defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		result, err = fn(sessCtx)
		if errCommit := session.CommitTransaction(sessCtx); errCommit != nil {
			result = nil
			return errCommit
		}
		return err
	}); err != nil {
		return nil, err
	}

	return result, err
}

func (c *MongoDBClient) InsertOne(ctx context.Context, tableName string, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	table := c.GetTable(tableName)
	return table.InsertOne(ctx, document, opts...)
}

func (c *MongoDBClient) UpdateOne(ctx context.Context, tableName string, filter, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	table := c.GetTable(tableName)
	return table.UpdateOne(ctx, filter, update, opts...)
}

func (c *MongoDBClient) FindOne(ctx context.Context, tableName string, filter, result interface{},
	opts ...*options.FindOneOptions) error {
	table := c.GetTable(tableName)
	return table.FindOne(ctx, filter, opts...).Decode(result)
}

func (c *MongoDBClient) DeleteOne(ctx context.Context, tableName string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	table := c.GetTable(tableName)
	return table.DeleteOne(ctx, filter, opts...)
}
