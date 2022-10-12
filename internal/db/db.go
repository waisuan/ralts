package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strings"
)

const (
	region         = "eu-central-1"
	localEndpoint  = "http://localhost:8000"
	localAccessKey = "dummy"
	localSecretKey = "dummy"
	TableName      = "locations"
	pkFieldName    = "pk"
)

type DataStore interface {
	GetByPk(pk string) (map[string]types.AttributeValue, error)
	UpdateItem(pk string, attributes map[string]string) error
}

type DatabaseClient struct {
	Client *dynamodb.Client
}

func New() *DatabaseClient {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg, func(options *dynamodb.Options) {
		options.Region = region
		options.Credentials = credentials.StaticCredentialsProvider{
			Value: aws.Credentials{AccessKeyID: localAccessKey, SecretAccessKey: localSecretKey},
		}
		options.EndpointResolver = dynamodb.EndpointResolverFromURL(localEndpoint)
	})

	_, err = svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(TableName)})
	if err != nil {
		_, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("pk"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("sk"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("pk"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("sk"),
					KeyType:       types.KeyTypeRange,
				},
			},
			TableName:   aws.String(TableName),
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			panic(err)
		}
	}

	return &DatabaseClient{Client: svc}
}

func (db *DatabaseClient) GetByPk(pk string) (map[string]types.AttributeValue, error) {
	out, err := db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			pkFieldName: &types.AttributeValueMemberS{Value: pk},
		},
	})

	if err != nil {
		return nil, err
	}

	return out.Item, nil
}

func (db *DatabaseClient) UpdateItem(pk string, attributes map[string]string) error {
	var updateExprAttr []string
	exprAttrVal := make(map[string]types.AttributeValue)
	for k, v := range attributes {
		value := fmt.Sprintf(":%s", k)
		updateExprAttr = append(updateExprAttr, fmt.Sprintf("%s = %s", k, value))
		exprAttrVal[value] = &types.AttributeValueMemberS{Value: v}
	}

	updateExpr := "set " + strings.Join(updateExprAttr, ", ")

	_, err := db.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			pkFieldName: &types.AttributeValueMemberS{Value: pk},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprAttrVal,
	})

	if err != nil {
		return err
	}

	return nil
}
