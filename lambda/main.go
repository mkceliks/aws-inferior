package main

import (
	"aws-inferior/config"
	"aws-inferior/entity"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var cognitoClient *cognitoidentityprovider.CognitoIdentityProvider

// Initialize the AWS Cognito client
func init() {
	// Create a new session using the default credentials and config from the environment
	sess := session.Must(session.NewSession())
	cognitoClient = cognitoidentityprovider.New(sess)
}

// signUpUser function to sign up a new user
func signUpUser(email, password string) error {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(config.ClientID),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}

	_, err := cognitoClient.SignUp(input)
	if err != nil {
		return fmt.Errorf("failed to sign up user: %v", err)
	}
	return nil
}

// signInUser function to authenticate a user
func signInUser(email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		ClientId: aws.String(config.ClientID),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(email),
			"PASSWORD": aws.String(password),
		},
	}

	result, err := cognitoClient.InitiateAuth(input)
	if err != nil {
		return nil, fmt.Errorf("failed to sign in user: %v", err)
	}
	return result, nil
}

// confirmSignUp function to confirm user sign-up
func confirmSignUp(email, code string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(config.ClientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
	}

	_, err := cognitoClient.ConfirmSignUp(input)
	if err != nil {
		return fmt.Errorf("failed to confirm sign-up: %v", err)
	}
	return nil
}

// SignUp handler function
func signUpUserHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var signUpReq entity.SignUpRequest
	if err := json.Unmarshal([]byte(req.Body), &signUpReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	// Sign up the user
	err := signUpUser(signUpReq.Email, signUpReq.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User signed up successfully!",
	}, nil
}

// SignIn handler function
func signInUserHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var signInReq entity.SignInRequest
	if err := json.Unmarshal([]byte(req.Body), &signInReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	resp, err := signInUser(signInReq.Email, signInReq.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}

	respBody, _ := json.Marshal(resp)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(respBody),
	}, nil
}

// ConfirmSignUp handler function
func confirmSignUpHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var confirmReq entity.ConfirmSignUpRequest
	if err := json.Unmarshal([]byte(req.Body), &confirmReq); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, nil
	}

	// Confirm sign-up
	err := confirmSignUp(confirmReq.Email, confirmReq.Code)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "User confirmed successfully!",
	}, nil
}

// Entry point for the Lambda function
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/signup":
		return signUpUserHandler(req)
	case "/signin":
		return signInUserHandler(req)
	case "/confirm-signup":
		return confirmSignUpHandler(req)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Not Found",
		}, nil
	}
}

func main() {
	// Start the Lambda function
	lambda.Start(handler)
}
