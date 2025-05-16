package orderAdapter

import (
	"context"
	"fmt"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ------------------------ Entities ------------------------
type mongoOrder struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ProductID primitive.ObjectID `bson:"product_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at"`
}

type mongoOrderRepository struct {
	collection *mongo.Collection
}

// ------------------------ Constructor ------------------------
func NewMongoOrderRepository(col *mongo.Collection) orderPort.OrderRepository {
	return &mongoOrderRepository{collection: col}
}

// ------------------------ Private Function -----------------------
func entities2MongoOrder(order *entities.Order) (*mongoOrder, error) {
	userID, err := primitive.ObjectIDFromHex(order.UserID)
	if err != nil {
		return nil, err
	}

	productID, err := primitive.ObjectIDFromHex(order.ProductID)
	if err != nil {
		return nil, err
	}

	return &mongoOrder{
		UserID:    userID,
		ProductID: productID,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,
	}, nil
}

func mongo2EntitiesOrder(mo *mongoOrder) *entities.Order {
	return &entities.Order{
		ID:        string(mo.ID.Hex()),
		UserID:    string(mo.UserID.Hex()),
		ProductID: string(mo.ProductID.Hex()),
		CreatedAt: mo.CreatedAt,
		UpdatedAt: mo.UpdatedAt,
		DeletedAt: mo.DeletedAt,
	}
}

func toEntitiesProduct(raw bson.M) (*entities.Product, error) {
	id, ok := raw["ID"].(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("invalid ID type")
	}

	return &entities.Product{
		ID:     id.Hex(),
		Title:  raw["title"].(string),
		Price:  raw["price"].(int32),
		Detail: raw["detail"].(string),
	}, nil
}

// ------------------------ Method ------------------------
func (m *mongoOrderRepository) Save(ctx context.Context, order *entities.Order) error {
	mo, err := entities2MongoOrder(order)
	if err != nil {
		return err
	}
	_, err = m.collection.InsertOne(ctx, mo)
	return err
}

func (m *mongoOrderRepository) UpdateOne(ctx context.Context, order *entities.Order, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	mo, err := entities2MongoOrder(order)
	if err != nil {
		return err
	}
	update := bson.M{"$set": mo}
	_, err = m.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (m *mongoOrderRepository) DeleteOne(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (m *mongoOrderRepository) DeleteAllOrderByUser(ctx context.Context, userId string) error {
	objUserID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteMany(ctx, bson.M{"user_id": objUserID})
	return err
}

func (m *mongoOrderRepository) DeleteAllOrderByProduct(ctx context.Context, productId string) error {
	objProductID, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteMany(ctx, bson.M{"product_id": objProductID})
	return err
}

func (m *mongoOrderRepository) FindById(ctx context.Context, id string) (*entities.Order, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	var mo mongoOrder
	if err := m.collection.FindOne(ctx, filter).Decode(&mo); err != nil {
		return nil, err
	}
	return mongo2EntitiesOrder(&mo), nil
}

func (m *mongoOrderRepository) FindByUserEmail(ctx context.Context, email string) (*entities.Order, error) {
	var mo mongoOrder
	filter := bson.M{"email": email}
	if err := m.collection.FindOne(ctx, filter).Decode(&mo); err != nil {
		return nil, err
	}
	return mongo2EntitiesOrder(&mo), nil

}

func (m *mongoOrderRepository) FindByUserId(ctx context.Context, userID string) ([]entities.Product, error) {
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var products []entities.Product

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: objUserID}}}},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "products"},
				{Key: "localField", Value: "product_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "product"},
			},
		}},
		{{Key: "$unwind", Value: "$product"}},
		{{
			Key: "$project", Value: bson.D{
				{Key: "ID", Value: "$product._id"},
				{Key: "title", Value: "$product.title"},
				{Key: "price", Value: "$product.price"},
				{Key: "detail", Value: "$product.detail"},
			},
		}},
	}

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, err
		}

		product, err := toEntitiesProduct(raw)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
