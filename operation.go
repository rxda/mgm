package mgm

import (
	"context"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func create(ctx context.Context, c *Collection, model Model, opts ...*options.InsertOneOptions) error {
	// Call to saving hook
	if err := callToBeforeCreateHooks(ctx, model); err != nil {
		return err
	}

	res, err := c.InsertOne(ctx, model, opts...)

	if err != nil {
		return err
	}

	// Set new id
	model.SetID(res.InsertedID)

	return callToAfterCreateHooks(ctx, model)
}

func createMany(ctx context.Context, c *Collection, models []Model, opts ...*options.InsertManyOptions) error {
	// Call to saving hook
	for k := range models {
		if err := callToBeforeCreateHooks(ctx, models[k]); err != nil {
			return err
		}
	}

	res, err := c.InsertMany(ctx, ModelsToInterfaces(models), opts...)

	if err != nil {
		return err
	}
	for i := range models {
		models[i].SetID(res.InsertedIDs[i])
	}

	for k := range models {
		if err := callToAfterCreateHooks(ctx, models[k]); err != nil {
			return err
		}
	}

	return nil
}

func first(ctx context.Context, c *Collection, filter interface{}, model Model, opts ...*options.FindOneOptions) error {
	return c.FindOne(ctx, filter, opts...).Decode(model)
}

func update(ctx context.Context, c *Collection, model Model, opts ...*options.UpdateOptions) error {
	// Call to saving hook
	if err := callToBeforeUpdateHooks(ctx, model); err != nil {
		return err
	}

	res, err := c.UpdateOne(ctx, bson.M{field.ID: model.GetID()}, bson.M{"$set": model}, opts...)

	if err != nil {
		return err
	}

	return callToAfterUpdateHooks(ctx, res, model)
}

func del(ctx context.Context, c *Collection, model Model) error {
	if err := callToBeforeDeleteHooks(ctx, model); err != nil {
		return err
	}
	res, err := c.DeleteOne(ctx, bson.M{field.ID: model.GetID()})
	if err != nil {
		return err
	}

	return callToAfterDeleteHooks(ctx, res, model)
}
