package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
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

	// lambda.Start(app.ApiHandler.CreateUser)

	lambda.Start(func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error){
		switch request.Path{
		case "/register":
			return app.ApiHandler.CreateUser(request)

			case "/login":
			return app.ApiHandler.LoginUser(request)
		
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body: "Not Found",
			},nil
		}
		
	})

	
}