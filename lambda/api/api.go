// api talks to the database and returns the response to the client
package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	db "usermanagement.tanvirrifat.io/database"
	"usermanagement.tanvirrifat.io/types"
)


type ApiHandler struct{
	db db.DynamoDBClient


}

func NewApiHandler() ApiHandler{
	return ApiHandler{
		db: db.NewDynamoDBClient(),

	
	}
}


func UploadToS3(photoBase64,username string)(string,error){
	   photoData,err:=base64.StdEncoding.DecodeString(photoBase64)

		 if err!=nil{
			 return "",fmt.Errorf("error decoding base64 string: %v",err)
		 }

	  sess:=session.Must(session.NewSession())
		svc:=s3.New(sess)

		// Upload to S3
	photoKey := "photos/" + username + ".jpg" 
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String("usermanagementbucketupload"),
		Key:           aws.String(photoKey),
		Body:          strings.NewReader(string(photoData)),
		ContentLength: aws.Int64(int64(len(photoData))),
		ContentType:   aws.String("image/jpeg"),
	})

	if err != nil {
		return "", err
	}

	// Generate S3 URL
	s3URL := "https://usermanagementbucketupload.s3.amazonaws.com/" + photoKey
	return s3URL, nil
}




// func (a ApiHandler) CreateUser(user types.User) error{
// 	// check if user already exists
// 	exist,err:=a.db.IsUserExist(user.Username)
// 	if err!=nil{
// 		return fmt.Errorf("error checking if user exists: %v",err)
// 	}

// 	if exist{
// 		return fmt.Errorf("user already exists")
// 	}

// 	// create user
//  createdUser,err:=types.NewUser(user)

//  if err!=nil{
// 	 return fmt.Errorf("error creating user: %v",err)
//  }

//  err=a.db.CreateUser(createdUser)

//  if err!=nil{
// 	 return fmt.Errorf("error creating user: %v",err)
//  }

//  return nil


	 


// }



func (a ApiHandler) CreateUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse,error) {
	// Check if the user already exists

	var registerUser types.User

	err:=json.Unmarshal([]byte(request.Body),&registerUser)

	if err!=nil{
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body: "Bad Request",
		},err
	}

	if registerUser.Email=="" || registerUser.Password=="" || registerUser.Username=="" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body: "Username or Email or Password is missing",
		},err
	}



	exist, err := a.db.IsUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		}, err

		
	}

	if exist {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusConflict,
			Body:       "User already exists",
		},err
	}

	// Upload the user's photo to S3
	photoURL, err := UploadToS3(registerUser.PhotoURL, registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		}, err
	
	}

	// Create a new user object with the S3 photo URL
	registerUser.PhotoURL = photoURL
	createdUser, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		}, err
		
	}

	// Store the user in the DynamoDB table
	err = a.db.CreateUser(createdUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User created successfully",
	}, nil
	
}
