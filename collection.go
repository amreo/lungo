package lungo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/256dpi/lungo/bsonkit"
)

var _ ICollection = &Collection{}

type Collection struct {
	ns     string
	name   string
	db     *Database
	client *Client
}

func (c *Collection) Aggregate(context.Context, interface{}, ...*options.AggregateOptions) (ICursor, error) {
	panic("not implemented")
}

func (c *Collection) BulkWrite(context.Context, []mongo.WriteModel, ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	panic("not implemented")
}

func (c *Collection) Clone(...*options.CollectionOptions) (ICollection, error) {
	panic("not implemented")
}

func (c *Collection) CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error) {
	panic("not implemented")
}

func (c *Collection) Database() IDatabase {
	return c.db
}

func (c *Collection) DeleteMany(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	panic("not implemented")
}

func (c *Collection) DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	panic("not implemented")
}

func (c *Collection) Distinct(context.Context, string, interface{}, ...*options.DistinctOptions) ([]interface{}, error) {
	panic("not implemented")
}

func (c *Collection) Drop(context.Context) error {
	panic("not implemented")
}

func (c *Collection) EstimatedDocumentCount(context.Context, ...*options.EstimatedDocumentCountOptions) (int64, error) {
	panic("not implemented")
}

func (c *Collection) Find(ctx context.Context, query interface{}, opts ...*options.FindOptions) (ICursor, error) {
	// merge options
	opt := options.MergeFindOptions(opts...)

	// assert unsupported options
	c.client.assertUnsupported(opt.AllowPartialResults == nil, "FindOptions.AllowPartialResults")
	c.client.assertUnsupported(opt.BatchSize == nil, "FindOptions.BatchSize")
	c.client.assertUnsupported(opt.Collation == nil, "FindOptions.Collation")
	c.client.assertUnsupported(opt.Comment == nil, "FindOptions.Comment")
	c.client.assertUnsupported(opt.CursorType == nil, "FindOptions.CursorType")
	c.client.assertUnsupported(opt.Hint == nil, "FindOptions.Hint")
	c.client.assertUnsupported(opt.Limit == nil, "FindOptions.Limit")
	c.client.assertUnsupported(opt.Max == nil, "FindOptions.Max")
	c.client.assertUnsupported(opt.MaxAwaitTime == nil, "FindOptions.MaxAwaitTime")
	c.client.assertUnsupported(opt.MaxTime == nil, "FindOptions.MaxTime")
	c.client.assertUnsupported(opt.Min == nil, "FindOptions.Min")
	c.client.assertUnsupported(opt.NoCursorTimeout == nil, "FindOptions.NoCursorTimeout")
	c.client.assertUnsupported(opt.OplogReplay == nil, "FindOptions.OplogReplay")
	c.client.assertUnsupported(opt.Projection == nil, "FindOptions.Projection")
	c.client.assertUnsupported(opt.ReturnKey == nil, "FindOptions.ReturnKey")
	c.client.assertUnsupported(opt.ShowRecordID == nil, "FindOptions.ShowRecordID")
	c.client.assertUnsupported(opt.Skip == nil, "FindOptions.Skip")
	c.client.assertUnsupported(opt.Snapshot == nil, "FindOptions.Snapshot")
	c.client.assertUnsupported(opt.Sort == nil, "FindOptions.Sort")

	// transform query
	qry, err := bsonkit.Transform(query)
	if err != nil {
		return nil, err
	}

	// TODO: Check supported operators.

	// get cursor
	csr, err := c.client.backend.find(c.ns, qry)
	if err != nil {
		return nil, err
	}

	return csr, nil
}

func (c *Collection) FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult {
	panic("not implemented")
}

func (c *Collection) FindOneAndDelete(context.Context, interface{}, ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	panic("not implemented")
}

func (c *Collection) FindOneAndReplace(context.Context, interface{}, interface{}, ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	panic("not implemented")
}

func (c *Collection) FindOneAndUpdate(context.Context, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	panic("not implemented")
}

func (c *Collection) Indexes() mongo.IndexView {
	panic("not implemented")
}

func (c *Collection) InsertMany(context.Context, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	panic("not implemented")
}

func (c *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	// merge options
	opt := options.MergeInsertOneOptions(opts...)

	// assert unsupported options
	c.client.assertUnsupported(opt.BypassDocumentValidation == nil, "InsertOneOptions.BypassDocumentValidation")

	// transform document
	doc, err := bsonkit.Transform(document)
	if err != nil {
		return nil, err
	}

	// ensure object id
	doc, id, err := ensureObjectID(doc)
	if err != nil {
		return nil, err
	}

	// write document
	err = c.client.backend.insertOne(c.ns, doc)
	if err != nil {
		return nil, err
	}

	return &mongo.InsertOneResult{
		InsertedID: id,
	}, nil
}

func (c *Collection) Name() string {
	panic("not implemented")
}

func (c *Collection) ReplaceOne(context.Context, interface{}, interface{}, ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	panic("not implemented")
}

func (c *Collection) UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("not implemented")
}

func (c *Collection) UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("not implemented")
}

func (c *Collection) Watch(context.Context, interface{}, ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	panic("not implemented")
}

func ensureObjectID(doc bson.D) (bson.D, primitive.ObjectID, error) {
	// check id
	var id primitive.ObjectID
	if v := bsonkit.Get(doc, "_id"); v != bsonkit.Missing {
		// check existing value
		oid, ok := v.(primitive.ObjectID)
		if !ok {
			return nil, oid, fmt.Errorf("only primitive.OjectID values are supported in _id field")
		} else if oid.IsZero() {
			return nil, oid, fmt.Errorf("found zero primitive.OjectID value in _id field")
		}

		// set id
		id = oid
	}

	// prepend id if zero
	if id.IsZero() {
		id = primitive.NewObjectID()
		bsonkit.Set(doc, "_id", id, true)
	}

	return doc, id, nil
}
