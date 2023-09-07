package product

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/yeom-c/golnag-dynamodb-api/internal/entity"
)

type Product struct {
	entity.Base
	Name string `json:"name"  dynamodbav:"name"`
}

func InterfaceToModel(data interface{}) (instance *Product, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"_id": p.ID}
}

func (p *Product) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Product) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       p.ID,
		"name":      p.Name,
		"createdAt": p.CreatedAt.Format(entity.GetTimeFormat()),
		"updatedAt": p.UpdatedAt.Format(entity.GetTimeFormat()),
	}
}

func ParseDynamoAttributeToStruct(response map[string]types.AttributeValue) (p Product, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return p, errors.New("Item not found")
	}
	err = attributevalue.UnmarshalMap(response, &p)
	// for key, value := range response {
	// 	logger.INFO("key ", key)
	// 	logger.INFO("value ", value)
	// 	// if key == "_id" {
	// 	// 	p.ID, err = uuid.Parse()
	// 	// 	if p.ID == uuid.Nil {
	// 	// 		err = errors.New("Item not found")
	// 	// 	}
	// 	// }
	// 	// if key == "name" {
	// 	// 	p.Name = value["S"]
	// 	// }
	// 	// if key == "createdAt" {
	// 	// 	p.CreatedAt, err = time.Parse(entity.GetTimeFormat(), *value.S)
	// 	// }
	// 	// if key == "updatedAt" {
	// 	// 	p.UpdatedAt, err = time.Parse(entity.GetTimeFormat(), *value.S)
	// 	// }
	// 	// if err != nil {
	// 	// 	return p, err
	// 	// }
	// }

	return p, err
}
