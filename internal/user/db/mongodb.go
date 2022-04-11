package db

import (
	"awesomeProject4/internal/apper"
	"awesomeProject4/internal/user"
	"awesomeProject4/pkg/logging"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("create user error: %v", err)
	}
	d.logger.Debug("convert insertedid to objektid")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("filed to convert objecktid to hex: %s", oid)
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("faliet to convert hex to objectid, hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			//TODO ErrEntityNotFound
			return u, fmt.Errorf("ErrEntityNotFound")
		}
		return u, fmt.Errorf("filed to find one user by id: %s, %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("filed to find one user by id: %s, %v", id, err)
	}
	return u, nil
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, apper.ErrNotFound
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("filed to read all documents from cursor, err: %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	//Принемаем user и обновляем пользователя по ID
	//Значит нужно у пользователя забрать ID
	//ID надо конвертировать из Hex в objectID
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to Object, ID = %s", user.ID)
	}
	//Создаём фильтр
	filter := bson.M{"_id": objectID}

	UserBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("falied to marshal user, error :%v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(UserBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("filed to unmarshal uer byte, error: %v", err)
	}
	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("filed to execute update user quare, error: %v", err)
	}

	if result.MatchedCount == 0 {
		return apper.ErrNotFound
	}

	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to Object, ID = %s", id)
	}
	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execution query, error: %v", err)
	}
	if result.DeletedCount == 0 {
		return apper.ErrNotFound
	}
	d.logger.Tracef("Delete %d documents", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
