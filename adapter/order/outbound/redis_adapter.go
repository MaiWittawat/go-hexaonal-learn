package orderAdapter

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wittawat/go-hex/core/entities"
	orderPort "github.com/wittawat/go-hex/core/port/order"
)

type redisOrderRepository struct {
	redisClient *redis.Client
	dbRepo      orderPort.OrderRepository
}

func NewRedisOrderRepository(redisClient *redis.Client, dbRepo orderPort.OrderRepository) orderPort.OrderRepository {
	return &redisOrderRepository{redisClient: redisClient, dbRepo: dbRepo}
}

func (r *redisOrderRepository) Save(ctx context.Context, order *entities.Order) error {
	err := r.dbRepo.Save(ctx, order)
	if err != nil {
		return err
	}

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	r.redisClient.Del(ctx, "redisAdapter::orders")
	r.redisClient.Set(ctx, "redisAdapter::order:"+order.ID, string(data), time.Minute*2)
	return nil
}

func (r *redisOrderRepository) FindByUserId(ctx context.Context, userId string) ([]entities.Product, error) {
	return nil, nil
}

// ****************** Problem Check Pls *******************
// Why find the user email in order table they have on field email for query
func (r *redisOrderRepository) FindByUserEmail(ctx context.Context, email string) (*entities.Order, error) {
	key := "redisAdapter::order:" + email
	var order *entities.Order

	if orderJson, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		if err = json.Unmarshal([]byte(orderJson), order); err == nil {
			log.Println("get data from redis")
			return order, nil
		}
	}

	order, err := r.dbRepo.FindByUserEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	if err = r.redisClient.Set(ctx, key, string(data), time.Minute*2).Err(); err != nil {
		return nil, err
	}
	log.Println("get data from database")
	return order, nil
}

func (r *redisOrderRepository) FindById(ctx context.Context, id string) (*entities.Order, error) {
	key := "redisAdapter::order:" + id
	var order *entities.Order

	if orderJson, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		if err = json.Unmarshal([]byte(orderJson), order); err == nil {
			log.Println("get data from redis")
			return order, nil
		}
	}

	order, err := r.dbRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	if err = r.redisClient.Set(ctx, key, string(data), time.Minute*2).Err(); err != nil {
		return nil, err
	}
	log.Println("get data from database")
	return order, nil
}

func (r *redisOrderRepository) UpdateOne(ctx context.Context, order *entities.Order, id string) error {
	err := r.dbRepo.UpdateOne(ctx, order, id)
	if err != nil {
		return err
	}

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	r.redisClient.Set(ctx, "redisAdapter::order:"+id, string(data), time.Minute*2)
	r.redisClient.Del(ctx, "redisAdapter::orders")
	return nil
}

func (r *redisOrderRepository) DeleteOne(ctx context.Context, id string) error {
	err := r.dbRepo.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	r.redisClient.Del(ctx, "redisAdapter::order:"+id)
	r.redisClient.Del(ctx, "redisAdapter::orders")
	return nil
}

func (r *redisOrderRepository) DeleteAllOrderByUser(ctx context.Context, userId string) error {
	err := r.dbRepo.DeleteAllOrderByUser(ctx, userId)
	if err != nil {
		return err
	}
	r.redisClient.Del(ctx, "redisAdapter::orders")
	return nil
}

func (r *redisOrderRepository) DeleteAllOrderByProduct(ctx context.Context, productId string) error {
	err := r.dbRepo.DeleteAllOrderByProduct(ctx, productId)
	if err != nil {
		return err
	}
	r.redisClient.Del(ctx, "redisAdapter::orders")
	return nil
}
