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

// Transaction represents a transaction with multiple sources
type Transaction struct {
    ID        string   `json:"id"`
    Sources   []string `json:"sources"`
    Processed bool     `json:"processed"`
}

// CreateTransaction creates a new transaction
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, id string, sources []string) error {
    transaction := Transaction{
        ID:        id,
        Sources:   sources,
        Processed: false,
    }

    transactionJSON, err := json.Marshal(transaction)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, transactionJSON)
}

// ValidateTransactions validates and processes transactions based on similarity
func (s *SmartContract) ValidateTransactions(ctx contractapi.TransactionContextInterface) error {
    // Query all transactions
    queryString := `{"selector":{}}`
    resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
    if err != nil {
        return err
    }
    defer resultsIterator.Close()

    var transactions []Transaction
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return err
        }

        var transaction Transaction
        err = json.Unmarshal(queryResponse.Value, &transaction)
        if err != nil {
            return err
        }
        transactions = append(transactions, transaction)
    }

    // Validate and process transactions
    for i, txn := range transactions {
        similarCount := 0
        for j, otherTxn := range transactions {
            if i != j && compareTransactions(txn.Sources, otherTxn.Sources) {
                similarCount++
            }
        }

        if similarCount > len(transactions)/2 {
            txn.Processed = true
            transactionJSON, err := json.Marshal(txn)
            if err != nil {
                return err
            }
            err = ctx.GetStub().PutState(txn.ID, transactionJSON)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

// compareTransactions compares two transactions to see if they are similar
func compareTransactions(sources1 []string, sources2 []string) bool {
    similarityCount := 0
    for _, src1 := range sources1 {
        for _, src2 := range sources2 {
            if src1 == src2 {
                similarityCount++
            }
        }
    }
    return similarityCount > len(sources1)/2
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
