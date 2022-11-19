package chat

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
	"some-api/internal/db"
	"time"
)

type Message struct {
	Pk        string `json:"-"`
	Sk        string `json:"-"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type Messages []Message

type Chat struct {
	db *db.DatabaseClient
}

func NewChat(db *db.DatabaseClient) *Chat {
	return &Chat{
		db: db,
	}
}

func (c *Chat) LoadAllMessages(day time.Time, now func() time.Time) (Messages, error) {
	var messages Messages

	pk := day.Format("20060102")
	sk := now().Add(-(time.Hour * 3)).Format("20060102150405")

	out, err := c.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(db.TableName),
		KeyConditionExpression: aws.String("pk = :pk and sk >= :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
			":sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		log.Error(fmt.Sprintf("Unable to load chat messages -> %s", err.Error()))
		return nil, err
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &messages)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to unmarshal chat messages -> %s", err.Error()))
		return nil, err
	}

	log.Info(fmt.Sprintf("%s, %s, %v", pk, sk, len(messages)))

	return messages, nil
}

func (c *Chat) SaveMessage(userId string, text string, now func() time.Time) error {
	today := now().Format("20060102")
	createdAt := now().Format("20060402150405")

	_, err := c.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(db.TableName),
		Item: map[string]types.AttributeValue{
			"pk":        &types.AttributeValueMemberS{Value: today},
			"sk":        &types.AttributeValueMemberS{Value: createdAt},
			"userId":    &types.AttributeValueMemberS{Value: userId},
			"text":      &types.AttributeValueMemberS{Value: text},
			"createdAt": &types.AttributeValueMemberS{Value: createdAt},
		},
	})
	if err != nil {
		log.Error(fmt.Sprintf("Unable to save chat message -> %s", err.Error()))
		return err
	}

	return nil
}
