package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Read the ABI JSON file
	abiFile, err := os.Open("erc20abi.json")
	if err != nil {
		log.Fatalf("Failed to open ABI file: %v", err)
	}
	defer abiFile.Close()

	byteValue, err := io.ReadAll(abiFile)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	// Unmarshal the ABI JSON
	var contractABI map[string]interface{}
	if err := json.Unmarshal(byteValue, &contractABI); err != nil {
		log.Fatalf("Failed to unmarshal ABI JSON: %v", err)
	}

	// Extract the "abi" field which contains the actual ABI
	abiJSON, err := json.Marshal(contractABI["abi"])
	if err != nil {
		log.Fatalf("Failed to marshal ABI: %v", err)
	}

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Function signature and parameters
	funcSignature := []byte("approve(address,uint256)")
	hash := crypto.Keccak256Hash(funcSignature).Hex()[:10] // First 4 bytes of keccak256 hash of the function signature

	// Define the spender address and the amount to approve
	spenderAddress := common.HexToAddress("0x6352a56caadC4F1E25CD6c75970Fa768A3304e64")
	amount := big.NewInt(1000) // 1000 tokens

	// Pack the parameters for the approve function
	data, err := parsedABI.Pack("approve", spenderAddress, amount)
	if err != nil {
		log.Fatalf("Failed to pack parameters: %v", err)
	}

	// Combine the function signature hash and packed parameters
	encodedData := fmt.Sprintf("%s%x", hash, data[4:])
	fmt.Printf("Encoded data: %s\n", encodedData)
}
