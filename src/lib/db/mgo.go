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

func (this MongoDBConf) GetMongoURI() string {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?", //retryWrites=false&replSet=rs0
		this.User, this.PW, this.Host, this.Port, this.DBName)
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
	//client.database.WriteConcern()

	return client, err
}

func (this *MongoDBClient) Copy(ref *MongoDBClient) {
	this.client = ref.client
	this.database = ref.database
}

func (this *MongoDBClient) Disconnect() {
	err := this.client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDBClient closed.")
}

func (this *MongoDBClient) GetClient() *mongo.Client {
	return this.client
}

func (this *MongoDBClient) UseDatabase(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	this.database = this.client.Database(name, opts...)
	return this.database
}

func (this *MongoDBClient) GetCurDatabase() *mongo.Database {
	return this.database
}

func (this *MongoDBClient) GetDatabaseByName(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return this.client.Database(name, opts...)
}

func (this *MongoDBClient) GetTable(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return this.database.Collection(name, opts...)
}

func (this *MongoDBClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return this.client.Ping(ctx, rp)
}

func (this *MongoDBClient) StartSession() (mongo.Session, error) {
	return this.client.StartSession()
}

func Test(sess mongo.Session, ctx context.Context) {
	log.Println("End Session")
	sess.EndSession(ctx)
}

func (this *MongoDBClient) WithTransaction(
	fn func(sessCtx mongo.SessionContext) (interface{}, error)) (result interface{}, err error) {
	var session mongo.Session
	var ctx = context.Background()

	if session, err = this.StartSession(); err != nil {
		return nil, err
	}

	if err = session.StartTransaction(); err != nil {
		return nil, err
	}

	// auto call AbortTransaction() if failed
	defer Test(session, ctx)
	//defer session.EndSession(ctx)

	if err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		result, err = fn(sessCtx)
		log.Println(result)
		if errCommit := session.CommitTransaction(sessCtx); errCommit != nil {
			result = nil
			log.Println("********")
			log.Println(errCommit)
			return errCommit
		}
		return err
	}); err != nil {
		return nil, err
	}
	log.Println("======********")
	log.Println(err)
	return result, err
}

func (this *MongoDBClient) InsertOne(ctx context.Context, tableName string, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	table := this.GetTable(tableName)
	return table.InsertOne(ctx, document, opts...)
}

func (this *MongoDBClient) UpdateOne(ctx context.Context, tableName string, filter, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	table := this.GetTable(tableName)
	return table.UpdateOne(ctx, filter, update, opts...)
}

func (this *MongoDBClient) FindOne(ctx context.Context, tableName string, filter, result interface{},
	opts ...*options.FindOneOptions) error {
	table := this.GetTable(tableName)
	return table.FindOne(ctx, filter, opts...).Decode(result)
}

func (this *MongoDBClient) DeleteOne(ctx context.Context, tableName string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	table := this.GetTable(tableName)
	return table.DeleteOne(ctx, filter, opts...)
}
