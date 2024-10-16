package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type UserManagementStackProps struct {
	awscdk.StackProps
}

func NewUserManagementStack(scope constructs.Construct, id string, props *UserManagementStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
// aws lambda code :
	 myLambda:=awslambda.NewFunction(stack,jsii.String("userManagementLambda"),
		&awslambda.FunctionProps{
			Runtime: awslambda.Runtime_PROVIDED_AL2023(),
			Code: awslambda.Code_FromAsset(jsii.String("lambda/function.zip"),nil),
			Handler: jsii.String("main"),
			
		})


// aws dynamodb code :

table:=awsdynamodb.NewTable(stack,jsii.String("userManagementTable"),&awsdynamodb.TableProps{
	PartitionKey: &awsdynamodb.Attribute{
		 Name: jsii.String("username"),
		 Type: awsdynamodb.AttributeType_STRING,
	},
	TableName: jsii.String("userManagementTable"),
})

table.GrantReadWriteData(myLambda)


// bucket s3

myBucket:=awss3.NewBucket(stack,jsii.String("userManagementBucket"),&awss3.BucketProps{
	Versioned: jsii.Bool(true),
	BucketName: jsii.String("usermanagementbucketupload"),

})

 myBucket.GrantReadWrite(myLambda,nil)


// api gateway:
api:=awsapigateway.NewRestApi(stack,jsii.String("userManagementApi"),&awsapigateway.RestApiProps{
	DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type","Authorization"),
			AllowMethods: jsii.Strings("GET","POST","DELETE","PUT","OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
	},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},

})

// integration apigateway to the lambda function
	integration:= awsapigateway.NewLambdaIntegration(myLambda,nil)

	// register route
	registerRoute:= api.Root().AddResource(jsii.String("register"),nil)
	registerRoute.AddMethod(jsii.String("POST"),integration,nil)

	// login route

		loginRoute:= api.Root().AddResource(jsii.String("login"),nil)
	loginRoute.AddMethod(jsii.String("POST"),integration,nil)


	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewUserManagementStack(app, "UserManagementStack", &UserManagementStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
