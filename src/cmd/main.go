package main

import(
  "github.com/RiAlGueFor/mutant-challenge/src/pkg/handlers"
  "os"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
  dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
  region:=os.Getenv("AWS_REGION")
  awsSession, err:=session.NewSession(&aws.Config{
    Region: aws.String(region)},)

  if err!=nil{
    return
  }

  dynaClient = dynamodb.New(awsSession)
  lambda.Start(handler)

}

const tableName = "LambaDNAValidationRecords"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error){
    switch req.HTTPMethod {
      case "GET":
        return handlers.GetStats(req,tableName,dynaClient)
      case "POST":
        return handlers.CheckMutantDNA(req,tableName,dynaClient)
      default:
        return handlers.UnhandledMethod()
    }
}
