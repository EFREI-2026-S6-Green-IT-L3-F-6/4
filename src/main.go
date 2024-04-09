package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SupplyChainContract struct {
    contractapi.Contract
}

type Product struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

func (s *SupplyChainContract) CreateProduct(ctx contractapi.TransactionContextInterface, id string, name string, description string) error {
    product := Product{
        ID:          id,
        Name:        name,
        Description: description,
    }

    productAsBytes, _ := json.Marshal(product)

    return ctx.GetStub().PutState(id, productAsBytes)
}

func (s *SupplyChainContract) QueryProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
    productAsBytes, err := ctx.GetStub().GetState(id)

    if err != nil {
        return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if productAsBytes == nil {
        return nil, fmt.Errorf("product %s does not exist", id)
    }

    product := new(Product)
    _ = json.Unmarshal(productAsBytes, product)

    return product, nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SupplyChainContract))

    if err != nil {
        fmt.Printf("Error create supplychain chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting supplychain chaincode: %s", err.Error())
    }
}
