package productAdapter

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wittawat/go-hex/core/entities"
	productPort "github.com/wittawat/go-hex/core/port/product"
)

type redisProductRepository struct {
	redisClient *redis.Client
	dbRepo      productPort.ProductRepository
}

func NewRedisProductRepository(redisClient *redis.Client, dbRepo productPort.ProductRepository) productPort.ProductRepository {
	return &redisProductRepository{redisClient: redisClient, dbRepo: dbRepo}
}

func (r *redisProductRepository) Save(ctx context.Context, product *entities.Product) error {
	err := r.dbRepo.Save(ctx, product)
	if err != nil {
		return err
	}

	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	r.redisClient.Del(ctx, "redisAdapter::products")
	r.redisClient.Set(ctx, "redisAdapter::product:"+product.ID, string(data), time.Minute*2)
	return nil
}

func (r *redisProductRepository) UpdateOne(ctx context.Context, product *entities.Product, id string) error {
	err := r.dbRepo.UpdateOne(ctx, product, id)
	if err != nil {
		return err
	}

	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	r.redisClient.Set(ctx, "redisAdapter::product:"+id, string(data), time.Minute*2)
	r.redisClient.Del(ctx, "redisAdapter::products")
	return nil
}

func (r *redisProductRepository) DeleteOne(ctx context.Context, id string) error {
	err := r.dbRepo.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	r.redisClient.Del(ctx, "redisAdapter::product:"+id)
	r.redisClient.Del(ctx, "redisAdapter::products")
	return nil
}

func (r *redisProductRepository) DeleteAll(ctx context.Context, userId string) error {
	err := r.dbRepo.DeleteAll(ctx, userId)
	if err != nil {
		return err
	}
	r.redisClient.Del(ctx, "redisAdapter::products")
	return nil
}

func (r *redisProductRepository) FindById(ctx context.Context, id string) (*entities.Product, error) {
	key := "redisAdapter::product:" + id
	var product *entities.Product

	if productJson, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		if err = json.Unmarshal([]byte(productJson), product); err == nil {
			log.Println("get data from redis")
			return product, nil
		}
	}

	product, err := r.dbRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(product)
	if err != nil {
		return nil, err
	}

	if err = r.redisClient.Set(ctx, key, string(data), time.Minute*2).Err(); err != nil {
		return nil, err
	}
	log.Println("get data from database")
	return product, nil
}

func (r *redisProductRepository) Find(ctx context.Context) ([]entities.Product, error) {
	key := "redisAdapter::products"
	var products []entities.Product

	if productsJson, err := r.redisClient.Get(ctx, key).Result(); err == nil {
		if err = json.Unmarshal([]byte(productsJson), &products); err == nil {
			log.Println("get data from redis")
			return products, nil
		}
	}

	products, err := r.dbRepo.Find(ctx)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	err = r.redisClient.Set(ctx, key, string(data), time.Minute*2).Err()
	if err != nil {
		return nil, err
	}

	log.Println("get data from database")
	return products, nil
}
