package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// Define the MentorshipAppStackProps struct, which will inherit from StackProps
type MentorshipAppStackProps struct {
	awscdk.StackProps
}

// Create Cognito User Pool
func CreateUserPool(stack awscdk.Stack) awscognito.UserPool {
	userPool := awscognito.NewUserPool(stack, jsii.String("inferior-userpool"), &awscognito.UserPoolProps{
		UserPoolName: jsii.String("inferior"),
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
		SelfSignUpEnabled: jsii.Bool(true),
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength:        jsii.Number(8),
			RequireSymbols:   jsii.Bool(false),
			RequireDigits:    jsii.Bool(true),
			RequireUppercase: jsii.Bool(true),
			RequireLowercase: jsii.Bool(true),
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("UserPoolId"), &awscdk.CfnOutputProps{
		Value:       userPool.UserPoolId(),
		Description: jsii.String("The ID of the Cognito User Pool"),
	})

	return userPool
}

// Create Cognito User Pool Client
func CreateUserPoolClient(stack awscdk.Stack, userPool awscognito.UserPool) {
	userPoolClient := awscognito.NewUserPoolClient(stack, jsii.String("inferior-client"), &awscognito.UserPoolClientProps{
		UserPool:       userPool,
		GenerateSecret: jsii.Bool(false),
	})

	awscdk.NewCfnOutput(stack, jsii.String("UserPoolClientId"), &awscdk.CfnOutputProps{
		Value:       userPoolClient.UserPoolClientId(),
		Description: jsii.String("The ID of the Cognito User Pool Client"),
	})
}

// Function to create a Lambda function from Docker image
func CreateDockerLambdaFunction(stack awscdk.Stack) awslambda.Function {
	myLambda := awslambda.NewFunction(stack, jsii.String("MyDockerLambdaFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2(),
		Code:    awslambda.Code_FromAssetImage(jsii.String("034362052544.dkr.ecr.us-east-1.amazonaws.com/auth-lambda"), &awslambda.AssetImageCodeProps{}),
		Handler: jsii.String("main"),
	})

	awscdk.NewCfnOutput(stack, jsii.String("LambdaFunctionArn"), &awscdk.CfnOutputProps{
		Value:       myLambda.FunctionArn(),
		Description: jsii.String("The ARN of the Docker Lambda Function"),
	})

	return myLambda
}

// Function to create an API Gateway
func CreateApiGateway(stack awscdk.Stack, myLambda awslambda.Function) {
	api := awsapigateway.NewLambdaRestApi(stack, jsii.String("MyApiGateway"), &awsapigateway.LambdaRestApiProps{
		Handler: myLambda,
	})

	awscdk.NewCfnOutput(stack, jsii.String("ApiGatewayUrl"), &awscdk.CfnOutputProps{
		Value:       api.Url(),
		Description: jsii.String("The URL of the API Gateway"),
	})
}

// Main stack creation function
func NewMentorshipAppStack(scope constructs.Construct, id string, props *MentorshipAppStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	userPool := CreateUserPool(stack)
	CreateUserPoolClient(stack, userPool)

	myLambda := CreateDockerLambdaFunction(stack)
	CreateApiGateway(stack, myLambda)

	return stack
}

// Main function for CDK App
func main() {
	app := awscdk.NewApp(nil)

	NewMentorshipAppStack(app, "inferior-stack", &MentorshipAppStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: jsii.String("034362052544"),
				Region:  jsii.String("us-east-1"),
			},
		},
	})

	app.Synth(nil)
}
