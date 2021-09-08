package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pborman/uuid"
)

type Drink struct {
	Id    string `json:"Id"`
	Name  string
	Price float64
}

const TABLE_NAME = "drinks"

func (app *application) CreateTable() {
	_, err := Dyna.Db.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Prints out full error message, including original error if there was one.
			app.errorLog.Println(awsErr.Error())
		} else {
			app.errorLog.Println(err.Error())
		}
	}
}

func (app *application) PutItem(drink Drink) {
	uId := uuid.New()
	_, err := Dyna.Db.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(uId),
			},
			"Name": {
				S: aws.String(drink.Name),
			},
			"Price": {
				N: aws.String(strconv.FormatFloat(drink.Price, 'f', 2, 64)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			app.errorLog.Println(awsErr.Error())
		} else {
			app.errorLog.Println(err.Error())
		}
	}
}

func (app *application) UpdateItem(drink Drink) {
	_, err := Dyna.Db.UpdateItem(&dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("Name"),
			"#P": aws.String("Price"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":Name": {
				S: aws.String(drink.Name),
			},
			":Price": {
				N: aws.String(strconv.FormatFloat(drink.Price, 'f', 2, 64)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(drink.Id),
			},
		},
		TableName:        aws.String(TABLE_NAME),
		UpdateExpression: aws.String("SET #N = :Name, #P = :Price"),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			app.errorLog.Println(awsErr.Error())
		} else {
			app.errorLog.Println(err.Error())
		}
	}
}

func (app *application) GetItem(id string) (drink Drink) {
	result, err := Dyna.Db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Prints out full error message, including original error if there was one.
			app.errorLog.Println(awsErr.Error())
		} else {
			app.errorLog.Println(err.Error())
		}
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &drink)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	return drink
}

func (app *application) GetAllItems() []Drink {
	result, err := Dyna.Db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			app.errorLog.Println(awsErr.Error())
		} else {
			app.errorLog.Println(err.Error())
		}
	}

	drinks := []Drink{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &drinks)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	return drinks
}
