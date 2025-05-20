package productAdapter

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	productPort "github.com/wittawat/go-hex/core/port/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ------------------------ Entities ------------------------
type mongoProduct struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Price     int32              `bson:"price"`
	Detail    string             `bson:"detail"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at"`
}

type mongoProductRepository struct {
	collection *mongo.Collection
}

// ------------------------ Constructor ------------------------
func NewMongoProductRepository(col *mongo.Collection, userCol *mongo.Collection) productPort.ProductRepository {
	if err := productFactoryMongo(col, userCol); err != nil {
		log.Println("failed to feed product to mongo: ", err)
	}
	return &mongoProductRepository{collection: col}
}

// ------------------------ Private Function ------------------------
func getSellerFormMongo(userCol *mongo.Collection) ([]primitive.ObjectID, error) {
	ctx := context.Background()
	filter := bson.M{"role": "seller"}
	cursor, err := userCol.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("failed to find sellers")
	}
	defer cursor.Close(ctx)

	type SellerID struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	var sellerDocs []SellerID
	if err := cursor.All(ctx, &sellerDocs); err != nil {
		return nil, errors.New("failed to decode sellers")
	}

	if len(sellerDocs) == 0 {
		return nil, errors.New("no seller found in users collection")
	}

	var sellerIDs []primitive.ObjectID
	for _, seller := range sellerDocs {
		sellerIDs = append(sellerIDs, seller.ID)
	}
	log.Println("product::getSellerFromMongo: success")
	return sellerIDs, nil
}

func productFactoryMongo(col *mongo.Collection, userCol *mongo.Collection) error {
	sellers, err := getSellerFormMongo(userCol)

	if err != nil {
		return err
	}

	count, err := col.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	var products []interface{}
	for i := 1; i <= 10; i++ {
		iStr := strconv.Itoa(i)
		randSell := rand.IntN(len(sellers))
		product := mongoProduct{
			Title:     "product" + iStr,
			Price:     int32(i * 10),
			Detail:    "detail for product" + iStr,
			CreatedBy: sellers[randSell],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}
		products = append(products, product)
	}
	col.InsertMany(context.Background(), products)
	log.Println("feed product to mongo: success")
	return nil
}

func entities2MongoProduct(p *entities.Product) (*mongoProduct, error) {
	objUserID, err := primitive.ObjectIDFromHex(p.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &mongoProduct{
		Title:     p.Title,
		Price:     p.Price,
		Detail:    p.Detail,
		CreatedBy: objUserID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
	}, nil
}

func mongo2EntitiesProduct(mp *mongoProduct) *entities.Product {
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
func (m *mongoProductRepository) Save(ctx context.Context, p *entities.Product) error {
	mp, err := entities2MongoProduct(p)
	if err != nil {
		return err
	}
	if _, err := m.collection.InsertOne(ctx, mp); err != nil {
		return err
	}
	return nil
}

func (m *mongoProductRepository) UpdateOne(ctx context.Context, p *entities.Product, id string) error {
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

func (m *mongoProductRepository) DeleteOne(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (m *mongoProductRepository) DeleteAll(ctx context.Context, userId string) error {
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

func (m *mongoProductRepository) FindById(ctx context.Context, id string) (*entities.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objID,
	}
	var mp mongoProduct
	if err := m.collection.FindOne(ctx, filter).Decode(&mp); err != nil {
		return nil, err
	}

	return mongo2EntitiesProduct(&mp), nil
}

func (m *mongoProductRepository) Find(ctx context.Context) ([]entities.Product, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []entities.Product
	for cursor.Next(ctx) {
		var mp mongoProduct
		if err := cursor.Decode(&mp); err != nil {
			continue
		}
		products = append(products, *mongo2EntitiesProduct(&mp))
	}
	return products, nil
}
