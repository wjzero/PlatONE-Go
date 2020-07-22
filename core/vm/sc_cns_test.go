package vm

import (
	"testing"

	"github.com/PlatONEnetwork/PlatONE-Go/common"

	"github.com/stretchr/testify/assert"
)

const (
	testName = "tofu"
)

var (
	testAddr1 = common.HexToAddress("0x0000000000000000000000000000000000000123")
	testAddr2 = common.HexToAddress("0x0000000000000000000000000000000000000456")
	testAddr3 = common.HexToAddress("0x0000000000000000000000000000000000000789")
	testAddr4 = common.HexToAddress("0x0000000000000000000000000000000000000101")

	testOrigin = common.HexToAddress("0x0000000000000000000000000000000000000afb")
	testCaller = common.HexToAddress("0x0000000000000000000000000000000000000afc")
)

var (
	cns = new(CnsWrapper)
	key = make([]string, 0)
)

var testCases = []*ContractInfo{
	{Name: testName, Version: "0.0.0.1", Address: testAddr1, Origin: testOrigin},
	{Name: testName, Version: "0.0.0.2", Address: testAddr2, Origin: testOrigin},
	{Name: testName, Version: "0.0.0.3", Address: testAddr3, Origin: testOrigin},
	{Name: testName, Version: "0.0.0.4", Address: testAddr4, Origin: testOrigin},
	{Name: "bob", Version: "0.0.0.1", Address: testAddr1, Origin: testOrigin},
}

func TestLatestVersion(t *testing.T) {
	testCase := []struct {
		ver1     string
		ver2     string
		expected int
	}{
		{"0.0.0.0", "0.0.2.0", -1},
		{"1.0.0.0", "0.0.0.1", 1},
		{"1.0.0.0", "1.0.0.0", 0},
		{"1.0.0.0", "0.0.0.1", 1},
		{"1.0.0.03011", "1.0.0.0301", 1},
		{"1.0.0.07141108", "1.0.0.07141130", -1},
	}

	for _, data := range testCase {
		result := verCompare(data.ver1, data.ver2)

		assert.Equal(t, data.expected, result, "cnsRecall FAILED")

		if result == 1 {
			t.Logf("ver1 %s is larger than ver2 %s\n", data.ver1, data.ver2)
		} else {
			t.Logf("ver1 %s is smaller than ver2 %s\n", data.ver1, data.ver2)
		}
	}
}

/*
func TestSerializeCnsInfo(t *testing.T) {
	cnsInfoArray := make([]*ContractInfo, 0)
	cnsInfo := newContractInfo("tofu", "0.0.0.1", common.HexToAddress("0x123"), common.HexToAddress("0x123"))
	cnsInfoArray = append(cnsInfoArray, cnsInfo, cnsInfo, cnsInfo)

	sBytes, err := serializeCnsInfo(codeOk, msgOk, cnsInfoArray)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", sBytes)
}*/

func TestCnsManager_cnsRegister(t *testing.T) {
	result, err := cns.cnsRegister("alice", "0.0.0.1", testAddr1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, success, result, "cnsRegister FAILED")
}

func TestCnsManager_getContractAddress(t *testing.T) {

	testCasesSub := []struct {
		name     string
		version  string
		expected common.Address
	}{
		{testName, "0.0.0.2", testAddr2},
		{testName, "latest", testAddr4},
	}

	for _, data := range testCasesSub {
		result, err := cns.getContractAddress(data.name, data.version)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, data.expected, common.HexToAddress(result), "getContractAddress FAILED")
		t.Log(result)
	}
}

func TestCnsManager_cnsRecall(t *testing.T) {

	curVersion := cns.base.cMap.getCurrentVer(testName)

	result, err := cns.cnsRedirect(testName, testCases[2].Version)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, success, result, "cnsRecall FAILED")

	actVersion := cns.base.cMap.getCurrentVer(testName)
	expVersion := testCases[2].Version
	assert.Equal(t, expVersion, actVersion, "cnsRecall FAILED")

	t.Logf("before: %s, after cnsRecall: %s\n", curVersion, actVersion)
}

func TestCnsManager_ifRegisteredByName(t *testing.T) {
	testCasesSub := []struct {
		name     string
		expected int32
	}{
		{testName, registered},
		{"tom", unregistered},
	}

	for _, data := range testCasesSub {
		result, err := cns.ifRegisteredByName(data.name)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, data.expected, result, "ifRegisteredByName FAILED")
	}
}

func TestCnsManager_getRegisteredContractsByRange(t *testing.T) {
	result, err := cns.getRegisteredContractsByRange(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("the registered contracts result is %s\n", result)
}
