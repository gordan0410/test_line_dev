package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"test_line_dev/domain"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	ctx        context.Context
	conn       *mongo.Client
	database   string
	collection string
}

func NewMogoDB(ctx context.Context, database, collection string) domain.MessageRepo {
	host := viper.GetString("mongoDB_host")
	port := viper.GetString("mongoDB_port")
	user := viper.GetString("mongoDB_username")
	pass := viper.GetString("mongoDB_password")
	conn, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host+":"+port))

	return &mongoDB{
		ctx:        ctx,
		conn:       conn,
		database:   database,
		collection: collection,
	}
}

func (m *mongoDB) Get(req domain.GetMessage) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := m.conn.Ping(ctx, nil); err != nil {
		return nil, errorHandling("Get", err)
	}

	var result primitive.M
	err := m.conn.Database(m.database).Collection(m.collection).FindOne(ctx, req).Decode(&result)
	if err != nil {
		return nil, errorHandling("Get", err)
	}

	return result, nil
}

func (m *mongoDB) GetAll(req domain.GetMessage, limit, offset int64) ([]map[string]interface{}, error) {
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

	var result []map[string]interface{}
	for cur.Next(ctx) {
		var res primitive.M
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

func (m *mongoDB) Create(req domain.CreateMessage) (string, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := m.conn.Ping(ctx, nil); err != nil {
		return "", errorHandling("Create", err)
	}

	res, err := m.conn.Database(m.database).Collection(m.collection).InsertOne(ctx, req)
	if err != nil {
		return "", errorHandling("Create", err)
	}

	oID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errorHandling("Create", errors.New("InsertedID convert failed"))
	}

	return oID.String(), nil
}

func (m *mongoDB) CountDocuments(req domain.GetMessage) (int, error) {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := m.conn.Ping(ctx, nil); err != nil {
		return 0, errorHandling("Create", err)
	}

	res, err := m.conn.Database(m.database).Collection(m.collection).CountDocuments(ctx, req)
	if err != nil {
		return 0, errorHandling("Create", err)
	}

	return int(res), nil
}

func (m *mongoDB) Close() {
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	if err := m.conn.Disconnect(ctx); err != nil {
		errorHandling("Close", err)
	}
}

func errorHandling(functionName string, err error) error {
	return fmt.Errorf("error occur in func:%s, err:%w", functionName, err)
}
