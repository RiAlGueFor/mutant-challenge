package handlers

import(
	"github.com/RiAlGueFor/mutant-challenge/src/pkg/mutantDNA"
  "github.com/RiAlGueFor/mutant-challenge/src/pkg/dynamoDB"
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

type DNAChain struct {
  DNA []string `json:"dna, omitempty"`
}

type Stats struct{
  MutantDNA int `json:"count_mutant_dna, omitempty"`
  HumanDNA int `json:"count_human_dna, omitempty"`
  Ratio float32 `json:"ratio, omitempty"`
}

type DNARecord struct {
  DNA string `json:"dna, omitempty"`
  IsMutant bool `json:"isMutant, omitempty"`
}

func GetStats(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error,
){
  countMutantDNA, err:=dynamoDB_API.FetchDNARecords(tableName,dynaClient,true)
  if err!=nil {
    return apiResponse(http.StatusBadRequest,ErrorBody{
      aws.String(err.Error()),
    })
  }
  countNoMutantDNA, err:=dynamoDB_API.FetchDNARecords(tableName,dynaClient,false)
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
  var dnaChain DNAChain
  if err := json.Unmarshal([]byte(req.Body), &dnaChain); err!=nil {
    return apiResponse(http.StatusForbidden, ErrorBody{
      aws.String(ErrorFailedToUnmarshalRecord),
    })
	}
  // 1 - Check if DNA is valid
  var isValid=validators.IsDNAValid(dnaChain.DNA)
  if !isValid {
    return apiResponse(http.StatusForbidden, ErrorBody{
      aws.String(ErrorInvalidDNAChain),
    })
  }
  // 2 - Check if DNA was previously validated. If it was, return the isMutant flag
  var dnaRecord DNARecord
  dnaJoin := strings.Join(dnaChain.DNA, "-")
  dnaJoin = strings.Replace(dnaJoin,"-","\",\"",-1)
  dnaRecord.DNA = "[\""+ dnaJoin +"\"]"
  currentDNA, _:=dynamoDB_API.FetchDNARecord(dnaRecord.DNA,tableName,dynaClient)
  if currentDNA!=nil && len(currentDNA.DNA)>0 {
    if currentDNA.IsMutant{
      return apiResponse(http.StatusOK, nil)
    } else {
      return apiResponse(http.StatusForbidden, nil)
    }
  }
  // 3 - If it wasn't validated, go on with the validation
  dnaRecord.IsMutant = true
  chainLetters := [3]string{ "A", "C", "G" }
  for k := 0 ; k < len(chainLetters); k++ {
    if (!ScanningDNA(dnaChain.DNA,chainLetters[k])) {
      dnaRecord.IsMutant=false
      break;
    }
  }
  // 4 - After Checking the DNA, Store DNA Chain and Validation Result on DynamoDB
  _, err:= dynamoDB_API.CreateRecordDNA(dnaRecord,tableName,dynaClient)
  if err!=nil{
    return apiResponse(http.StatusBadRequest, ErrorBody{
      aws.String(err.Error()),
    })
  }

  if dnaRecord.IsMutant {
    return apiResponse(http.StatusOK, nil)
  } else {
    return apiResponse(http.StatusForbidden, nil)
  }
}

func UnhandledMethod()(*events.APIGatewayProxyResponse, error){
  return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
