package bsonkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConvert(t *testing.T) {
	res := Convert(bson.M{
		"foo": "bar",
		"bar": bson.A{
			bson.M{
				"foo": "bar",
			},
		},
		"baz": bson.D{
			bson.E{Key: "foo", Value: bson.M{
				"foo": "bar",
			}},
		},
	})
	assert.Equal(t, &bson.D{
		bson.E{Key: "bar", Value: bson.A{
			bson.D{
				bson.E{Key: "foo", Value: "bar"},
			},
		}},
		bson.E{Key: "baz", Value: bson.D{
			bson.E{Key: "foo", Value: bson.D{
				bson.E{Key: "foo", Value: "bar"},
			}},
		}},
		bson.E{Key: "foo", Value: "bar"},
	}, res)

	assert.Panics(t, func() {
		Convert(bson.M{
			"foo": uint(1),
		})
	})
}

func TestConvertList(t *testing.T) {
	res := ConvertList([]bson.M{
		{
			"foo": "bar",
		},
		{
			"baz": bson.D{
				bson.E{Key: "foo", Value: bson.M{
					"foo": "bar",
				}},
			},
		},
	})
	assert.Equal(t, List{
		&bson.D{
			bson.E{Key: "foo", Value: "bar"},
		},
		&bson.D{
			bson.E{Key: "baz", Value: bson.D{
				bson.E{Key: "foo", Value: bson.D{
					bson.E{Key: "foo", Value: "bar"},
				}},
			}},
		},
	}, res)
}
