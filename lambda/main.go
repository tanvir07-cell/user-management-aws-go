package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"usermanagement.tanvirrifat.io/app"
	"usermanagement.tanvirrifat.io/types"
)

// type User struct{
// 	ID string
// 	Username string
// 	Email string
// 	Password string

// }

func HandleUser(u types.User){
	fmt.Println(u)
	user,err:=types.NewUser(u)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(user)
	
}


func main(){

	app:= app.NewApp()

	lambda.Start(app.ApiHandler.CreateUser)

	
}