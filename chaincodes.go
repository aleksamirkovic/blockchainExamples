// +build main1

package main

import (
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"net/http"
	"strings"
	"encoding/json"
	"time"
)

var logger = shim.NewLogger("RealMarket")

type CarExample struct{}

func main() {
	fmt.Println("Welcome to the chaincode")
}

type Car struct {
	carID string
	name string
	ownerID string
	price float64
	color string
}
type Person struct {
	personID string
	firstName string
	lastName string
	balance float64
}

//change color function
{
	
}

func main() {
	// Create a new Smart Contract
	err := shim.Start(new(CarExample))
	if err != nil {
		logger.Error("Error creating new Smart Contract")
	}
}

func (c *CarExample) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Invoke")

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	logger.Info("Invoked function: " + function)

	// Route to the appropriate handler function to interact with the ledger appropriately
	switch function {
	case "transferCarAsset":
		return c.transferCarAsset(APIstub, args)
	case "getAssetById":
		return c.getAssetById(APIstub, args)
	default:
		logger.Error("Invalid Smart Contract function name.")
		return Error(http.StatusNotImplemented, "Invalid Smart Contract function name.")
	}

}

func (s *RealMarket) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Init")
	car := Car{"car1","Audi","ownerId",344.4,"blue"}
	buyerOk := Person{"buyerOkId","Mark","Rush",350}
	buyerNotOk := Person{"buyerNotOkId""John","Green",250}
	owner := Person{"ownerId","Steven","Tyler",600}
	APIstub.PutState(car.carID,car)
	APIstub.PutState(buyerOk.personID,buyerOk)
	APIstub.PutState(buyerNotOk.personID,buyerNotOk)
	APIstub.PutState(owner.personID,owner)
	return Success(http.StatusOK, "initialized", nil)
}

func (s *CarExample) transferCarAsset(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Info("saveContract")

	if len(args) != 2 {
		logger.Error("Incorrect number of arguments. Expecting 2, but was " + strconv.Itoa(len(args)))
		return Error(http.StatusBadRequest, "Incorrect number of arguments. Expecting 2, but was "+strconv.Itoa(len(args)))
	}

	// Check arguments
	if len(args[0]) <= 0 {
		logger.Error("Car ID is required")
		return Error(http.StatusBadRequest, "Car ID is required")
	}
	if len(args[1]) <= 0 {
		logger.Error("Buyer ID is required")
		return Error(http.StatusBadRequest, "Buyer  is required")
	}

	var carID = args[0]
	var buyerID = args[1]
	var car *Car

	//get buyer person from world state
	buyer, err := APIstub.GetState(buyerID)
	if err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	//get car from world state
	car, err = APIstub.GetState(carID)
	if err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	ownerID := car.ownerID
	carPrice := car.price

	//get owner person from world state
	owner, err := APIstub.GetState(ownerID)
	if err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	//check if buyer has enough funds
	if buyer.balance < carPrice {
		logger.Error("Not enough funds")
		return Error(http.StatusBadRequest, "Not enough funds")
	}

	//if everything is ok then change balances
	//and change ownership over the car
	buyer.balance = buyer.balance - carPrice
	owner.balance = owner.balance + carPrice
	car.ownerID = buyerID

	//persist owner with updated balance
	err = APIstub.PutState(ownerID, owner)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return Error(http.StatusBadRequest, err.Error())
	}
	//persist car with updated ownerID
	err = APIstub.PutState(carID, car)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return Error(http.StatusBadRequest, err.Error())
	}
	//persist owner with updated balance
	err = APIstub.PutState(buyerID, buyer)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return Error(http.StatusBadRequest, err.Error())
	}

	return Success(http.StatusOK, "Ownership changed.", car)
}

