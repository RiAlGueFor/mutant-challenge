package mutantDNA

import (
  "testing"
  "strings"
  "strconv"
  // "os"
  // "github.com/aws/aws-lambda-go/events"
  // "github.com/aws/aws-lambda-go/lambda"
  // "github.com/aws/aws-sdk-go/aws"
  // "github.com/aws/aws-sdk-go/aws/session"
  // "github.com/aws/aws-sdk-go/service/dynamodb"
  // "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type TestDNA struct {
    in  string
    out bool
}

// const tableName = "LambaDNAValidationRecords"

var testDNA = []TestDNA{
  {"CAGTAC,TCAGGT,ACGGGG,TTACGG,AAAATG,CCCCAC",true},
  {"CTGCTA,CAGTAC,TCAGGT,AGCAGG,ATTCAG,TCACTC",true},
  {"ATGCTT,CTTTAC,TCTGGT,AGTAGG,TAGCGG,GTAGGC",false},
  {"ATGCGA,CAGTAC,TCAGGT,AGCAGG,ATACAG,TCACTG",true},
}

func TestScanningDNA(t *testing.T) {
  // region:=os.Getenv("AWS_REGION")
  // awsSession, err:=session.NewSession(&aws.Config{
  //   Region: aws.String(region)},)
  //
  // if err!=nil{
  //   t.Errorf("Error trying to connect to AWS ", err.Error())
  // }
  //
  // dynaClient = dynamodb.New(awsSession)

    var strSlice []string
    var isCorrect bool
    chainLetters := [3]string{ "A", "C", "G" }
    for i, test := range testDNA {
      strSlice = strings.Split(test.in, ",")
      for k := 0 ; k < len(chainLetters); k++ {
        isCorrect = ScanningDNA(strSlice,chainLetters[k])
        if isCorrect != test.out {
            t.Errorf("index %q, got %q, wanted %q", i,strconv.FormatBool(isCorrect), strconv.FormatBool(test.out))
            break;
        }
      }
    }
}
