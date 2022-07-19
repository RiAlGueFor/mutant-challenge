package mutantDNA

import (
  dynamock "github.com/gusaul/go-dynamock"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var mock *dynamock.DynaMock

func init() {
	Dyna = new(MyDynamo)
	Dyna.Db, mock = dynamock.New()
}

func TestGetRecords(t *testing.T) {

  expectKey := map[string]*dynamodb.AttributeValue{
		"dna": {
			N: aws.String("[\"ATGCGA\",\"CAGTAC\",\"TCAGGT\",\"AGCAGG\",\"TTACGG\",\"TCACTC\"]"),
		},
	}

	expectedResult := aws.String("true")
	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"isMutant": {
				S: expectedResult,
			},
		},
	}

	//lets start dynamock in action
	mock.ExpectGetItem().ToTable("LambaDNAValidationRecords").WithKeys(expectKey).WillReturns(result)

  FetchDNARecords("LambaDNAValidationRecords",dynaClient,true)
	actualResult, _ := GetName("1")
	if actualResult != expectedResult {
		t.Errorf("Test Fail")
	}
}
