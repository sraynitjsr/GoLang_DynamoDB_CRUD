package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}

	client := dynamodb.NewFromConfig(cfg)

	tableName := "MyTable"

	err = putItem(client, tableName)
	if err != nil {
		fmt.Println("Error putting item:", err)
		return
	}

	item, err := getItem(client, tableName)
	if err != nil {
		fmt.Println("Error getting item:", err)
		return
	}
	fmt.Println("Item Retrieved:", item)
}

func putItem(client *dynamodb.Client, tableName string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: "1"},
			"Name": &types.AttributeValueMemberS{Value: "Virat Kohli"},
		},
	}

	_, err := client.PutItem(context.TODO(), input)
	return err
}

func getItem(client *dynamodb.Client, tableName string) (map[string]types.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: "1"},
		},
	}

	output, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return output.Item, nil
}
