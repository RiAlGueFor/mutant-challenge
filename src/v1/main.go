package main

import (
  "encoding/json"
  "net/http"
  "regexp"
  "strings"
  "errors"
  "os"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
  dynaClient dynamodbiface.DynamoDBAPI
  ErrorFailedToFetchRecord = "Failed To Fetch Record"
  ErrorMethodNotAllowed = "Method not allowed"
  ErrorFailedToUnmarshalRecord = "Failed to Unmarshall Record"
  ErrorInvalidDNAChain = "Invalid DNA Chain"
  ErrorResourceNotFound = "Resource Not Found"
  ErrorCouldNotMarshalItem = "Could not Marshal Item"
  ErrorCouldNotDynamoPutItem = "Could not Dynamo Put Item"
  ErrorDNAValidationAlreadyExists= "DNA Validation Already Exists"
)

type DNAChain struct {
  DNA []string `json:"dna, omitempty"`
}

type DNARecord struct {
  DNA string `json:"dna, omitempty"`
  IsMutant bool `json:"isMutant, omitempty"`
}

type Stats struct{
  MutantDNA float32 `json:"count_mutant_dna, omitempty"`
  HumanDNA float32 `json:"count_human_dna, omitempty"`
  Ratio float32 `json:"ratio, omitempty"`
}

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
        return GetStats(req,tableName,dynaClient)
      case "POST":
        return CheckMutantDNA(req,tableName,dynaClient)
      default:
        return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
    }
}

func apiResponse(status int, body interface{})(*events.APIGatewayProxyResponse, error){
  resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type":"application/json"}}
	resp.StatusCode = status
  if body!=nil{
    stringBody, _ := json.Marshal(body)
    resp.Body = string(stringBody)
  }
  return &resp,nil
}

func ScanningDNA(dna []string, letter2Check string) bool {
  oCounter:=0
  vCounter:=0
  aux_i:=0
  aux_j:=0

  repString := strings.Repeat(letter2Check, 4)

  for i := 0; i < len(dna) ; i++{
    // Horizontal ScanningDNA
    if(strings.Contains(dna[i], repString)){
      return true
    }

    oCounter = 0
    vCounter = 0
    for j := 0; j < len(dna[i]) ; j++{
      if (string(dna[i][j]) == letter2Check){
        oCounter++
        vCounter++

        // Vertical ScanningDNA
        if(len(dna)-i>=4){
          aux_i=i
          for vCounter < 4 {
            aux_i++
            if (string(dna[aux_i][j]) == letter2Check){
              vCounter++
            } else {
              break;
            }
          }
          if(vCounter == 4){
            return true
          } else {
            vCounter = 0
          }
        }
        // Oblique ScanningDNA
        if(len(dna[i])-j>=4){
          // Go on

          // Oblique ScanningDNA
          if(len(dna)-i>=4){
            // Downside
            aux_i = i
            aux_j = j
            for oCounter < 4 {
          		aux_i++
              aux_j++
              if (string(dna[aux_i][aux_j]) == letter2Check){
                oCounter++
              } else {
                oCounter=0
                break;
              }
          	}

            if(oCounter == 4){
              return true
            } else {
              oCounter=0
            }
          }

          if(i>=3){
            // Upside
            aux_i=i
            aux_j=j
            for oCounter < 4 {
          		aux_i--
              aux_j++
              if (string(dna[aux_i][aux_j]) == letter2Check){
                oCounter++
              } else {
                oCounter=0
                break;
              }
          	}

            if(oCounter == 4){
              return true
            } else {
              oCounter=0
            }
          }
        } else {
          oCounter=0
          break;
        }
      }
    }
  }
  return false
}

type ErrorBody struct{
  ErrorMsg *string `json:"error, omitempty"`
}

func GetStats(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse, error,
){
  countMutantDNA, err:=FetchDNARecords(tableName,dynaClient,true)
  if err!=nil {
    return apiResponse(http.StatusBadRequest,ErrorBody{
      aws.String(err.Error()),
    })
  }
  countNoMutantDNA, err:=FetchDNARecords(tableName,dynaClient,false)
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
  var isValid=IsDNAValid(dnaChain.DNA)
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
  currentDNA, _:=FetchDNARecord(dnaRecord.DNA,tableName,dynaClient)
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
  _, err:= CreateRecordDNA(dnaRecord,tableName,dynaClient)
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

func FetchDNARecord(dnaString string, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*DNARecord, error){
  input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"dna":{
				S: aws.String(dnaString),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err!= nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(DNARecord)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchDNARecords(tableName string, dynaClient dynamodbiface.DynamoDBAPI, isMutant bool)(float32, error) {
  filt := expression.Name("isMutant").Equal(expression.Value(isMutant))
  expr, err := expression.NewBuilder().WithFilter(filt).Build()

  if err != nil {
      return 0, errors.New(ErrorFailedToFetchRecord)
  }

  params := &dynamodb.ScanInput{
      ExpressionAttributeNames:  expr.Names(),
      ExpressionAttributeValues: expr.Values(),
      FilterExpression:          expr.Filter(),
      ProjectionExpression:      expr.Projection(),
      TableName:                 aws.String(tableName),
  }

  result, err := dynaClient.Scan(params)
  if err!=nil {
    return 0, errors.New(ErrorFailedToFetchRecord)
  }
  return float32(len(result.Items)), nil
}

func CreateRecordDNA(dnaRecord DNARecord, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*DNARecord, error){

  av, err := dynamodbattribute.MarshalMap(dnaRecord)

  if err!=nil{
    return nil, errors.New(ErrorCouldNotMarshalItem)
  }

  input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}

  _, err = dynaClient.PutItem(input)
  if err!=nil{
    return nil, errors.New(ErrorCouldNotDynamoPutItem)
  }
  return &dnaRecord, nil
}

func IsDNAValid(dna []string) bool{
  match := regexp.MustCompile(`[^ATCG]`)

  var initialLength int = len(dna[0])

	for k := 0 ; k < len(dna); k++ {
    if(len(dna[k])!=initialLength) {
      return false
    }
    if(match.MatchString(dna[k])) {
      return false
    }
	}
  return true
}
