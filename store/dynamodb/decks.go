package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	handy "github.com/ztsu/handy-go/store"
)

const (
	decksTableName = "Decks"
)

type DeckDynamoDBStore struct {
	db *dynamodb.DynamoDB
}

func NewDeckDynamoDBStore(accessKeyID, secret string) (*DeckDynamoDBStore, error) {
	cfg := &aws.Config{
		Region:                        aws.String(endpoints.EuCentral1RegionID),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials(accessKeyID, secret, ""),
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	db := dynamodb.New(sess)

	return &DeckDynamoDBStore{db}, nil
}

func (store *DeckDynamoDBStore) Add(deck *handy.Deck) error {
	deckAv, err := dynamodbattribute.MarshalMap(deck)
	if err != nil {
		return err
	}

	idNotExistsExpr, _ := notExists("id")

	putDeck := &dynamodb.Put{
		TableName:                aws.String(decksTableName),
		Item:                     deckAv,
		ConditionExpression:      idNotExistsExpr.Condition(),
		ExpressionAttributeNames: idNotExistsExpr.Names(),
	}

	txInp := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: putDeck},
		},
	}

	_, err = store.db.TransactWriteItems(txInp)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
			case dynamodb.ErrCodeTransactionCanceledException: // TODO Extract CancellationReasons
				return handy.ErrDeckAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (store *DeckDynamoDBStore) Get(id string) (*handy.Deck, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(decksTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}

	output, err := store.db.GetItem(input)
	if err != nil {
		return nil, err
	}

	deck := &handy.Deck{}

	err = dynamodbattribute.UnmarshalMap(output.Item, deck)
	if err != nil {
		return nil, err
	}

	if deck.ID == "" {
		return nil, handy.ErrDeckNotFound
	}

	return deck, nil
}

func (store *DeckDynamoDBStore) Save(deck *handy.Deck) error {
	deckAv, err := dynamodbattribute.MarshalMap(deck)
	if err != nil {
		return err
	}

	idExists, _ := exists("id")

	putDeck := &dynamodb.Put{
		TableName:                aws.String(decksTableName),
		Item:                     deckAv,
		ConditionExpression:      idExists.Condition(),
		ExpressionAttributeNames: idExists.Names(),
	}

	txInp := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: putDeck},
		},
	}

	_, err = store.db.TransactWriteItems(txInp)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case dynamodb.ErrCodeTransactionCanceledException:
				return handy.ErrDeckAlreadyExists
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return handy.ErrDeckNotFound
			}
		}

		return err
	}

	return nil
}

func (store *DeckDynamoDBStore) Delete(id string) error {
	deck, err := store.Get(id)
	if err != nil {
		return err
	}

	idExistsExpr, _ := exists("id")

	deleteDeck := &dynamodb.Delete{
		TableName: aws.String(decksTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(deck.ID)},
		},
		ConditionExpression:      idExistsExpr.Condition(),
		ExpressionAttributeNames: idExistsExpr.Names(),
	}

	txInp := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Delete: deleteDeck},
		},
	}

	_, err = store.db.TransactWriteItems(txInp)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return handy.ErrDeckNotFound
			}
		}

		return err
	}

	return nil
}
