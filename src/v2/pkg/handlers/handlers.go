package handlers

import(
	"github.com/RiAlGueFor/mutant-challenge/src/v2/pkg/mutantDNA"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
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
  MutantDNA float32 `json:"count_mutant_dna, omitempty"`
  HumanDNA float32 `json:"count_human_dna, omitempty"`
  Ratio float32 `json:"ratio, omitempty"`
}

func GetStats(req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error,){
	mutantDNA.ConfigureDynamoDB()
	countMutantDNA, err:=mutantDNA.FetchDNARecords(true)
  if err!=nil {
    return apiResponse(http.StatusBadRequest,ErrorBody{
      aws.String(err.Error()),
    })
  }
	countNoMutantDNA, err:=mutantDNA.FetchDNARecords(false)
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

func CheckMutantDNA(req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse, error){
	result, err := mutantDNA.InitScanning(req)
	if result!=nil{
		if result.IsMutant{
      return apiResponse(http.StatusOK, nil)
    } else {
      return apiResponse(http.StatusForbidden, nil)
    }
	} else if err!=nil {
    if err.Error()!="" {
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
