package orm

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var setupIndexFunc []func(im *IndexModifier) error

type IndexModifier struct {
	Key    bson.D
	Option *options.IndexOptions
}

func IndexKey(i any) string {
	if _, ok := i.(bson.D); !ok {
		panic("the key must be the bson.D type")
	}
	return buildKeyFromD(i.(bson.D))
}

func buildKeyFromD(k bson.D) string {
	// hack the key builder
	// see the default behavior in mongo.IndexModel.Options.Name
	// > The default value is "[field1]_[direction1]_[field2]_[direction2]..."
	var keys []string
	for _, e := range k {
		key := bytes.NewBuffer(nil)
		key.WriteString(e.Key)
		key.WriteString("_")
		direction := strconv.FormatInt(int64(e.Value.(int)), 10)
		key.WriteString(direction)
		keys = append(keys, key.String())
	}
	return strings.Join(keys, "_")
}

func SetupIndexModel(c *mongo.Collection, keys []mongo.IndexModel) {
	ctx := context.Background()
	setupIndexFunc = append(setupIndexFunc, func(im *IndexModifier) error {
		for index, k := range keys {
			if buildKeyFromD(im.Key) == buildKeyFromD(k.Keys.(bson.D)) && im.Option != nil {
				keys[index].Options = im.Option
			}
		}
		if err := ensureIndex(ctx, c, keys); err != nil {
			return err
		}
		return nil
	})
}

type IndexOptionsHandler func(im *IndexModifier) error

func WithIndexNameHandler(index bson.D, opt *options.IndexOptions) IndexOptionsHandler {
	return func(im *IndexModifier) error {
		im.Key = index
		im.Option = opt
		return nil
	}
}

func EnsureAllIndex(fns ...IndexOptionsHandler) error {
	opt := &IndexModifier{}
	for _, f := range fns {
		f(opt)
	}
	for _, f := range setupIndexFunc {
		if err := f(opt); err != nil {
			return err
		}
	}
	return nil
}

// ensureIndex will ensure the index model provided is on the given collection.
// we should directly create the provided index , per the mongo doc, mongo will :
// 1. (before version v4.2), create the index even if the index (name) already exist , and not returns any errors.
// 2. (after version v4.2), returns an already exist error if the index (name) already exist , but the duplicate index will still be created.
func ensureIndex(ctx context.Context, c *mongo.Collection, keys []mongo.IndexModel) error {

	idxs := c.Indexes()
	cur, err := idxs.List(ctx, options.ListIndexes().SetBatchSize(1))
	if err != nil {
		return fmt.Errorf("ensureIndex: unable to list indexes: %w", err)
	}

	var idx []bson.M
	if err := cur.All(ctx, &idx); err != nil {
		return fmt.Errorf("ensureIndex: range indexes: %w", err)
	}

	exKeys := make(map[string]mongo.IndexModel)
	for _, m := range idx {
		exKeys[m["name"].(string)] = mongo.IndexModel{}
	}

	for _, k := range keys {
		var keyName string
		switch {
		case k.Options == nil, k.Options.Name == nil:
			keyName = buildKeyFromD(k.Keys.(bson.D))
		case k.Options.Name != nil:
			keyName = *k.Options.Name
		}
		if _, ok := exKeys[keyName]; !ok {
			if _, err := idxs.CreateOne(ctx, k); err != nil {
				return fmt.Errorf("ensureIndex: unable to create index: %w", err)
			}
		}
	}

	return nil
}
