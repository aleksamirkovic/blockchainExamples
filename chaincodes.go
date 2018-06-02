/*
   Title: Smart Contract handling transaction CRUD
   Author: Branko Terzic
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("SmartContract")

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the transaction model
type Car struct {
	carID   string
	name    string
	ownerID string
	price   float64
	color   string
}
type Person struct {
	personID  string
	firstName string
	lastName  string
	balance   float64
}

/*
 * The Init method is called when the Smart Contract "txMan" is instantiated by the blockchain network
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {

	car := Car{"car1", "Audi", "ownerId", 344.4, "blue"}
	buyerOk := Person{"buyerOkId", "Mark", "Rush", 350}
	owner := Person{"ownerId", "Steven", "Tyler", 600}

	carBytes, err := json.Marshal(car)
	buyerOkBytes, err := json.Marshal(buyerOk)
	ownerBytes, err := json.Marshal(owner)

	if err != nil {
		fmt.Println("Error")
	}

	APIstub.PutState(car.carID, carBytes)
	APIstub.PutState(buyerOk.personID, buyerOkBytes)
	APIstub.PutState(owner.personID, ownerBytes)

	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "transferCarAsset" {
		return s.transferCarAsset(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) transferCarAsset(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		logger.Error("Incorrect number of arguments. Expecting 2, but was " + strconv.Itoa(len(args)))
	}

	// Check arguments

	if len(args[0]) <= 0 {
		logger.Error("Car ID is required")
		return shim.Error("Car ID is required")

	}
	if len(args[1]) <= 0 {
		return shim.Error("buyer ID is required")
	}

	var carID = args[0]
	var buyerID = args[1]
	var car *Car
	var buyer *Person
	var owner *Person

	//get buyer person from world state
	buyerBytes, err := APIstub.GetState(buyerID)
	if err != nil {
		return shim.Error("No such person in world state")
	}

	//get car from world state
	carBytes, err := APIstub.GetState(carID)
	if err != nil {
		return shim.Error("No such car in world state")
	}
	json.Unmarshal(carBytes, &car)
	json.Unmarshal(buyerBytes, &buyer)

	ownerID := car.ownerID
	carPrice := car.price

	//get owner person from world state
	ownerBytes, err := APIstub.GetState(ownerID)

	if err != nil {
		return shim.Error("No such person in world state")
	}
	json.Unmarshal(ownerBytes, &owner)

	//check if buyer has enough funds
	if buyer.balance < carPrice {
		logger.Error("Not enough funds")
		return shim.Error("Error")
	}

	//if everything is ok then change balances
	//and change ownership over the car
	buyer.balance = buyer.balance - carPrice
	owner.balance = owner.balance + carPrice
	car.ownerID = buyerID

	//persist owner with updated balance
	ownerBytes, err = json.Marshal(owner)
	err = APIstub.PutState(ownerID, ownerBytes)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return shim.Error("PutState error: " + err.Error())
	}
	//persist car with updated ownerID
	carBytes, err = json.Marshal(car)
	err = APIstub.PutState(carID, carBytes)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return shim.Error("PutState error: " + err.Error())
	}
	//persist owner with updated balance
	buyerBytes, err = json.Marshal(buyer)
	err = APIstub.PutState(buyerID, buyerBytes)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return shim.Error("PutState error: " + err.Error())
	}

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
