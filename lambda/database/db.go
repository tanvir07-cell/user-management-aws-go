package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
      "photo_url":{
        S: aws.String(user.PhotoURL),
      },
    },
  })

   
  if err!=nil{
    return err
  }

  return nil

}


func (u DynamoDBClient) GetUser(username string)(types.User,error){
  var user types.User

	

	result,err:= u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username":{
				S: aws.String(username),
			},
		},
	})

	// aws error
	if err!=nil{
		return user,err
		 

	}

	if result.Item==nil{
		return user, fmt.Errorf("user not found")
	}

	// if user exist then umnarshal it from the json
	// like we did json.parse()
	err = dynamodbattribute.UnmarshalMap(result.Item,&user)

	if err!=nil{
		return user,err
	}
	return user,nil

}

