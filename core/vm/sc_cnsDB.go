package vm

import (
	"fmt"
	"github.com/PlatONEnetwork/PlatONE-Go/common"
	"github.com/PlatONEnetwork/PlatONE-Go/log"
	"strconv"
	"strings"
)

const (
	NAME = "cnsManager"
	CNS_TOTAL = "total"
	CNS_LATEST = "latest"
)

type cnsMap struct{
	/// StateDB
	stateDB 	StateDB
	CodeAddr	*common.Address
}

func NewCnsMap(db StateDB, addr *common.Address) *cnsMap {
	return &cnsMap{db, addr}
}

func (c *cnsMap) setState(key, value []byte) {
	str := fmt.Sprintf("[CNS] setState %v, key is %v, value is %v", *c.CodeAddr, key, value)
	log.Debug(str)
	///c.stateDB.SetState(*c.CodeAddr, key, value)

	c.stateDB.SetState(*c.CodeAddr, key, value)
}

func (c *cnsMap) getState(key []byte) []byte {
	///return c.GetState(*c.CodeAddr, key)

	///value := c.stateDB.GetState(*c.CodeAddr, key)

	value := c.stateDB.GetState(*c.CodeAddr, key)

	str := fmt.Sprintf("[CNS] %s getState, key is %v, value is %v", c.CodeAddr.Hex(), key, value)
	log.Debug(str)
	return value
}

func (c *cnsMap) getKey(index int) []byte {
	indexStr := strconv.Itoa(index)
	value := c.getState(wrapper(indexStr))

	return value
}

func (c *cnsMap) find(key []byte) *ContractInfo {
	value := c.getState(key)
	if value == nil {
		return nil
	}

	cnsInfo, _ := decodeCnsInfo(value)
	return cnsInfo
}

func (c *cnsMap) get(index int) *ContractInfo {
	value := c.getKey(index)
	if value == nil {
		return nil
	}

	return c.find(value)
}

func (c *cnsMap) total() int {
	value := c.getState(totalWrapper())

	if value == nil || len(value) == 0 {
		return 0
	}

	totalStr := string(value)
	total, _ := strconv.Atoi(totalStr)
	return total
}

func (c *cnsMap) insert(key, value []byte) {
	total := c.total()
	index := strconv.Itoa(total)

	c.setState(key, value)
	c.setState(wrapper(index), key)

	update := strconv.Itoa(total + 1)
	c.setState(totalWrapper(), []byte(update))
}

func (c *cnsMap) update(key, value []byte) {		// todo overwrite?
	c.setState(key, value)
}

func wrapper(str string) []byte {
	return []byte(NAME + str)
}

func totalWrapper() []byte {
	return []byte(NAME + CNS_TOTAL)
}

// utils

func (cMap *cnsMap) isNameRegByOthers(name, origin string) bool {
	for index := 0; index < cMap.total(); index++{
		cnsInfo := cMap.get(index)
		if cnsInfo.Name == name && cnsInfo.Origin != origin {
			return true
		}
	}

	return false
}

func isNameRegByOthers_Method2(c *cnsMap, name, origin string) bool {
	for index := 0; index < c.total(); index++{
		key := c.getKey(index)
		existedName := strings.Split(string(key), ":")[0]
		if existedName == name {
			cnsInfo := c.find(key)
			if cnsInfo.Origin != origin {
				return true
			} else {
				return false
			}
		}
	}

	return false
}

func (c *cnsMap) getLargestVersion(name string) string {
	tempVersion := "0.0.0.0"

	for index := 0; index < c.total(); index++{
		cnsInfo := c.get(index)
		if cnsInfo.Name == name{
			if verCompare(cnsInfo.Version, tempVersion) == 1 {
				tempVersion = cnsInfo.Version
			}
		}
	}

	return tempVersion
}

func (c *cnsMap) getLargestVersion_Method2(name string) string {
	tempVersion := "0.0.0.0"

	for index := 0; index < c.total(); index++{
		key := c.getKey(index)
		existedName := strings.Split(string(key), ":")[0]
		existedVersion := strings.Split(string(key), ":")[1]
		if existedName == name{
			if verCompare(existedVersion, tempVersion) == 1 {
				tempVersion = existedVersion
			}
		}
	}

	return tempVersion
}

//-------------------------------------------------------
// Method 2
func latestWrapper(name string) []byte {
	return []byte(NAME + CNS_LATEST + name)
}

func (c *cnsMap) getLatestVer(name string) string {
	ver := c.getState(latestWrapper(name))
	return string(ver)
}

func (c *cnsMap) updateLatestVer(name, ver string) {
	c.setState(latestWrapper(name), []byte(ver))
}