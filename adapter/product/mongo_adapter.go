package productAdapter

import (
	"context"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ------------------------ Entities ------------------------
type MongoProduct struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Price     int32              `bson:"price"`
	Detail    string             `bson:"detail"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at"`
}

type MongoProductRepository struct {
	collection *mongo.Collection
}

// ------------------------ Constructor ------------------------
func NewMongoProductRepository(col *mongo.Collection) *MongoProductRepository {
	return &MongoProductRepository{collection: col}
}

// ------------------------ Private Function ------------------------
func entities2MongoProduct(p *entities.Product) (*MongoProduct, error) {
	objUserID, err := primitive.ObjectIDFromHex(p.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &MongoProduct{
		Title:     p.Title,
		Price:     p.Price,
		Detail:    p.Detail,
		CreatedBy: objUserID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
	}, nil
}

func mongo2EntitiesProduct(mp *MongoProduct) *entities.Product {
	return &entities.Product{
		ID:        string(mp.ID.Hex()),
		Title:     mp.Title,
		Price:     mp.Price,
		Detail:    mp.Detail,
		CreatedBy: string(mp.CreatedBy.Hex()),
		CreatedAt: mp.CreatedAt,
		UpdatedAt: mp.UpdatedAt,
		DeletedAt: mp.DeletedAt,
	}
}

// ------------------------ Method ------------------------
func (m *MongoProductRepository) Save(ctx context.Context, p *entities.Product) error {
	mp, err := entities2MongoProduct(p)
	if err != nil {
		return err
	}
	if _, err := m.collection.InsertOne(ctx, mp); err != nil {
		return err
	}
	return nil
}

func (m *MongoProductRepository) UpdateOne(ctx context.Context, p *entities.Product, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	mp, err := entities2MongoProduct(p)
	if err != nil {
		return err
	}
	update := bson.M{"$set": mp}
	_, err = m.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (m *MongoProductRepository) DeleteOne(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (m *MongoProductRepository) DeleteAll(ctx context.Context, userId string) error {
	objUserID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	filter := bson.M{
		"created_by": objUserID,
	}
	_, err = m.collection.DeleteMany(ctx, filter)
	return err
}

func (m *MongoProductRepository) FindById(ctx context.Context, id string) (*entities.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objID,
	}
	var mp MongoProduct
	if err := m.collection.FindOne(ctx, filter).Decode(&mp); err != nil {
		return nil, err
	}

	return mongo2EntitiesProduct(&mp), nil
}

func (m *MongoProductRepository) Find(ctx context.Context) ([]entities.Product, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []entities.Product
	for cursor.Next(ctx) {
		var mp MongoProduct
		if err := cursor.Decode(&mp); err != nil {
			continue
		}
		products = append(products, *mongo2EntitiesProduct(&mp))
	}
	return products, nil
}
