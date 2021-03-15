package main

import (
	"encoding/base64"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const ADMIN = "admin"
const CLIENT = "client"
const OWNER = "owner"

// Framework example simple Chaincode implementation
type Framework struct {
}

func (t *Framework) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Framework Init")
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	creator := base64.StdEncoding.EncodeToString(creatorByte)

	err = stub.PutState(creator, []byte(ADMIN))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(OWNER, []byte(creator))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(creatorByte)
}

func (t *Framework) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Framework Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "chmod" {
		return t.chmod(stub, args)
	} else if function == "save" {
		return t.save(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}
	return shim.Success(nil)
}

func (t *Framework) chmod(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	creator := base64.StdEncoding.EncodeToString(creatorByte)
	roleBytes, err := stub.GetState(creator)
	if err != nil {
		return shim.Error(fmt.Sprintf("get state of %s err", creator))
	}

	if roleBytes == nil || string(roleBytes) != ADMIN {
		return shim.Error("permission deny")
	}

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	userMsp := args[0]
	userRole := args[1]

	ownerBytes, err := stub.GetState(OWNER)
	if err != nil {
		return shim.Error(fmt.Sprintf("get state of %s err", OWNER))
	}

	if userMsp == string(ownerBytes) {
		return shim.Error("cannot chmod owner")
	}

	err = stub.PutState(userMsp, []byte(userRole))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *Framework) save(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	creator := base64.StdEncoding.EncodeToString(creatorByte)
	roleBytes, err := stub.GetState(creator)
	if err != nil {
		return shim.Error(fmt.Sprintf("get state of %s err", creator))
	}

	if roleBytes == nil || string(roleBytes) != ADMIN {
		return shim.Error("permission deny")
	}

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	id := args[0]
	myData := args[1]

	err = stub.PutState(id, []byte(myData))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *Framework) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	creatorByte, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	creator := base64.StdEncoding.EncodeToString(creatorByte)
	roleBytes, err := stub.GetState(creator)
	if err != nil {
		return shim.Error(fmt.Sprintf("get state of %s err", creator))
	}

	if roleBytes == nil {
		return shim.Error("permission deny")
	}
	if string(roleBytes) != ADMIN && string(roleBytes) != CLIENT {
		return shim.Error("permission deny")
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	id := args[0]

	myData, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprintf("get state of %s err", id))
	}

	return shim.Success(myData)
}

func main() {
	err := shim.Start(new(Framework))
	if err != nil {
		fmt.Printf("Error starting Framework chaincode: %s", err)
	}
}
