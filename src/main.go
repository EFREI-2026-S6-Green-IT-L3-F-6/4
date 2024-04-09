package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SupplyChainContract structure définit le contrat intelligent
type SupplyChainContract struct {
    contractapi.Contract
}

// Product représente un bien dans la chaîne d'approvisionnement
type Product struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Owner       string `json:"owner"`
}

// Transaction enregistre un changement de propriété
type Transaction struct {
    ProductID   string    `json:"productId"`
    FromOwner   string    `json:"fromOwner"`
    ToOwner     string    `json:"toOwner"`
    Timestamp   time.Time `json:"timestamp"`
}

// InitLedger ajoute un ensemble de base de produits dans le registre
func (s *SupplyChainContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    products := []Product{
        {ID: "prod1", Name: "Laptop", Description: "Un ordinateur portable haute performance", Owner: "Fabricant"},
        {ID: "prod2", Name: "Smartphone", Description: "Un smartphone innovant", Owner: "Fabricant"},
    }

    for _, product := range products {
        productJSON, err := json.Marshal(product)
        if err != nil {
            return err
        }
        err = ctx.GetStub().PutState(product.ID, productJSON)
        if err != nil {
            return fmt.Errorf("échec de l'enregistrement dans l'état global. %v", err)
        }
    }
    return nil
}

// CreateProduct ajoute un nouveau produit dans le registre
func (s *SupplyChainContract) CreateProduct(ctx contractapi.TransactionContextInterface, id string, name string, description string, owner string) error {
    product := Product{
        ID:          id,
        Name:        name,
        Description: description,
        Owner:       owner,
    }

    productJSON, err := json.Marshal(product)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, productJSON)
}

// QueryProduct retourne le produit stocké dans le registre avec l'identifiant donné
func (s *SupplyChainContract) QueryProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
    productJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("échec de la lecture de l'état global: %v", err)
    }
    if productJSON == nil {
        return nil, fmt.Errorf("le produit %s n'existe pas", id)
    }

    var product Product
    err = json.Unmarshal(productJSON, &product)
    if err != nil {
        return nil, err
    }

    return &product, nil
}

// TransferProduct enregistre le changement de propriété d'un produit
func (s *SupplyChainContract) TransferProduct(ctx contractapi.TransactionContextInterface, productId string, newOwner string) error {
    productJSON, err := ctx.GetStub().GetState(productId)
    if err != nil {
        return err
    }
    if productJSON == nil {
        return fmt.Errorf("produit non trouvé")
    }

    var product Product
    err = json.Unmarshal(productJSON, &product)
    if err != nil {
        return err
    }

    // Créer un enregistrement de transaction
    transaction := Transaction{
        ProductID:   product.ID,
        FromOwner:   product.Owner,
        ToOwner:     newOwner,
        Timestamp:   time.Now(),
    }
    transactionJSON, err := json.Marshal(transaction)
    if err != nil {
        return err
    }

    // Mettre à jour le propriétaire du produit
    product.Owner = newOwner
    productJSON, err = json.Marshal(product)
    if err != nil {
        return err
    }
    err = ctx.GetStub().PutState(product.ID, productJSON)
    if err != nil {
        return err
    }

    // Utiliser une clé composite pour les transactions pour permettre des requêtes par plage
    txnKey, err := ctx.GetStub().CreateCompositeKey("txn", []string{product.ID, transaction.Timestamp.Format(time.RFC3339)})
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(txnKey, transactionJSON)
}

// La fonction main est seulement pertinente en mode test unitaire. Incluse ici pour complétude.
func main() {
    chaincode, err := contractapi.NewChaincode(new(SupplyChainContract))

    if err != nil {
        fmt.Printf("Erreur lors de la création du contrat intelligent pour la chaîne d'approvisionnement: %s", err)
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Erreur lors du démarrage du contrat intelligent pour la chaîne d'approvisionnement: %s", err.Error())
    }
}
