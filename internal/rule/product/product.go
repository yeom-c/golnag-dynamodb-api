package product

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/yeom-c/golnag-dynamodb-api/internal/entity"
	"github.com/yeom-c/golnag-dynamodb-api/internal/entity/product"
)

type Rule struct{}

func NewRule() *Rule {
	return &Rule{}
}

func (r *Rule) ConvertIoReaderToStruct(data io.Reader, model interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("body is invalid")
	}
	return model, json.NewDecoder(data).Decode(model)
}

func (r *Rule) GetMock() interface{} {
	return product.Product{
		Base: entity.Base{
			ID:        uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: uuid.New().String(),
	}
}

func (r *Rule) Migrate(connection *dynamodb.Client) error {
	return r.createTable(connection)
}

func (r *Rule) Validate(model interface{}) error {
	productModel, err := product.InterfaceToModel(model)
	if err != nil {
		return err
	}

	return validation.ValidateStruct(productModel,
		validation.Field(&productModel.ID, validation.Required, is.UUIDv4),
		validation.Field(&productModel.Name, validation.Required, validation.Length(3, 50)),
	)
}

func (r *Rule) createTable(connection *dynamodb.Client) error {
	entity := &product.Product{}
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("_id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("_id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(entity.TableName()),
	}

	response, err := connection.CreateTable(context.Background(), input)
	if err != nil && strings.Contains(err.Error(), "Table already exists") {
		return nil
	}
	if response != nil && response.TableDescription.TableStatus == types.TableStatusCreating {
		time.Sleep(3 * time.Second)
		err = r.createTable(connection)
		if err != nil {
			return err
		}
	}

	return err
}
