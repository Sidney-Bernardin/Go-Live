package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// hexToObjID converts the given hex to a primitive.ObjectID.
func (*databaseRepository) hexToObjID(hex string) (primitive.ObjectID, error) {

	// Check if the hex is valid.
	if !primitive.IsValidObjectID(hex) {
		return primitive.NilObjectID, errors.New("hex is invalid")
	}

	// Convert hex into a primitive.ObjectID.
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "cannot create object id from a hex")
	}

	return id, nil
}

// bsonProjection converts the given projection into BSON.
func (*databaseRepository) bsonProjection(projection ...string) bson.M {
	ret := bson.M{"_id": 1}
	for _, v := range projection {
		if v != "" {
			ret[v] = 1
		}
	}
	return ret
}

// deleteFilesByName deletes all chunks and metadata associated with files that
// have the given filename.
func (*databaseRepository) deleteFilesByName(bucket *gridfs.Bucket, filename string) error {

	ctx := context.Background()

	// Find the files.
	cursor, err := bucket.GetFilesCollection().Find(ctx, bson.M{"filename": filename})
	if err != nil {
		return errors.Wrap(err, "cannot find file")
	}

	for cursor.Next(ctx) {

		// Decode the cursor's current file.
		var file map[string]any
		if err := cursor.Decode(&file); err != nil {
			return errors.Wrap(err, "cannot decode file")
		}

		// Delete the file.
		if err := bucket.Delete(file["_id"]); err != nil {
			return errors.Wrap(err, "cannot delete file")
		}
	}

	return nil
}
