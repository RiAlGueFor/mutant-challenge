package main

import(
  "github.com/RiAlGueFor/mutant-challenge/src/v2/pkg/handlers"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
)

func main() {
  lambda.Start(handler)

}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error){
    switch req.HTTPMethod {
      case "GET":
        return handlers.GetStats(req)
      case "POST":
        return handlers.CheckMutantDNA(req)
      default:
        return handlers.UnhandledMethod()
    }
}
