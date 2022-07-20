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

var actualDNA DNARecord

func TestGetRecord(t *testing.T) {
  var expectedDNA string
  expectedDNA = "[\"ATGCGA\",\"CAGTAC\",\"TCAGGT\",\"AGCAGG\",\"TTACGG\",\"TCACTC\"]"
  expectKey := map[string]*dynamodb.AttributeValue{
		"dna": {
			S: aws.String(expectedDNA),
		},
	}

	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"dna": {
				S: aws.String(expectedDNA),
			},
		},
	}
  actualDNA.DNA = expectedDNA
  actualDNA.IsMutant = true

	//lets start dynamock in action
	mock.ExpectGetItem().ToTable("LambaDNAValidationRecords").WithKeys(expectKey).WillReturns(result)
  actualResult, _ := FetchDNARecord(expectedDNA)
  if (actualResult.DNA != actualDNA.DNA) && (actualResult.IsMutant != actualDNA.IsMutant){
		t.Errorf("Test Fail")
	}
}

func TestGetMutantRecords(t *testing.T) {
  var expectedCountMutantDNA float32
  expectedCountMutantDNA=2

  var expectedDNA string
  expectedDNA = "[\"ATGCGA\",\"CAGTAC\",\"TCAGGT\",\"AGCAGG\",\"TTACGG\",\"TCACTC\"]"

  result := dynamodb.ScanOutput{
    Items: []map[string]*dynamodb.AttributeValue{
      {
        "dna": {
          S: aws.String(expectedDNA),
        },
      },
      {
        "dna": {
          S: aws.String(expectedDNA),
        },
      },
    },
	}

  mock.ExpectScan().Table("LambaDNAValidationRecords").WillReturns(result)
  actualResultMutant, _ := FetchDNARecords(true)
  if actualResultMutant!=expectedCountMutantDNA {
		t.Errorf("Test Fail")
	}
}

func TestGetHumanRecords(t *testing.T) {
  var expectedCountHumanDNA float32
  expectedCountHumanDNA=4

  var expectedDNA string
  expectedDNA = "[\"ATGCGA\",\"CAGTAC\",\"TCAGGT\",\"AGCAGG\",\"TTACGG\",\"TCACTC\"]"

  result := dynamodb.ScanOutput{
    Items: []map[string]*dynamodb.AttributeValue{
      {
        "dna": {
          S: aws.String(expectedDNA),
        },
      },
      {
        "dna": {
          S: aws.String(expectedDNA),
        },
      },
      {
        "dna": {
          S: aws.String(expectedDNA),
        },
      },
      {
        "dna": {
          S: aws.String(expectedDNA),
        },
      },
    },
	}

  mock.ExpectScan().Table("LambaDNAValidationRecords").WillReturns(result)
  actualResultHuman, _ := FetchDNARecords(false)
  if actualResultHuman!=expectedCountHumanDNA {
		t.Errorf("Test Fail")
	}
}

func TestCreateRecord(t *testing.T) {
  var expectedDNA string
  expectedDNA = "[\"ATGCGA\",\"CAGTAC\",\"TCAGGT\",\"AGCAGG\",\"TTACGG\",\"TCACTC\"]"

	result := dynamodb.PutItemOutput{
		Attributes: map[string]*dynamodb.AttributeValue{
			"dna": {
				S: aws.String(expectedDNA),
			},
		},
	}
  actualDNA.DNA = expectedDNA
  actualDNA.IsMutant = true

  mock.ExpectPutItem().ToTable("LambaDNAValidationRecords").WillReturns(result)
  actualResult, _ := CreateRecordDNA(actualDNA)
  if (actualResult.DNA != actualDNA.DNA) && (actualResult.IsMutant != actualDNA.IsMutant){
		t.Errorf("Test Fail")
	}
}
