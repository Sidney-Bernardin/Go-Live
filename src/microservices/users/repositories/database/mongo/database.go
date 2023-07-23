package mongo

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"users/configuration"
	"users/domain"
)

type databaseRepository struct {
	config *configuration.Configuration
	client *mongo.Client
}

func NewDatabaseRepository(config *configuration.Configuration) (domain.DatabaseRepository, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.DatabaseTimeout))
	defer cancel()

	// Connect to MongoDB.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		return nil, errors.Wrap(err, "cannot cannot to mongodb")
	}

	// Create a TTL index for the sessions collection.
	_, err = client.Database("users").Collection("sessions").
		Indexes().
		CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.M{"expire_at": 1},
			Options: options.Index().SetExpireAfterSeconds(0),
		})

	if err != nil {
		return nil, errors.Wrap(err, "cannot index the sessions collection")
	}

	// Ping MongoDB.
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "cannot ping mongodb")
	}

	return &databaseRepository{config, client}, nil
}

func (repo *databaseRepository) InsertUserWithSession(user *domain.User, session *domain.Session) (string, error) {

	db := repo.client.Database("users")
	user.MongoID = primitive.NewObjectID()

	// Insert the user.
	if _, err := db.Collection("users").InsertOne(context.Background(), user); err != nil {
		return "", errors.Wrap(err, "cannot insert user")
	}

	session.MongoID = primitive.NewObjectID()
	session.MongoUserID = user.MongoID

	// Insert the session.
	_, err := db.Collection("sessions").InsertOne(context.Background(), session)
	return session.MongoID.Hex(), errors.Wrap(err, "cannot insert session")
}

func (repo *databaseRepository) GetUser(userIDHex string, fields ...string) (*domain.User, error) {

	// Convert the user-ID into a primitive.ObjectID.
	userID, err := repo.hexToObjID(userIDHex)
	if err != nil {
		return nil, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidID,
			Detail: "The given user-ID isn't valid.",
		}
	}

	// Find the user.
	res := repo.client.Database("users").Collection("users").FindOne(
		context.Background(),
		bson.M{"_id": userID},
		options.FindOne().SetProjection(repo.bsonProjection(fields...)),
	)

	// Decode the result into a user.
	var user *domain.User
	if err := res.Decode(&user); err != nil {

		// Check if the user wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{
				Type: domain.PDTypeUserDoesntExist,
			}
		}

		return nil, errors.Wrap(err, "cannot find user")
	}

	user.ID = user.MongoID.Hex()
	return user, nil
}

func (repo *databaseRepository) SearchUsers(username string, skip, limit int64, fields ...string) ([]*domain.User, error) {

	ctx := context.Background()

	// Find the users with regex.
	cursor, err := repo.client.Database("users").Collection("users").Find(
		ctx,
		bson.M{"username": bson.M{"$regex": username}},
		options.Find().
			SetSkip(skip).
			SetLimit(limit).
			SetProjection(repo.bsonProjection(fields...)))

	if err != nil {
		return nil, errors.Wrap(err, "cannot find users")
	}

	// Decode the results into users.
	users := []*domain.User{}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, errors.Wrap(err, "cannot decode users")
	}

	// Set each user's ID to it's MongoID.
	for _, u := range users {
		u.ID = u.MongoID.Hex()
	}

	return users, nil
}

func (repo *databaseRepository) GetUserByUsername(username string, fields ...string) (*domain.User, error) {

	// Find the user.
	res := repo.client.Database("users").Collection("users").FindOne(
		context.Background(),
		bson.M{"username": username},
		options.FindOne().SetProjection(repo.bsonProjection(fields...)),
	)

	// Decode the result into a user.
	var user *domain.User
	if err := res.Decode(&user); err != nil {

		// Check if the user wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{
				Type: domain.PDTypeUserDoesntExist,
			}
		}

		return nil, errors.Wrap(err, "cannot find user")
	}

	user.ID = user.MongoID.Hex()
	return user, nil
}

func (repo *databaseRepository) GetSessionsUser(sessionIDHex string, projection ...string) (*domain.User, error) {

	// Convert the session-ID into a primitive.ObjectID.
	sessionID, err := repo.hexToObjID(sessionIDHex)
	if err != nil {
		return nil, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidID,
			Detail: "The given session-ID isn't valid.",
		}
	}

	ctx := context.Background()
	db := repo.client.Database("users")

	// Find the session.
	res := db.Collection("sessions").FindOne(ctx, bson.M{"_id": sessionID})

	// Decode the result into a session.
	var session *domain.Session
	if err := res.Decode(&session); err != nil {

		// Check if the session wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{
				Type: domain.PDTypeSessionDoesntExist,
			}
		}

		return nil, errors.Wrap(err, "cannot find session")
	}

	// Find the user.
	res = db.Collection("users").FindOne(
		ctx,
		bson.M{"_id": session.MongoUserID},
		options.FindOne().SetProjection(repo.bsonProjection(projection...)),
	)

	// Decode the result into a user.
	var user *domain.User
	if err := res.Decode(&user); err != nil {

		// Check if the user wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{
				Type: domain.PDTypeUserDoesntExist,
			}
		}

		return nil, errors.Wrap(err, "cannot find user")
	}

	user.ID = user.MongoID.Hex()
	return user, nil
}

func (repo *databaseRepository) CheckForTakenUserFields(fields map[string]any) error {

	for key, value := range fields {

		// Find a user with a matching field.
		err := repo.client.Database("users").Collection("users").
			FindOne(context.Background(), bson.M{key: value}).Err()

		if err != nil {

			// Check if a user wasn't found.
			if err == mongo.ErrNoDocuments {
				continue
			}

			return errors.Wrap(err, "cannot find user")
		}

		return domain.ProblemDetail{
			Type:   domain.PDTypeFieldTaken,
			Detail: fmt.Sprintf("The %s has been taken by another user.", key),
		}
	}

	return nil
}

func (repo *databaseRepository) InsertSession(session *domain.Session) (string, error) {

	// Convert the session's user-ID into a primitive.ObjectID.
	userID, err := repo.hexToObjID(session.UserID)
	if err != nil {
		return "", domain.ProblemDetail{
			Type:   domain.PDTypeInvalidID,
			Detail: "The given user-ID isn't valid.",
		}
	}

	session.MongoID = primitive.NewObjectID()
	session.MongoUserID = userID

	// Insert the session.
	_, err = repo.client.Database("users").Collection("sessions").
		InsertOne(context.Background(), session)

	return session.MongoID.Hex(), errors.Wrap(err, "cannot insert session")
}

func (repo *databaseRepository) DeleteSession(sessionIDHex string) error {

	// Convert the session-ID into a primitive.ObjectID.
	sessionID, err := repo.hexToObjID(sessionIDHex)
	if err != nil {
		return domain.ProblemDetail{
			Type:   domain.PDTypeInvalidID,
			Detail: "The given session-ID isn't valid.",
		}
	}

	// Delete the session.
	_, err = repo.client.
		Database("users").
		Collection("sessions").
		DeleteOne(context.Background(), bson.M{"_id": sessionID})

	return errors.Wrap(err, "cannot delete session")
}

func (repo *databaseRepository) InsertProfilePicture(userIDHex string, profilePicture []byte) error {

	if !primitive.IsValidObjectID(userIDHex) {
		return domain.ProblemDetail{
			Type:   domain.PDTypeInvalidID,
			Detail: "The given user-ID isn't valid.",
		}
	}

	// Create a gridfs bucket.
	bucket, err := gridfs.NewBucket(repo.client.Database("users"))
	if err != nil {
		return errors.Wrap(err, "cannot create gridfs bucket")
	}

	// Delete any duplicates of the profile-picture.
	if err := repo.deleteFilesByName(bucket, userIDHex); err != nil {
		return errors.Wrap(err, "cannot delete files by name")
	}

	// Create an upload stream.
	uploadStream, err := bucket.OpenUploadStream(userIDHex)
	if err != nil {
		return errors.Wrap(err, "cannot create upload stream")
	}
	defer uploadStream.Close()

	// Write the profile-picture to the upload stream.
	_, err = uploadStream.Write(profilePicture)
	return errors.Wrap(err, "cannot write to upload stream")
}

func (repo *databaseRepository) GetProfilePicture(userIDHex string) (*bytes.Buffer, error) {

	if !primitive.IsValidObjectID(userIDHex) {
		return nil, domain.ProblemDetail{
			Type:   domain.PDTypeInvalidID,
			Detail: "The given user-ID isn't valid.",
		}
	}

	// Create a gridfs bucket.
	bucket, err := gridfs.NewBucket(repo.client.Database("users"))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create gridfs bucket")
	}

	// Download the profile-picture.
	var buf bytes.Buffer
	_, err = bucket.DownloadToStreamByName(userIDHex, &buf)

	// Check if the profile-picture wasn't found.
	if err == gridfs.ErrFileNotFound {
		return nil, domain.ProblemDetail{
			Type: domain.PDTypeProfilePictureDoesntExist,
		}
	}

	return &buf, errors.Wrap(err, "cannot get profile-picture")
}
