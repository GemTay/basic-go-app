package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Drink struct {
	Id    int `json:",omitempty"`
	Name  string
	Price float64
}

const TABLE_NAME = "drinks"

func CreateTable() {
	_, err := Dyna.Db.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("N"),
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
			log.Println("Error:", awsErr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}
}

func PutItem(drink Drink) {
	_, err := Dyna.Db.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(len(GetAllItems()))),
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
			log.Println("Error:", awsErr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}
}

func UpdateItem(drink Drink) {
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
				N: aws.String(strconv.Itoa(drink.Id)),
			},
		},
		TableName:        aws.String(TABLE_NAME),
		UpdateExpression: aws.String("SET #N = :Name, #P = :Price"),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}
}

func GetItem(id int) (drink Drink) {
	result, err := Dyna.Db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Prints out full error message, including original error if there was one.
			log.Println("Error:", awsErr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &drink)
	if err != nil {
		fmt.Println(err.Error())
	}

	return drink
}

func GetAllItems() []Drink {
	result, err := Dyna.Db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(TABLE_NAME),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Println("Error:", awsErr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}

	drinks := []Drink{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &drinks)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return drinks
}