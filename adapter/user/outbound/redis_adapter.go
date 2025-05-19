package userAdapter

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wittawat/go-hex/core/entities"
	userPort "github.com/wittawat/go-hex/core/port/user"
)

type redisUserRepository struct {
	redisClient *redis.Client
	dbRepo      userPort.UserRepository
}

func NewRedisUserRepository(redisClient *redis.Client, dbRepo userPort.UserRepository) userPort.UserRepository {
	return &redisUserRepository{redisClient: redisClient, dbRepo: dbRepo}
}

func (r *redisUserRepository) Save(ctx context.Context, user *entities.User) error {
	err := r.dbRepo.Save(ctx, user)
	if err != nil {
		return err
	}
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	r.redisClient.Del(ctx, "redisAdapter::users")
	r.redisClient.Set(ctx, "redisAdapter::user:"+user.ID, string(data), time.Minute*2)
	return nil
}

func (r *redisUserRepository) FindById(ctx context.Context, id string) (*entities.User, error) {
	key := "redisAdapter::user:" + id
	var user *entities.User

	if userJson, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		if err = json.Unmarshal([]byte(userJson), user); err == nil {
			log.Println("get data from redis")
			return user, nil
		}
	}

	user, err := r.dbRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	if err = r.redisClient.Set(ctx, key, string(data), time.Minute*2).Err(); err != nil {
		return nil, err
	}

	log.Panicln("get data from database")
	return user, nil
}

func (r *redisUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	key := "redisAdapter::user:" + email
	var user *entities.User

	if userJson, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		if err = json.Unmarshal([]byte(userJson), user); err == nil {
			log.Println("get data from redis")
			return user, nil
		}
	}

	user, err := r.dbRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	if err = r.redisClient.Set(ctx, key, string(data), time.Minute*2).Err(); err != nil {
		return nil, err
	}

	log.Println("get data from database")
	return user, nil
}

func (r *redisUserRepository) Find(ctx context.Context) ([]entities.User, error) {
	key := "redisAdapter::products"
	var users []entities.User

	if usersJson, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		if err = json.Unmarshal([]byte(usersJson), &users); err == nil {
			log.Println("get data from redis")
			return users, nil
		}
	}

	users, err := r.dbRepo.Find(ctx)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	err = r.redisClient.Set(ctx, key, string(data), time.Minute*2).Err()
	if err != nil {
		return nil, err
	}

	log.Println("get data from database")
	return users, nil
}

func (r *redisUserRepository) UpdateOne(ctx context.Context, user *entities.User, id string) error {
	err := r.dbRepo.UpdateOne(ctx, user, id)
	if err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	r.redisClient.Set(ctx, "redisAdapter::user:"+id, string(data), time.Minute*2)
	r.redisClient.Del(ctx, "redisAdapter::users")
	return nil
}

func (r *redisUserRepository) DeleteOne(ctx context.Context, id string) error {
	err := r.dbRepo.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	r.redisClient.Del(ctx, "redisAdapter::user:"+id)
	r.redisClient.Del(ctx, "redisAdapter::users")
	return nil
}
