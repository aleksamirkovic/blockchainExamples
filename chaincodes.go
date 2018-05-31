// +build main1

package main

import "fmt"


func main() {
	fmt.Println("Welcome to the chaincode")
}

type Car struct {
	price float64
	color string
	name string
	ownerID string
}
type Person struct {
	firstName string
	lastName string
	balance float64
	personID string
}

//change color function
{
	
}

func (s *CarExample) transferCarAsset(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Info("saveContract")

	if len(args) != 2 {
		logger.Error("Incorrect number of arguments. Expecting 2, but was " + strconv.Itoa(len(args)))
		return Error(http.StatusBadRequest, "Incorrect number of arguments. Expecting 2, but was "+strconv.Itoa(len(args)))
	}

	// Check arguments
	if len(args[0]) <= 0 {
		logger.Error("[1] Car ID is required")
		return Error(http.StatusBadRequest, "[1] Car ID is required")
	}
	if len(args[1]) <= 0 {
		logger.Error("[2] Contract string(JSON) is required")
		return Error(http.StatusBadRequest, "[2] Contract string(JSON) is required")
	}

	var carID = args[0]
	var buyerID = args[1]
	var car *Car

	buyer, err := APIstub.GetState(buyerID)
	if err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	// Parse contract json string
	car, err = APIstub.GetState(carID)
	if err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	ownerID := car.ownerID
	carPrice := car.price

	owner, err := APIstub.GetState(ownerID)
	if err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	if buyer.balance < carPrice {
		logger.Error("Not enough funds")
		return Error(http.StatusBadRequest, "ContractID parameter differs from Contract.ContractID")
	}

	buyer.balance = buyer.balance - carPrice
	owner.balance = owner.balance + carPrice
	car.ownerID = buyerID

	err = APIstub.PutState(ownerID, owner)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return Error(http.StatusBadRequest, err.Error())
	}
	err = APIstub.PutState(carID, car)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return Error(http.StatusBadRequest, err.Error())
	}
	err = APIstub.PutState(buyerID, buyer)

	if err != nil {
		logger.Info("PutState error: " + err.Error())
		return Error(http.StatusBadRequest, err.Error())
	}

	return Success(http.StatusOK, "Ownership changed.", car)
}

