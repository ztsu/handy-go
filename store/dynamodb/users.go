package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	handy "github.com/ztsu/handy-go/store"
)

const (
	usersTableName       = "Users"
	usersEmailsTableName = "UsersEmails"
)

type emailIdx struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

type UserDynamoDBStore struct {
	db *dynamodb.DynamoDB
}

func NewUserDynamoDBStore(accessKeyID, secret string) (*UserDynamoDBStore, error) {
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

	return &UserDynamoDBStore{db}, nil
}

func (store *UserDynamoDBStore) Get(id string) (*handy.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(usersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	}

	output, err := store.db.GetItem(input)
	if err != nil {
		return nil, err
	}

	user := &handy.User{}

	err = dynamodbattribute.UnmarshalMap(output.Item, user)
	if err != nil {
		return nil, err
	}

	if user.ID == "" {
		return nil, handy.ErrUserNotFound
	}

	return user, nil
}

func (store *UserDynamoDBStore) Add(user *handy.User) error {
	userAv, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	idNotExistsExpr, _ := notExists("id")
	emailNotExistsExpr, _ := notExists("email")

	putUser := &dynamodb.Put{
		TableName:                aws.String(usersTableName),
		Item:                     userAv,
		ConditionExpression:      idNotExistsExpr.Condition(),
		ExpressionAttributeNames: idNotExistsExpr.Names(),
	}

	emailIdxAv, err := dynamodbattribute.MarshalMap(emailIdx{Email: user.Email, ID: user.ID})
	if err != nil {
		return err
	}

	putEmail := &dynamodb.Put{
		TableName:                aws.String(usersEmailsTableName),
		Item:                     emailIdxAv,
		ConditionExpression:      emailNotExistsExpr.Condition(),
		ExpressionAttributeNames: emailNotExistsExpr.Names(),
	}

	txInp := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: putUser},
			{Put: putEmail},
		},
	}

	_, err = store.db.TransactWriteItems(txInp)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
			case dynamodb.ErrCodeTransactionCanceledException: // TODO Extract CancellationReasons
				return handy.ErrUserAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (store *UserDynamoDBStore) Save(user *handy.User) error {
	existed, err := store.Get(user.ID)
	if err != nil {
		return err
	}

	userAv, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	idExists, _ := exists("id")

	putUser := &dynamodb.Put{
		TableName:                aws.String(usersTableName),
		Item:                     userAv,
		ConditionExpression:      idExists.Condition(),
		ExpressionAttributeNames: idExists.Names(),
	}

	txInp := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Put: putUser},
			//,
		},
	}

	if user.Email != existed.Email {
		emailNotExistsExpr, _ := notExists("email")

		emailIdxAv, err := dynamodbattribute.MarshalMap(emailIdx{Email: user.Email, ID: user.ID})
		if err != nil {
			return err
		}

		putEmail := &dynamodb.Put{
			TableName:                aws.String(usersEmailsTableName),
			Item:                     emailIdxAv,
			ConditionExpression:      emailNotExistsExpr.Condition(),
			ExpressionAttributeNames: emailNotExistsExpr.Names(),
		}

		deleteOldEmail := &dynamodb.Delete{
			TableName: aws.String(usersEmailsTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"email": {S: aws.String(existed.Email)},
			},
		}

		txInp.TransactItems = append(
			txInp.TransactItems,
			&dynamodb.TransactWriteItem{
				Put: putEmail,
			},
			&dynamodb.TransactWriteItem{
				Delete: deleteOldEmail,
			},
		)
	}

	_, err = store.db.TransactWriteItems(txInp)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case dynamodb.ErrCodeTransactionCanceledException:
				return handy.ErrUserAlreadyExists
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return handy.ErrUserNotFound
			}
		}

		return err
	}

	return nil
}

func (store *UserDynamoDBStore) Delete(id string) error {
	user, err := store.Get(id)
	if err != nil {
		return err
	}

	idExistsExpr, _ := exists("id")

	deleteUser := &dynamodb.Delete{
		TableName: aws.String(usersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(user.ID)},
		},
		ConditionExpression:      idExistsExpr.Condition(),
		ExpressionAttributeNames: idExistsExpr.Names(),
	}

	deleteEmail := &dynamodb.Delete{
		TableName: aws.String(usersEmailsTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(user.Email)},
		},
	}

	txInp := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{Delete: deleteUser},
			{Delete: deleteEmail},
		},
	}

	_, err = store.db.TransactWriteItems(txInp)
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return handy.ErrUserNotFound
			}
		}

		return err
	}

	return nil
}

func exists(name string) (expression.Expression, error) {
	return expression.NewBuilder().WithCondition(expression.AttributeExists(expression.Name(name))).Build()
}

func notExists(name string) (expression.Expression, error) {
	return expression.NewBuilder().WithCondition(expression.AttributeNotExists(expression.Name(name))).Build()
}
