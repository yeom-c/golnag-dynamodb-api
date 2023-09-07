package rule

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Interface interface {
	ConvertIoReaderToStruct(data io.Reader, model interface{}) (body interface{}, err error)
	GetMock() interface{}
	Migrate(connection *dynamodb.Client) error
	Validate(model interface{}) error
}
