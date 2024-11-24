#!/bin/bash

# Hyperledger Fabric configurations
CHANNEL_NAME="mychannel"
CC_NAME="mychaincode"
ORDERER_NAME="myorderer"

# Function to create a transaction
CreateTransaction() {
    TXN_ID=$1
    SOURCE1=$2
    SOURCE2=$3
    SOURCE3=$4
    echo "Creating transaction $TXN_ID..."
    peer chaincode invoke -o $ORDERER_NAME -C $CHANNEL_NAME -n $CC_NAME -c "{\"function\":\"CreateTransaction\",\"Args\":[\"$TXN_ID\",\"$SOURCE1\",\"$SOURCE2\",\"$SOURCE3\"]}"
}

# Function to validate a transaction
ValidateTransaction() {
    TXN_ID=$1
    echo "Validating transaction $TXN_ID..."
    peer chaincode invoke -o $ORDERER_NAME -C $CHANNEL_NAME -n $CC_NAME -c "{\"function\":\"ValidateTransaction\",\"Args\":[\"$TXN_ID\"]}"
}

# Main script
echo "Starting bash script..."

# Create transactions
CreateTransaction txn1 commandA commandA commandB
CreateTransaction txn2 commandC commandD commandC
CreateTransaction txn3 commandE commandF commandG

# Wait for transactions to be committed (you may need to adjust the wait time)
sleep 5

# Validate transactions
ValidateTransaction txn1
ValidateTransaction txn2
ValidateTransaction txn3

echo "Bash script completed."
