package dynamoDB_API

import (
  "github.com/RiAlGueFor/mutant-challenge/src/pkg/mutantDNA"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
)


var (
  ErrorFailedToFetchRecord = "Failed To Fetch Record"
  ErrorFailedToUnmarshalRecord = "Failed to Unmarshall Record"
  ErrorCouldNotMarshalItem = "Could not Marshal Item"
  ErrorCouldNotDynamoPutItem = "Could not Dynamo Put Item"
)

func FetchDNARecord(dnaString string, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*mutantDNA.DNARecord, error){
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

	item := new(mutantDNA.DNARecord)
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

func CreateRecordDNA(dnaRecord mutantDNA.DNARecord, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*mutantDNA.DNARecord, error){
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
