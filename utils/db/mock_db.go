package db

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type MockDatabaseClient struct {
	GetByPkFn    func(pk string) (map[string]types.AttributeValue, error)
	UpdateItemFn func(pk string, attributes map[string]string) error
}

func (md *MockDatabaseClient) GetByPk(pk string) (map[string]types.AttributeValue, error) {
	return md.GetByPkFn(pk)
}

func (md *MockDatabaseClient) UpdateItem(pk string, attributes map[string]string) error {
	return md.UpdateItemFn(pk, attributes)
}
