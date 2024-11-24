package main

import (
    "encoding/json"
    "fmt"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a transaction
type SmartContract struct {
    contractapi.Contract
}

// Transaction represents a transaction with three sources
type Transaction struct {
    ID       string `json:"id"`
    Source1  string `json:"source1"`
    Source2  string `json:"source2"`
    Source3  string `json:"source3"`
    Processed bool  `json:"processed"`
}

// CreateTransaction creates a new transaction
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, id string, source1 string, source2 string, source3 string) error {
    transaction := Transaction{
        ID:       id,
        Source1:  source1,
        Source2:  source2,
        Source3:  source3,
        Processed: false,
    }

    transactionJSON, err := json.Marshal(transaction)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, transactionJSON)
}

// ValidateTransaction validates and processes the transaction if at least two sources are equal
func (s *SmartContract) ValidateTransaction(ctx contractapi.TransactionContextInterface, id string) error {
    transactionJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return fmt.Errorf("failed to read from world state: %v", err)
    }
    if transactionJSON == nil {
        return fmt.Errorf("transaction %s does not exist", id)
    }

    var transaction Transaction
    err = json.Unmarshal(transactionJSON, &transaction)
    if err != nil {
        return err
    }

    if (transaction.Source1 == transaction.Source2) || (transaction.Source1 == transaction.Source3) || (transaction.Source2 == transaction.Source3) {
        transaction.Processed = true
    } else {
        return fmt.Errorf("transaction %s could not be validated", id)
    }

    transactionJSON, err = json.Marshal(transaction)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, transactionJSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating SmartContract chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting SmartContract chaincode: %s", err.Error())
    }
}
