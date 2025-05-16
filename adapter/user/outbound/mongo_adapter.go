package userAdapter

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	userPort "github.com/wittawat/go-hex/core/port/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ------------------------ Entities ------------------------
type mongoUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at"`
}

type mongoUserRepository struct {
	collection *mongo.Collection
}

// ------------------------ Constructor ------------------------
func NewMongoUserRepository(col *mongo.Collection) userPort.UserRepository {
	return &mongoUserRepository{collection: col}
}

// ------------------------ Private Function -----------------------
func entities2MongoUser(user *entities.User) *mongoUser {
	return &mongoUser{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func mongo2EntitiesUser(mg *mongoUser) *entities.User {
	return &entities.User{
		ID:        string(mg.ID.Hex()),
		Username:  mg.Username,
		Email:     mg.Email,
		Password:  mg.Password,
		Role:      mg.Role,
		CreatedAt: mg.CreatedAt,
		UpdatedAt: mg.UpdatedAt,
		DeletedAt: mg.DeletedAt,
	}
}

// ------------------------ Method ------------------------
func (m *mongoUserRepository) setEmailUniqueIndex() error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	name, err := m.collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Println("Failed to create unique index on email:", err)
		return err
	}
	log.Println("Created index:", name)
	return nil
}

func (m *mongoUserRepository) Save(ctx context.Context, user *entities.User) error {
	if err := m.setEmailUniqueIndex(); err != nil {
		return err
	}
	_, err := m.collection.InsertOne(ctx, entities2MongoUser(user))
	if mongo.IsDuplicateKeyError(err) {
		return errors.New("email is already exists")
	}
	return err
}

func (m *mongoUserRepository) UpdateOne(ctx context.Context, user *entities.User, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.M{"$set": entities2MongoUser(user)}
	_, err = m.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (m *mongoUserRepository) DeleteOne(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (m *mongoUserRepository) FindById(ctx context.Context, id string) (*entities.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objID,
	}
	var mu mongoUser
	if err := m.collection.FindOne(ctx, filter).Decode(&mu); err != nil {
		log.Println("Fail to get user by id: ", err)
		return nil, err
	}
	return mongo2EntitiesUser(&mu), nil
}

func (m *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	filter := bson.M{
		"email": email,
	}
	var mu mongoUser
	if err := m.collection.FindOne(ctx, filter).Decode(&mu); err != nil {
		log.Println("Fail to get user by email: ", err)
		return nil, err
	}
	return mongo2EntitiesUser(&mu), nil
}

func (m *mongoUserRepository) Find(ctx context.Context) ([]entities.User, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []entities.User
	for cursor.Next(ctx) {
		var mu mongoUser
		if err := cursor.Decode(&mu); err != nil {
			continue
		}
		users = append(users, *mongo2EntitiesUser(&mu))
	}
	return users, nil
}
