package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	Config "github.com/yeom-c/golnag-dynamodb-api/config"
	"github.com/yeom-c/golnag-dynamodb-api/internal/repository/adapter"
	"github.com/yeom-c/golnag-dynamodb-api/internal/repository/instance"
	"github.com/yeom-c/golnag-dynamodb-api/internal/route"
	"github.com/yeom-c/golnag-dynamodb-api/internal/rule"
	"github.com/yeom-c/golnag-dynamodb-api/internal/rule/product"
	"github.com/yeom-c/golnag-dynamodb-api/util/logger"
)

func main() {
	config := Config.GetConfig()
	connection := instance.GetConnection()
	repository := adapter.NewAdapter(connection)

	logger.INFO("waiting for the service to start...", nil)

	errors := Migrate(connection)
	if len(errors) > 0 {
		for _, err := range errors {
			logger.PANIC("Error on migration:...", err)
		}
	}

	logger.PANIC("", checkTables(connection))

	port := fmt.Sprintf(":%d", config.Port)
	router := route.NewRouter().SetRouter(repository)
	logger.INFO("service is running on port", port)

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}

func Migrate(connection *dynamodb.Client) []error {
	errors := []error{}
	callMigrateAndAppendError(&errors, connection, &product.Rule{})

	return errors
}

func callMigrateAndAppendError(errors *[]error, connection *dynamodb.Client, rule rule.Interface) {
	err := rule.Migrate(connection)
	if err != nil {
		*errors = append(*errors, err)
	}
}

func checkTables(connection *dynamodb.Client) error {
	response, err := connection.ListTables(context.Background(), &dynamodb.ListTablesInput{})
	if response != nil {
		if len(response.TableNames) == 0 {
			logger.INFO("Tables not found", nil)
		}

		for _, tableName := range response.TableNames {
			logger.INFO("Table found: ", tableName)
		}
	}

	return err
}
