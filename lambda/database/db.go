package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"usermanagement.tanvirrifat.io/types"
)


const TABLE_NAME="userManagementTable"


type DynamoDBClient struct{
  databaseStore *dynamodb.DynamoDB
}

func NewDynamoDBClient() DynamoDBClient{
  return DynamoDBClient{
    databaseStore: dynamodb.New(session.Must(session.NewSession())),
  }
}


func (u DynamoDBClient) IsUserExist(username string)(bool,error){
   res,err:=u.databaseStore.GetItem(&dynamodb.GetItemInput{
    TableName: aws.String(TABLE_NAME),
    Key: map[string]*dynamodb.AttributeValue{
      "username":{
        S: aws.String(username),
      },
    },
   })

   // user exist 
   // but aws error
    if err!=nil{
      return true,err
    }

    // user  doesn't exist
    if res.Item==nil{
      return false,nil
    }

    // user exist

    return true,nil

}


func (u DynamoDBClient) CreateUser(user types.User) error{

  
  _,err:=u.databaseStore.PutItem(&dynamodb.PutItemInput{
    TableName: aws.String(TABLE_NAME),
    Item: map[string]*dynamodb.AttributeValue{
      "username":{
        S: aws.String(user.Username),
      },
      "email":{
        S: aws.String(user.Email),
      },
      "password":{
        S: aws.String(user.Password),

      },
    },
  })

   
  if err!=nil{
    return err
  }

  return nil

}

