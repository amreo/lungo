package mongokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/256dpi/lungo/bsonkit"
)

func matchTest(t *testing.T, doc, query bson.M, result interface{}) {
	t.Run("Mongo", func(t *testing.T) {
		coll := testCollection()
		_, err := coll.InsertOne(nil, doc)
		assert.NoError(t, err)
		n, err := coll.CountDocuments(nil, query)
		if result == nil {
			assert.Error(t, err, query)
			assert.Zero(t, n, query)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, result, n == 1, query)
		}
	})

	t.Run("Lungo", func(t *testing.T) {
		res, err := Match(bsonkit.Convert(doc), bsonkit.Convert(query))
		if result == nil {
			assert.Error(t, err, query)
			assert.False(t, res, query)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, result, res, query)
		}
	})
}

func TestMatchBasic(t *testing.T) {
	// empty query filter
	matchTest(t, bson.M{
		"foo": "bar",
	}, bson.M{}, true)

	// empty top level and operator
	matchTest(t, bson.M{
		"foo": "bar",
	}, bson.M{
		"$and": bson.A{},
	}, nil)
}

func TestMatchEq(t *testing.T) {
	matchTest(t, bson.M{
		"foo": "bar",
	}, bson.M{
		"foo": "bar",
	}, true)

	matchTest(t, bson.M{
		"foo": "bar",
	}, bson.M{
		"foo": bson.M{
			"$eq": "bar",
		},
	}, true)

	matchTest(t, bson.M{
		"foo": "bar",
	}, bson.M{
		"foo": "baz",
	}, false)

	matchTest(t, bson.M{
		"foo": "bar",
	}, bson.M{
		"foo": bson.M{
			"$eq": "baz",
		},
	}, false)

	matchTest(t, bson.M{
		"foo": bson.M{
			"bar": bson.M{
				"$eq": "baz",
			},
		},
	}, bson.M{
		"foo": bson.M{
			"bar": bson.M{
				"$eq": "baz",
			},
		},
	}, true)

	matchTest(t, bson.M{
		"foo": bson.M{
			"bar": "baz",
		},
	}, bson.M{
		"foo": bson.M{
			"$eq": bson.M{
				"bar": "baz",
			},
		},
	}, true)
}