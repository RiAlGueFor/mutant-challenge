package mutantDNA

import(
  "errors"
  "os"
	"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
  "github.com/aws/aws-sdk-go/aws/session"
)
// MyDynamo struct hold dynamodb connection
type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

// Dyna - object from MyDynamo
var Dyna *MyDynamo

const tableName = "LambaDNAValidationRecords"

// ConfigureDynamoDB - init func for open connection to aws dynamodb
func ConfigureDynamoDB() error {
	Dyna = new(MyDynamo)
  region:=os.Getenv("AWS_REGION")
  awsSession, err:=session.NewSession(&aws.Config{
    Region: aws.String(region)},)

  if err!=nil{
    return errors.New(ErrorCouldNotConnectToAWS)
  }

	svc := dynamodb.New(awsSession)
	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
  return nil
}

// func FetchDNARecord(dnaString string, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*DNARecord, error){
func FetchDNARecord(dnaString string)(*DNARecord, error){
  input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"dna":{
				S: aws.String(dnaString),
			},
		},
		TableName: aws.String(tableName),
	}

  result, err := Dyna.Db.GetItem(input)
  // result, err := dynaClient.GetItem(input)
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

// func FetchDNARecords(tableName string, dynaClient dynamodbiface.DynamoDBAPI, isMutant bool)(float32, error) {
func FetchDNARecords(isMutant bool)(float32, error) {
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
  result, err := Dyna.Db.Scan(params)
  // result, err := dynaClient.Scan(params)
  if err!=nil {
    return 0, errors.New(ErrorFailedToFetchRecord)
  }
  return float32(len(result.Items)), nil
}

// func CreateRecordDNA(dnaRecord DNARecord, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*DNARecord, error){
func CreateRecordDNA(dnaRecord DNARecord)(*DNARecord, error){
  av, err := dynamodbattribute.MarshalMap(dnaRecord)

  if err!=nil{
    return nil, errors.New(ErrorCouldNotMarshalItem)
  }

  input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(tableName),
	}
  _, err = Dyna.Db.PutItem(input)
  // _, err = dynaClient.PutItem(input)
  if err!=nil{
    return nil, errors.New(ErrorCouldNotDynamoPutItem)
  }
  return &dnaRecord, nil
}
