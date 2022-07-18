package handlers

import(
	"github.com/RiAlGueFor/mutant-challenge/src/v2/pkg/mutantDNA"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
  ErrorMethodNotAllowed = "Method not allowed"
  ErrorFailedToUnmarshalRecord = "Failed to Unmarshall Record"
  ErrorInvalidDNAChain = "Invalid DNA Chain"
)

type ErrorBody struct{
  ErrorMsg *string `json:"error, omitempty"`
}

type Stats struct{
  MutantDNA int `json:"count_mutant_dna, omitempty"`
  HumanDNA int `json:"count_human_dna, omitempty"`
  Ratio float32 `json:"ratio, omitempty"`
}

func GetStats(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error,
){
  countMutantDNA, err:=mutantDNA.FetchDNARecords(tableName,dynaClient,true)
  if err!=nil {
    return apiResponse(http.StatusBadRequest,ErrorBody{
      aws.String(err.Error()),
    })
  }
  countNoMutantDNA, err:=mutantDNA.FetchDNARecords(tableName,dynaClient,false)
  if err!=nil {
    return apiResponse(http.StatusBadRequest,ErrorBody{
      aws.String(err.Error()),
    })
  }

  var result Stats
  result.MutantDNA=countMutantDNA
  result.HumanDNA=countNoMutantDNA
  if countNoMutantDNA>0{
    result.Ratio=float32(countMutantDNA/countNoMutantDNA)
  } else if countMutantDNA>0 {
    result.Ratio=1
  } else {
    result.Ratio=0
  }
  return apiResponse(http.StatusOK,result)
}

func CheckMutantDNA(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*events.APIGatewayProxyResponse, error){
  result, err := mutantDNA.InitScanning(req, tableName, dynaClient)
	if err!=nil {
    if err!="" {
      return apiResponse(http.StatusBadRequest, ErrorBody{
  			aws.String(err.Error()),
  		})
    } else {
      apiResponse(http.StatusBadRequest, nil)
    }
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod()(*events.APIGatewayProxyResponse, error){
  return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
