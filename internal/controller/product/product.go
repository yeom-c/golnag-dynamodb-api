package product

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/yeom-c/golnag-dynamodb-api/internal/entity/product"
	"github.com/yeom-c/golnag-dynamodb-api/internal/repository/adapter"
)

type Controller struct {
	repository adapter.Interface
}

type Interface interface {
	ListAll() (entities []product.Product, err error)
	ListOne(id uuid.UUID) (entity product.Product, err error)
	Create(entity *product.Product) (string, error)
	Update(id uuid.UUID, entity *product.Product) error
	Delete(id uuid.UUID) error
}

func NewController(repository adapter.Interface) Interface {
	return &Controller{
		repository,
	}
}

func (c *Controller) ListAll() (entities []product.Product, err error) {
	entities = []product.Product{}
	var entityProduct product.Product

	filter := expression.Name("name").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return entities, err
	}

	response, err := c.repository.FindAll(condition, entityProduct.TableName())
	if err != nil {
		return entities, err
	}

	if response != nil {
		for _, value := range response.Items {
			entity, err := product.ParseDynamoAttributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}
	return entities, nil
}

func (c *Controller) ListOne(id uuid.UUID) (entity product.Product, err error) {
	entity.ID = id.String()
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return product.ParseDynamoAttributeToStruct(response.Item)
}

func (c *Controller) Create(entity *product.Product) (string, error) {
	entity.CreatedAt = time.Now()
	c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.ID, nil
}

func (c *Controller) Update(id uuid.UUID, entity *product.Product) error {
	found, err := c.ListOne(id)
	if err != nil {
		return err
	}

	found.ID = id.String()
	found.Name = entity.Name
	found.UpdatedAt = time.Now()
	_, err = c.repository.CreateOrUpdate(found.GetMap(), found.TableName())

	return err
}

func (c *Controller) Delete(id uuid.UUID) error {
	entity, err := c.ListOne(id)
	if err != nil {
		return err
	}

	_, err = c.repository.Delete(entity.GetFilterId(), entity.TableName())

	return err
}
