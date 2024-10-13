// api talks to the database and returns the response to the client
package api

import (
	"encoding/base64"
	"fmt"
	"strings"

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



func (a ApiHandler) CreateUser(user types.User) error {
	// Check if the user already exists
	exist, err := a.db.IsUserExist(user.Username)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	if exist {
		return fmt.Errorf("user already exists")
	}

	// Upload the user's photo to S3
	photoURL, err := UploadToS3(user.PhotoURL, user.Username)
	if err != nil {
		return fmt.Errorf("error uploading photo to S3: %v", err)
	}

	// Create a new user object with the S3 photo URL
	user.PhotoURL = photoURL
	createdUser, err := types.NewUser(user)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	// Store the user in the DynamoDB table
	err = a.db.CreateUser(createdUser)
	if err != nil {
		return fmt.Errorf("error creating user in DynamoDB: %v", err)
	}

	return nil
}
