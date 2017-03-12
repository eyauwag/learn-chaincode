/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/binary"
	"strconv"

	//"container/list"
	"errors"
	"fmt"
	//"strconv"

	//"unsafe"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var chaincodeLogger = logging.MustGetLogger("main")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type DonatorRecord struct {
	donator string
	money   int
	index   int
	project string
	time    timestamp.Timestamp
}

type Donator struct {
	donator string
	time    timestamp.Timestamp
	money   int
	count   int
	record  []DonatorRecord
}

type Receiver struct {
	ProfiteReceiver string
	time            timestamp.Timestamp
	money           int
	project         string
}

type Project struct {
	Project string
	time    timestamp.Timestamp
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A string // Entity
	var Aval int // Asset holding
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running %s" + function)
	chaincodeLogger.Info("invoke is running %s" + function)
	// Handle different functions
	switch function {
	case "init":
		return t.Init(stub, "init", args)
	case "registerDonater":
		return t.registerDonater(stub, args)
	case "registerProfiteReceiver":
		return t.registerProfiteReceiver(stub, args)
	case "donate":
		return t.donate(stub, args)
	case "assign":
		return t.assign(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	switch function {
	case "trackproject":
		return t.trackproject(stub, args)
	case "trackreceiver":
		return t.trackreceiver(stub, args)
	case "trackdonator":
		return t.trackdonator(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) registerDonater(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var donatorname string
	var err error
	fmt.Println("running registerDonater()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	var newdonator Donator
	newdonator.donator = args[0]

	//timestamp.Timestamp* timeStamp;
	timeStamp, err := stub.GetTxTimestamp()
	newdonator.time = *timeStamp
	newdonator.money = 0
	newdonator.count = 0

	var bin_buf bytes.Buffer

	binary.Write(&bin_buf, binary.BigEndian, newdonator)
	err = stub.PutState(donatorname, bin_buf.Bytes())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) registerProfiteReceiver(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var receiverName string
	var err error
	fmt.Println("running registerProfiteReceiver()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	var newreceiver Receiver
	newreceiver.ProfiteReceiver = args[0]

	timeStamp, err := stub.GetTxTimestamp()
	newreceiver.time = *timeStamp
	newreceiver.money = 0

	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, newreceiver)
	err = stub.PutState(receiverName, bin_buf.Bytes())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) donate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var money int
	var err error
	fmt.Println("running registerProfiteReceiver()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name of the donator and money and project to set")
	}
	var record DonatorRecord
	record.donator = args[0]
	record.project = args[1]
	money, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}
	record.money = money
	timeStamp, err := stub.GetTxTimestamp()
	record.time = *timeStamp
	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, record)

	err = stub.PutState(record.donator, bin_buf.Bytes())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) assign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("running registerProfiteReceiver()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	//projectName := args[0]
	receiverName := args[1]
	/*money, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}*/
	timeStamp, err := stub.GetTxTimestamp()
	//bs := unsafe.StringBytes(timeStamp.String())
	err = stub.PutState(receiverName, []byte(timeStamp.String()))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) trackproject(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var project, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the project to query")
	}

	project = args[0]
	valAsbytes, err := stub.GetState(project)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + project + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) trackreceiver(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var receiver, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the receiver to query")
	}

	//key = args[0]
	valAsbytes, err := stub.GetState(receiver)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + receiver + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) trackdonator(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var donatorname, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting donator name to query")
	}

	donatorname = args[0]
	valAsbytes, err := stub.GetState(donatorname)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + donatorname + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
