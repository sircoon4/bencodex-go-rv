package bencodex

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sircoon4/bencodex-go/util"
	"github.com/stretchr/testify/assert"
)

const encodedDataFilesPath = "spec/testsuite/*.dat"
const decodedDataFilesPath = "spec/testsuite/*.repr.json"

func TestBencodexEncode(t *testing.T) {
	testFiles, err := filepath.Glob(decodedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}

	testResults, err := filepath.Glob(encodedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}

	for i, file := range testFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			fmt.Println()
			fmt.Println("Test File:", file)

			// Read the test file
			jsonData, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}

			data, err := util.UnmarshalJsonRepr(jsonData)
			if err != nil {
				t.Fatal(err)
			}

			// Encode the data
			encoded, err := Encode(data)
			if err != nil {
				t.Fatal(err)
			}

			// Read the test result file
			result, err := os.ReadFile(testResults[i])
			if err != nil {
				t.Fatal(err)
			}

			// Compare the encoded data with the test result
			assert.Equal(t, result, encoded)
		})
	}
}

func TestBencodexDecode(t *testing.T) {
	testFiles, err := filepath.Glob(encodedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}

	testResultFiles, err := filepath.Glob(decodedDataFilesPath)
	if err != nil {
		t.Fatal(err)
	}

	for i, file := range testFiles {
		t.Run(filepath.Base(file), func(t *testing.T) {
			fmt.Println()
			fmt.Println("Test File:", file)

			// Read the test file
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}

			// Decode the encoded data
			decoded, err := Decode(data)
			if err != nil {
				t.Fatal(err)
			}

			// Read the test file
			jsonData, err := os.ReadFile(testResultFiles[i])
			if err != nil {
				t.Fatal(err)
			}

			result, err := util.UnmarshalJsonRepr(jsonData)
			if err != nil {
				t.Fatal(err)
			}

			// Compare the original data with the decoded data
			customizedAssertEqual(t, result, decoded)
		})
	}
}

func TestNilString(t *testing.T) {
	var encodedE, encodedA []byte
	var decoded any
	var err error

	encodedE = []byte("u0:")
	decoded, err = Decode(encodedE)
	if err != nil {
		t.Fatal(err)
	}
	encodedA, err = Encode(decoded)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, encodedE, encodedA)
}

func TestNilList(t *testing.T) {
	var encodedE, encodedA []byte
	var decoded any
	var err error

	encodedE = []byte("le")
	decoded, err = Decode(encodedE)
	if err != nil {
		t.Fatal(err)
	}
	encodedA, err = Encode(decoded)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, encodedE, encodedA)
}

func TestNilDictionary(t *testing.T) {
	var encodedE, encodedA []byte
	var decoded any
	var err error

	encodedE = []byte("de")
	decoded, err = Decode(encodedE)
	if err != nil {
		t.Fatal(err)
	}
	encodedA, err = Encode(decoded)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, encodedE, encodedA)
}
