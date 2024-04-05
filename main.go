package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/gin-gonic/gin"
    "net/http"
)

type Item struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var svc *dynamodb.DynamoDB

func main() {
    sess, err := session.NewSession(&aws.Config{
        Endpoint: aws.String("http://localhost:8000"),
    })
    if err != nil {
        panic(err)
    }

    svc = dynamodb.New(sess)

    router := gin.Default()

    router.POST("/items", createItemHandler)
    router.GET("/items/:id", getItemHandler)
    router.PUT("/items/:id", updateItemHandler)
    router.DELETE("/items/:id", deleteItemHandler)

    router.Run(":8080")
}

func createItemHandler(c *gin.Context) {
    var item Item
    if err := c.BindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := putItem(item)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, item)
}

func getItemHandler(c *gin.Context) {
    id := c.Param("id")

    item, err := getItem(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if item == nil {
        c.Status(http.StatusNotFound)
        return
    }

    c.JSON(http.StatusOK, item)
}

func updateItemHandler(c *gin.Context) {
    id := c.Param("id")

    var item Item
    if err := c.BindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    item.ID = id
    err := putItem(item)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, item)
}

func deleteItemHandler(c *gin.Context) {
    id := c.Param("id")

    err := deleteItem(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func putItem(item Item) error {
    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        return err
    }

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String("my-table-name"),
    }

    _, err = svc.PutItem(input)
    return err
}

func getItem(id string) (*Item, error) {
    input := &dynamodb.GetItemInput{
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                S: aws.String(id),
            },
        },
        TableName: aws.String("my-table-name"),
    }

    result, err := svc.GetItem(input)
    if err != nil {
        return nil, err
    }

    if len(result.Item) == 0 {
        return nil, nil
    }

    var item Item
    err = dynamodbattribute.UnmarshalMap(result.Item, &item)
    if err != nil {
        return nil, err
    }

    return &item, nil
}

func deleteItem(id string) error {
    input := &dynamodb.DeleteItemInput{
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                S: aws.String(id),
            },
        },
        TableName: aws.String("my-table-name"),
    }

    _, err := svc.DeleteItem(input)
    return err
}
