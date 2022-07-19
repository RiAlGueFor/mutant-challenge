package mutantDNA

import(
  "github.com/RiAlGueFor/mutant-challenge/src/v2/pkg/validators"
  "encoding/json"
  "errors"
  "strings"
	"github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
  ErrorFailedToUnmarshalRecord = "Failed to Unmarshall Record"
  ErrorInvalidDNAChain = "Invalid DNA Chain"
  ErrorFailedToFetchRecord = "Failed To Fetch Record"
  ErrorCouldNotMarshalItem = "Could not Marshal Item"
  ErrorCouldNotDynamoPutItem = "Could not Dynamo Put Item"
)

type DNAChain struct {
  DNA []string `json:"dna, omitempty"`
}
type DNARecord struct {
  DNA string `json:"dna, omitempty"`
  IsMutant bool `json:"isMutant, omitempty"`
}

func InitScanning(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*DNARecord, error){
// func InitScanning(req events.APIGatewayProxyRequest)(*DNARecord, error){
  var dnaChain DNAChain
  if err := json.Unmarshal([]byte(req.Body), &dnaChain); err!=nil {
    return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
  // 1 - Check if DNA is valid
  var isValid=validators.IsDNAValid(dnaChain.DNA)
  if !isValid {
    return nil, errors.New(ErrorInvalidDNAChain)
  }
  // 2 - Check if DNA was previously validated. If it was, return the isMutant flag
  var dnaRecord DNARecord
  dnaJoin := strings.Join(dnaChain.DNA, "-")
  dnaJoin = strings.Replace(dnaJoin,"-","\",\"",-1)
  dnaRecord.DNA = "[\""+ dnaJoin +"\"]"

  ConfigureDynamoDB()
  currentDNA, _:=FetchDNARecord(dnaRecord.DNA,tableName,dynaClient)
  // currentDNA, _:=FetchDNARecord(dnaRecord.DNA)
  if currentDNA!=nil && len(currentDNA.DNA)>0 {
    if currentDNA.IsMutant{
      return &dnaRecord, nil
    } else {
      return nil, errors.New("")
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
  // insertItem, err:= CreateRecordDNA(dnaRecord)
  if err!=nil{
    return nil, err
  }

  if insertItem.IsMutant {
    return &dnaRecord, nil
  } else {
    return nil, errors.New("")
  }
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
