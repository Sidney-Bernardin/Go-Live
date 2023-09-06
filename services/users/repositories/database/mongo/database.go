package mongo

import (
	"bytes"
	"context"
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
	client *mongo.Client

	usersDB      *mongo.Database
	usersColl    *mongo.Collection
	sessionsColl *mongo.Collection
}

func NewDatabaseRepository(config *configuration.Config) (domain.DatabaseRepository, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.DBConnTimeout))
	defer cancel()

	// Connect to MongoDB.
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURL))
	if err != nil {
		return nil, errors.Wrap(err, "cannot cannot to mongodb")
	}

	// Ping MongoDB.
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "cannot ping mongodb")
	}

	usersDB := client.Database("users")

	return &databaseRepository{
		client:       client,
		usersDB:      usersDB,
		usersColl:    usersDB.Collection("users"),
		sessionsColl: usersDB.Collection("sessions"),
	}, nil
}

// InsertUser inserts user.
func (repo *databaseRepository) InsertUser(ctx context.Context, user *domain.User) (string, error) {
	user.MongoID = primitive.NewObjectID()
	_, err := repo.usersColl.InsertOne(ctx, user)
	return user.MongoID.Hex(), errors.Wrap(err, "cannot insert user")
}

// GetUser finds userID's User.
func (repo *databaseRepository) GetUser(ctx context.Context, userIDHex string, fields ...string) (*domain.User, error) {

	var (
		user *domain.User

		userID, _ = primitive.ObjectIDFromHex(userIDHex)
		opts      = options.FindOne().SetProjection(repo.bsonProjection(fields...))
	)

	// Find the User-ID's User.
	err := repo.usersColl.FindOne(ctx, bson.M{"_id": userID}, opts).Decode(&user)
	if err != nil {

		// Check if the User wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{Problem: domain.ProblemUserDoesntExist}
		}

		return nil, errors.Wrap(err, "cannot find user")
	}

	user.ID = user.MongoID.Hex()
	return user, nil
}

// GetUserByUsername finds username's User.
func (repo *databaseRepository) GetUserByUsername(ctx context.Context, username string, fields ...string) (*domain.User, error) {

	var (
		user *domain.User
		opts = options.FindOne().SetProjection(repo.bsonProjection(fields...))
	)

	// Find the username's User.
	err := repo.usersColl.FindOne(ctx, bson.M{"username": username}, opts).Decode(&user)
	if err != nil {

		// Check if the User wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{Problem: domain.ProblemUserDoesntExist}
		}

		return nil, errors.Wrap(err, "cannot find user")
	}

	user.ID = user.MongoID.Hex()
	return user, nil

}

// SearchUsers finds Users with usernames similar to the one given.
func (repo *databaseRepository) SearchUsers(ctx context.Context, username string, skip, limit int64, fields ...string) ([]*domain.User, error) {

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetProjection(repo.bsonProjection(fields...))

	// Find Users with usernames similar to the username.
	cursor, err := repo.usersColl.Find(ctx, bson.M{"username": bson.M{"$regex": username}}, opts)
	if err != nil {
		return nil, errors.Wrap(err, "cannot find users")
	}

	// Loop over and docode the cursor's documents.
	users := []*domain.User{}
	for cursor.Next(ctx) {

		// Decode the cursor's current document.
		var user *domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, errors.Wrap(err, "cannot decode user")
		}

		user.ID = user.MongoID.Hex()
		users = append(users, user)
	}

	return users, errors.Wrap(cursor.Err(), "cursor error")
}

// GetUserBySession finds sessionIDHex's User.
func (repo *databaseRepository) GetUserBySession(ctx context.Context, sessionIDHex string, projection ...string) (*domain.User, error) {

	var (
		session      *domain.Session
		sessionID, _ = primitive.ObjectIDFromHex(sessionIDHex)
	)

	// Find the Session-ID's Session.
	err := repo.sessionsColl.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&session)
	if err != nil {

		// Check if the Session wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{Problem: domain.ProblemSessionDoesntExist}
		}

		return nil, errors.Wrap(err, "cannot find session")
	}

	var (
		user *domain.User
		opts = options.FindOne().SetProjection(repo.bsonProjection(projection...))
	)

	// Find the User of the Sessions's User-ID.
	err = repo.usersColl.FindOne(ctx, bson.M{"_id": session.MongoUserID}, opts).Decode(&user)
	if err != nil {

		// Check if the User wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{Problem: domain.ProblemUserDoesntExist}
		}

		return nil, errors.Wrap(err, "cannot find user")
	}

	user.ID = user.MongoID.Hex()
	return user, nil
}

// UserExists checks if a User with one of the given key-value pairs exists.
func (repo *databaseRepository) UserExists(ctx context.Context, kvPairs ...string) (bool, string, error) {

	// Make sure each key has a value in the key-value pairs.
	if len(kvPairs)%2 != 0 {
		return false, "", errors.New("their must be exactly one value for each key")
	}

	opts := options.Count().SetLimit(1)

	// Loop over each key-value pair.
	for i := 0; i < len(kvPairs); i += 2 {

		// Count the Users that have the current key-value pair.
		count, err := repo.usersColl.CountDocuments(ctx, bson.M{kvPairs[i]: kvPairs[i+1]}, opts)
		if err != nil {
			return false, "", errors.Wrap(err, "cannot count users")
		}

		// Check if any Users were counted.
		if count > 0 {
			return true, kvPairs[i], nil
		}
	}

	return false, "", nil
}

// InsertSession inserts session.
func (repo *databaseRepository) InsertSession(ctx context.Context, session *domain.Session) (string, error) {

	session.MongoID = primitive.NewObjectID()
	session.MongoUserID, _ = primitive.ObjectIDFromHex(session.UserID)

	// Insert the Session.
	_, err := repo.sessionsColl.InsertOne(ctx, session)
	return session.MongoID.Hex(), errors.Wrap(err, "cannot insert session")
}

// DeleteSession deletes sessionIDHex's Session.
func (repo *databaseRepository) DeleteSession(ctx context.Context, sessionIDHex string) error {
	sessionID, _ := primitive.ObjectIDFromHex(sessionIDHex)
	_, err := repo.sessionsColl.DeleteOne(ctx, bson.M{"_id": sessionID})
	return errors.Wrap(err, "cannot delete session")
}

// UpdateProfilePicture writes profilePicture to the profile-picture file of userIDHex's User.
func (repo *databaseRepository) UpdateProfilePicture(ctx context.Context, userIDHex string, profilePicture []byte) error {

	var (
		user      *domain.User
		userID, _ = primitive.ObjectIDFromHex(userIDHex)
	)

	// Find the User-ID's User.
	err := repo.usersColl.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {

		// Check if the User wasn't found.
		if err == mongo.ErrNoDocuments {
			return domain.ProblemDetail{Problem: domain.ProblemUserDoesntExist}
		}

		return errors.Wrap(err, "cannot find user")
	}

	// Create a gridfs Bucket.
	bucket, err := gridfs.NewBucket(repo.usersDB)
	if err != nil {
		return errors.Wrap(err, "cannot create gridfs bucket")
	}

	// Check if a new ID needs to be generated for the profile picture's file.
	if user.MongoProfilePictureID == primitive.NilObjectID {

		user.MongoProfilePictureID = primitive.NewObjectID()

		// Update the User's profile-picture-ID to the new one.
		_, err = repo.usersColl.UpdateByID(ctx, userID, bson.M{
			"$set": bson.M{"profile_picture_id": user.MongoProfilePictureID},
		})

		if err != nil {
			return errors.Wrap(err, "cannot update user")
		}
	}

	// Delete the User's previous profile picture.
	if err := bucket.Delete(user.MongoProfilePictureID); err != nil {
		return errors.Wrap(err, "cannot create upload stream")
	}

	// Create a gridfs UploadStream with the User's profile-picture-ID.
	uploadStream, err := bucket.OpenUploadStreamWithID(user.MongoProfilePictureID, user.Username)
	if err != nil {
		return errors.Wrap(err, "cannot create upload stream")
	}
	defer uploadStream.Close()

	// Write the profile picture to the UploadStream.
	_, err = uploadStream.Write(profilePicture)
	return errors.Wrap(err, "cannot write to upload stream")
}

// GetProfilePicture gets the profile picture file of userIDHex's User.
func (repo *databaseRepository) GetProfilePicture(ctx context.Context, userIDHex string) (*bytes.Buffer, error) {

	var (
		user      *domain.User
		userID, _ = primitive.ObjectIDFromHex(userIDHex)
	)

	// Find the User-ID's User.
	err := repo.usersColl.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {

		// Check if the User wasn't found.
		if err == mongo.ErrNoDocuments {
			return nil, domain.ProblemDetail{Problem: domain.ProblemUserDoesntExist}
		}

		return nil, errors.Wrap(err, "cannot find user")
	}

	// Create a gridfs Bucket.
	bucket, err := gridfs.NewBucket(repo.usersDB)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create gridfs bucket")
	}

	// Download the profile-picture-ID's profile-picture.
	var buf bytes.Buffer
	if _, err = bucket.DownloadToStream(user.MongoProfilePictureID, &buf); err != nil {

		if err == gridfs.ErrFileNotFound {
			return nil, domain.ProblemDetail{Problem: domain.ProblemProfilePictureDoesntExist}
		}

		return &buf, errors.Wrap(err, "cannot get profile-picture")
	}

	return &buf, nil
}
