/*
	Database layer for the url shortner. The GoLinkShortner API uses a redis database as a backend.
	It follows a common connection pattern where a main session is created then other sessions are created by copying the information of the main session while utilizing a different socket from a socket pool
*/
package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"gopkg.in/redis.v5"
	"log"
)

var ErrDuplicate = errors.New("duplicate slug")
var ErrNotFound = errors.New("slug not found")
var ErrTableNotFound = errors.New("database table not found")
var ErrConnection = errors.New("connection issue")
var ErrOverload = errors.New("overloaded database")
var ErrGeneric = errors.New("generic error")

type Database struct {
	Redis  *redis.Client
	Dynamo *dynamodb.DynamoDB
	Table  string
}

type Item struct {
	slug        string
	destination string
}

func NewDatabase() *Database {
	redis := NewRedis()
	dynamo := NewDynamoDb()

	return &Database{
		Redis:  redis,
		Dynamo: dynamo,
		Table:  *dynamoTable,
	}
}

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: *redisEndpoint,
	})
}

func NewDynamoDb() *dynamodb.DynamoDB {

	session := NewSessionAWS(*dynamoRegion, *dynamoEndpoint)
	return dynamodb.New(session)
}

func (db *Database) Get(slug string) (destination string, err error) {

	destination, err = db.GetRedis(slug)

	// return cached url
	if err == nil {
		log.Printf("Returning cached %s", destination)
		return destination, nil
	}
	log.Printf("Cache miss for %s with error %v", slug, err)

	// retrieve from database
	destination, err = db.GetDynamoDb(slug)

	if err == nil {
		log.Printf("Caching and returning source %s", destination)
		err = db.PutRedis(slug, destination)
		return destination, err
	}

	log.Printf("Source miss for %s with error  %v ", slug)
	return destination, err
}

func (db *Database) Put(slug string, destination string) (err error) {

	err = db.PutDynamoDb(slug, destination)

	if err != nil {
		return err
	}

	log.Printf("Inserted to source %s", destination)
	return nil
}

func (db *Database) GetRedis(slug string) (long_url string, err error) {

	val, err := db.Redis.Get(slug).Result()

	if err == redis.Nil {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	} else {
		return val, err
	}

}

func (db *Database) PutRedis(slug string, destination string) (err error) {

	val, err := db.Redis.Exists(slug).Result()

	if err != nil {

		log.Printf("RedisPut error: %v", err)
		return ErrGeneric

	} else if val == false {

		errSet := db.Redis.Set(slug, destination, 0).Err()

		if errSet != nil {
			log.Println(errSet)
		}

		log.Printf("Saving short:%s long:%s", slug, destination)

	} else {

		return ErrDuplicate

	}

	return

}

func (db *Database) GetDynamoDb(slug string) (destination string, err error) {

	input := BuildDynamoGetItemInput(slug, db.Table)

	result, err := db.Dynamo.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Printf("DynamoGet error: %v %v", result, aerr)
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return "", ErrOverload
			case dynamodb.ErrCodeResourceNotFoundException:
				return "", ErrTableNotFound // Initialize db?
			case dynamodb.ErrCodeInternalServerError:
			default:
				return "", ErrGeneric
			}
		}
		return "", err
	}

	var item Item
	if err = dynamodbattribute.Unmarshal(result.Item["destination"], &item.destination); err != nil {
		log.Printf("UnmarshalMap(GetItem response) err=%q", err)
	}

	return item.destination, nil
}

func (db *Database) PutDynamoDb(slug string, destination string) (err error) {

	input := BuildDynamoPutItemInput(slug, destination, db.Table)

	result, err := db.Dynamo.PutItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {

			log.Printf("DynamoPut error: %v %v", result, aerr)

			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return ErrDuplicate
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return ErrOverload
			case dynamodb.ErrCodeResourceNotFoundException:
				return ErrTableNotFound
			case dynamodb.ErrCodeInternalServerError:
			default:
				return ErrGeneric
			}
		}
		return err
	}

	return nil

}

func NewSessionAWS(region string, endpoint string) *session.Session {

	return session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
		//CredentialsChainVerboseErrors: aws.Bool(true),
		//Endpoint: aws.String("https://dynamodb.eu-west-1.amazonaws.com"),
		//LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	}))

}

func BuildDynamoPutItemInput(slug string, destination string, table string) *dynamodb.PutItemInput {

	return &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"slug": {
				S: aws.String(slug),
			},
			"destination": {
				S: aws.String(destination),
			},
		},
		TableName:           aws.String(table),
		ConditionExpression: aws.String("attribute_not_exists(slug)"),
	}
}

func BuildDynamoGetItemInput(slug string, table string) *dynamodb.GetItemInput {

	return &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"slug": {
				S: aws.String(slug),
			},
		},
		TableName: aws.String(table),
	}

}
