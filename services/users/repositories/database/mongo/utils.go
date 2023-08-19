package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

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
