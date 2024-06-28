package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/planetarium/bencodex-go"
	"github.com/planetarium/bencodex-go/util"
)

// Parse the serialized payload of a block transaction from the GraphQL query response
// Get response from https://9c-main-rpc-1.nine-chronicles.com/graphql/explorer
func blockTransactionsParseExample() {
	const path9c = "https://9c-main-rpc-1.nine-chronicles.com/graphql/explorer"

	// Make GraphQL query request
	query := `{
		blockQuery{
			blocks(desc: true, limit: 1) {
				transactions {
					serializedPayload
				}
			}
		}
	}`

	// Create the request body
	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}

	// Send the request
	resp, err := http.Post(path9c, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the response body
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return
	}

	// Print the response body
	fmt.Printf("%#v\n", response)
	fmt.Println()

	var serializedPayloadEncodedList [][]byte
	for _, transaction := range response.Data.BlockQuery.Blocks[0].Transactions {
		serializedPayloadEncoded, err := base64.StdEncoding.DecodeString(transaction.SerializedPayload)
		if err != nil {
			fmt.Println("Error decoding serialized payload:", err)
			return
		}
		serializedPayloadEncodedList = append(serializedPayloadEncodedList, serializedPayloadEncoded)
	}

	var serializedPayloadList []any
	for _, serializedPayloadEncoded := range serializedPayloadEncodedList {
		value, err := bencodex.Decode(serializedPayloadEncoded)
		if err != nil {
			fmt.Println("Error decoding bencodex value:", err)
			return
		}
		serializedPayloadList = append(serializedPayloadList, value)
	}

	files, err := filepath.Glob("bencodex_map_data_*.json")
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}
	files, err = filepath.Glob("bencodex_encoded_*.dat")
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}
	files, err = filepath.Glob("bencodex_serializedPayloadEncoded_*.dat")
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return
		}
	}

	for i, serializedPayload := range serializedPayloadList {
		fmt.Printf("Serialized Payload %d\n:%v\n\n", i, serializedPayload)

		mapData, err := util.ConvertToBencodexMapData(serializedPayload)
		if err != nil {
			fmt.Println("Error converting to Bencodex map data:", err)
			return
		}

		fmt.Printf("Bencodex Map Data %d:\n%v\n\n", i, mapData)

		jsonData, err := json.MarshalIndent(mapData, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling Bencodex map data:", err)
			return
		}

		err = os.WriteFile(fmt.Sprintf("bencodex_map_data_%d.json", i), jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing Bencodex map data:", err)
			return
		}
	}

	files, err = filepath.Glob("bencodex_map_data_*.json")
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}
	for _, file := range files {
		i := 0
		_, err := fmt.Sscanf(file, "bencodex_map_data_%d.json", &i)
		if err != nil {
			fmt.Println("Error extracting number from file name:", err)
			return
		}

		jsonData, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		var preData map[string]any
		err = json.Unmarshal(jsonData, &preData)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON data: %v", err)
		}
		data, err := util.ParseBencodexMapData(preData)
		if err != nil {
			fmt.Printf("Error parsing Bencodex map data: %v", err)
		}

		// Encode the data
		encoded, err := bencodex.Encode(data)
		if err != nil {
			fmt.Println("Error encoding data:", err)
		}

		serializedPayloadEncoded := serializedPayloadEncodedList[i]
		if !bytes.Equal(encoded, serializedPayloadEncoded) {
			fmt.Println("Encoded data does not match serialized payload")
		}
	}
}