package vm

import (
	"math/big"
	"strings"

	"github.com/PlatONEnetwork/PlatONE-Go/common"
	"github.com/PlatONEnetwork/PlatONE-Go/common/syscontracts"
	"github.com/PlatONEnetwork/PlatONE-Go/core/state"
)

const (
	fwOpSuccess       CodeType = 0
	fwNoPermission    CodeType = 1
	fwInvalidArgument CodeType = 2
)

type FireWall struct {
	stateDB     StateDB
	caller      common.Address // msg.From()	contract.caller
	blockNumber *big.Int
}

func NewFireWall(evm *EVM, contract *Contract) *FireWall {
	return &FireWall{
		stateDB:     evm.StateDB,
		caller:      contract.CallerAddress,
		blockNumber: evm.BlockNumber,
	}
}

func (u *FireWall) isOwner(contractAddr common.Address) bool {
	contractOwnerAddr := u.stateDB.GetContractCreator(contractAddr)
	callerAddr := u.caller

	if callerAddr.Hex() == contractOwnerAddr.Hex() {
		return true
	}

	return false
}

func (u *FireWall) openFirewall(contractAddr common.Address) error {
	if !u.isOwner(contractAddr) {
		u.emitNotifyEvent(fwNoPermission, fwErrNotOwner.Error())
		return fwErrNotOwner
	}

	u.stateDB.OpenFirewall(contractAddr)

	u.emitNotifyEvent(fwOpSuccess, "fw start success")
	return nil
}

func (u *FireWall) closeFirewall(contractAddr common.Address) error {
	if !u.isOwner(contractAddr) {
		u.emitNotifyEvent(fwNoPermission, fwErrNotOwner.Error())
		return fwErrNotOwner
	}

	u.stateDB.CloseFirewall(contractAddr)

	u.emitNotifyEvent(fwOpSuccess, "fw close success")
	return nil
}

func (u *FireWall) fwClear(contractAddr common.Address, action string) error {
	if !u.isOwner(contractAddr) {
		u.emitNotifyEvent(fwNoPermission, fwErrNotOwner.Error())
		return fwErrNotOwner
	}

	act, err := state.NewAction(action)
	if err != nil {
		u.emitNotifyEvent(fwInvalidArgument, err.Error())
		return err
	}

	u.stateDB.FwClear(contractAddr, act)

	u.emitNotifyEvent(fwOpSuccess, "fw clear success")
	return nil
}

func (u *FireWall) fwAdd(contractAddr common.Address, action, lst string) error {
	if !u.isOwner(contractAddr) {
		u.emitNotifyEvent(fwNoPermission, fwErrNotOwner.Error())
		return fwErrNotOwner
	}

	act, err := state.NewAction(action)
	if err != nil {
		u.emitNotifyEvent(fwInvalidArgument, err.Error())
		return err
	}

	list, err := convertToFwElem(lst)
	if err != nil {
		u.emitNotifyEvent(fwInvalidArgument, err.Error())
		return err
	}

	u.stateDB.FwAdd(contractAddr, act, list)

	u.emitNotifyEvent(fwOpSuccess, "fw add success")
	return nil
}

func (u *FireWall) fwDel(contractAddr common.Address, action, lst string) error {
	if !u.isOwner(contractAddr) {
		u.emitNotifyEvent(fwNoPermission, fwErrNotOwner.Error())
		return fwErrNotOwner
	}

	act, err := state.NewAction(action)
	if err != nil {
		u.emitNotifyEvent(fwInvalidArgument, err.Error())
		return err
	}

	list, err := convertToFwElem(lst)
	if err != nil {
		u.emitNotifyEvent(fwInvalidArgument, err.Error())
		return err
	}

	u.stateDB.FwDel(contractAddr, act, list)

	u.emitNotifyEvent(fwOpSuccess, "fw delete success")
	return nil
}

// todo: input arguments type
func (u *FireWall) fwSet(contractAddr common.Address, act, lst string) error {
	if !u.isOwner(contractAddr) {
		u.emitNotifyEvent(fwNoPermission, fwErrNotOwner.Error())
		return fwErrNotOwner
	}

	action, err := state.NewAction(act)
	if err != nil {
		u.emitNotifyEvent(fwInvalidArgument, err.Error())
		return err
	}

	list, err := convertToFwElem(lst)
	if err != nil {
		u.emitNotifyEvent(fwInvalidArgument, err.Error())
		return err
	}

	u.stateDB.FwSet(contractAddr, action, list)

	u.emitNotifyEvent(fwOpSuccess, "fw reset success")
	return nil
}

func (u *FireWall) fwImport(contractAddr common.Address, data []byte) error {
	if !u.isOwner(contractAddr) {
		u.emitNotifyEvent(fwNoPermission, fwErrNotOwner.Error())
		return fwErrNotOwner
	}

	err := u.stateDB.FwImport(contractAddr, data)

	u.emitNotifyEvent(fwOpSuccess, "fw import success")
	return err
}

func (u *FireWall) getFwStatus(contractAddr common.Address) (*state.FwStatus, error) {
	if !u.isOwner(contractAddr) {
		return nil, fwErrNotOwner
	}

	fwStatus := u.stateDB.GetFwStatus(contractAddr)
	return &fwStatus, nil
}

func (u *FireWall) emitNotifyEvent(code CodeType, msg string) {
	topic := "Notify"
	emitEvent(syscontracts.FirewallManagementAddress, u.stateDB, u.blockNumber.Uint64(), topic, code, msg)
}

func convertToFwElem(l string) ([]state.FwElem, error) {
	var list = make([]state.FwElem, 0)

	elements := strings.Split(l, "|")
	for _, e := range elements {
		tmp := strings.Split(e, ":")
		if len(tmp) != 2 {
			/// log.Warn("FW : error, wrong function parameters!")
			return nil, ErrFwRule
		}

		addr := ZeroAddress
		addrStr := tmp[0]
		api := tmp[1]
		if addrStr == "*" {
			addr = state.FwWildchardAddr
		}
		fwElem := state.FwElem{Addr: addr, FuncName: api}
		list = append(list, fwElem)
	}

	return list, nil
}
