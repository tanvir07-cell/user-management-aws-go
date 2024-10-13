// api talks to the database and returns the response to the client
package api

import (
	"fmt"

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

func (a ApiHandler) CreateUser(user types.User) error{
	// check if user already exists
	exist,err:=a.db.IsUserExist(user.Username)
	if err!=nil{
		return fmt.Errorf("error checking if user exists: %v",err)
	}

	if exist{
		return fmt.Errorf("user already exists")
	}

	// create user
 createdUser,err:=types.NewUser(user)

 if err!=nil{
	 return fmt.Errorf("error creating user: %v",err)
 }

 err=a.db.CreateUser(createdUser)

 if err!=nil{
	 return fmt.Errorf("error creating user: %v",err)
 }

 return nil


	 


}



