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

type MongoDBClient struct {
	client   *mongo.Client
	database *mongo.Database
}

type MongoDBSession struct {
	session *mongo.Session
}

func NewMongoDB() (client *MongoDBClient, err error) {
	client = &MongoDBClient{}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	clientOption := options.Client().ApplyURI("mongodb://test:password@127.0.0.1:27017/admin")
	client.client, err = mongo.Connect(ctx, clientOption)

	//defer client.Disconnect()
	if err != nil {
		log.Fatal(err)
		return client, err
	}
	log.Println("hello db")
	return client, err

	/* sess := globalSess.Clone()
	return &MgoDBCntlr{
		sess: sess,
		client:   sess.DB(DBNAME),
	}*/
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

func (this *MongoDBClient) UseDatabase(name string, opts ...*options.DatabaseOptions) {
	this.database = this.client.Database(name, opts...)
}

func (this *MongoDBClient) GetTable(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return this.database.Collection(name, opts...)
}

func (this *MongoDBClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return this.client.Ping(ctx, rp)
}

func (this *MongoDBClient) StartSession() (*MongoDBSession, error) {
	sess, err := this.client.StartSession()

	if err != nil {
		return nil, err
	}

	mgsess := MongoDBSession{}
	mgsess.session = &sess

	return &mgsess, err
}
