package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conn *mongo.Client

type MongoDBRepo interface {
	Get(req primitive.D) (bson.D, error)
	GetAll(req primitive.D, limit, offset int64) ([]bson.D, error)
	Create(req interface{}) (string, error)
	Close()
}

type MongoDB struct {
	ctx        context.Context
	conn       *mongo.Client
	database   string
	collection string
}

func NewMogoDB(ctx context.Context, database, collection string) MongoDBRepo {
	host := viper.GetString("mongoDB_host")
	port := viper.GetString("mongoDB_port")
	user := viper.GetString("mongoDB_username")
	pass := viper.GetString("mongoDB_password")
	conn, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host+":"+port))

	return &MongoDB{
		ctx:        ctx,
		conn:       conn,
		database:   database,
		collection: collection,
	}
}

func (m *MongoDB) Get(req primitive.D) (bson.D, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := m.conn.Ping(ctx, nil); err != nil {
		return nil, errorHandling("Get", err)
	}

	var result bson.D
	err := m.conn.Database(m.database).Collection(m.collection).FindOne(ctx, req).Decode(&result)
	if err != nil {
		return nil, errorHandling("Get", err)
	}

	return result, nil
}

func (m *MongoDB) GetAll(req primitive.D, limit, offset int64) ([]bson.D, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := m.conn.Ping(ctx, nil); err != nil {
		return nil, errorHandling("Get", err)
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset)

	cur, err := m.conn.Database(m.database).Collection(m.collection).Find(ctx, req, opts)
	defer cur.Close(ctx)
	if err != nil {
		return nil, errorHandling("GetAll", err)
	}

	var result []bson.D
	for cur.Next(ctx) {
		var res bson.D
		err := cur.Decode(&res)
		if err != nil {
			return nil, errorHandling("GetAll", err)
		}
		result = append(result, res)
	}
	if err := cur.Err(); err != nil {
		return nil, errorHandling("GetAll", err)
	}

	return result, nil
}

func (m *MongoDB) Create(req interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := m.conn.Ping(ctx, nil); err != nil {
		return "", errorHandling("Create", err)
	}

	res, err := conn.Database(m.database).Collection(m.collection).InsertOne(ctx, req)
	if err != nil {
		return "", errorHandling("Create", err)
	}

	oID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errorHandling("Create", errors.New("InsertedID convert failed"))
	}

	return oID.String(), nil
}

func (m *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := conn.Disconnect(ctx); err != nil {
		errorHandling("Close", err)
	}
}

func errorHandling(functionName string, err error) error {
	return fmt.Errorf("error occur in func:%s, err:%w", functionName, err)
}
