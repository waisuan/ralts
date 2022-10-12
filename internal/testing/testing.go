package testing

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"some-api/internal/db"
)

func ClearDB(dbClient *db.DatabaseClient) {
	_, err := dbClient.Client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{TableName: aws.String(db.TableName)})
	if err != nil {
		panic(err)
	}
}
