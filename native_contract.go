package ontology_go_sdk

import (
	"fmt"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	cutils "github.com/ontio/ontology/core/utils"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
)

var (
	ONG_CONTRACT_ADDRESS, _ = utils.AddressFromHexString("0200000000000000000000000000000000000000")
)

var (
	ONG_CONTRACT_VERSION = byte(0)
)

type NativeContract struct {
	ontSdk *OntologySdk
	Ong    *Ong
}

func newNativeContract(ontSdk *OntologySdk) *NativeContract {
	native := &NativeContract{ontSdk: ontSdk}
	native.Ong = &Ong{native: native, ontSdk: ontSdk}
	return native
}

func (this *NativeContract) NewNativeInvokeTransaction(
	chainID uint64,
	gasPrice,
	gasLimit uint64,
	version byte,
	contractAddress common.Address,
	method string,
	params []interface{},
) (*types.MutableTransaction, error) {
	if params == nil {
		params = make([]interface{}, 0, 1)
	}
	//Params cannot empty, if params is empty, fulfil with empty string
	if len(params) == 0 {
		params = append(params, "")
	}
	invokeCode, err := cutils.BuildNativeInvokeCode(contractAddress, version, method, params)
	if err != nil {
		return nil, fmt.Errorf("BuildNativeInvokeCode error:%s", err)
	}
	return this.ontSdk.NewInvokeTransaction(chainID, gasPrice, gasLimit, invokeCode), nil
}

func (this *NativeContract) InvokeNativeContract(
	chainID uint64,
	gasPrice,
	gasLimit uint64,
	singer *Account,
	version byte,
	contractAddress common.Address,
	method string,
	params []interface{},
) (common.Uint256, error) {
	tx, err := this.NewNativeInvokeTransaction(chainID, gasPrice, gasLimit, version, contractAddress, method, params)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.ontSdk.SignToTransaction(tx, singer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.ontSdk.SendTransaction(tx)
}

func (this *NativeContract) PreExecInvokeNativeContract(
	chainID uint64,
	contractAddress common.Address,
	version byte,
	method string,
	params []interface{},
) (*sdkcom.PreExecResult, error) {
	tx, err := this.NewNativeInvokeTransaction(chainID, 0, 0, version, contractAddress, method, params)
	if err != nil {
		return nil, err
	}
	return this.ontSdk.PreExecTransaction(tx)
}

type Ong struct {
	ontSdk *OntologySdk
	native *NativeContract
}

func (this *Ong) NewTransferTransaction(chainID uint64, gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &ont.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.NewMultiTransferTransaction(chainID, gasPrice, gasLimit, []*ont.State{state})
}

func (this *Ong) Transfer(chainID uint64, gasPrice, gasLimit uint64, from *Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferTransaction(chainID, gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.ontSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.ontSdk.SendTransaction(tx)
}

func (this *Ong) NewMultiTransferTransaction(chainID uint64, gasPrice, gasLimit uint64, states []*ont.State) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		chainID,
		gasPrice,
		gasLimit,
		ONG_CONTRACT_VERSION,
		ONG_CONTRACT_ADDRESS,
		ont.TRANSFER_NAME,
		[]interface{}{states})
}

func (this *Ong) MultiTransfer(chainID uint64, gasPrice, gasLimit uint64, states []*ont.State, signer *Account) (common.Uint256, error) {
	tx, err := this.NewMultiTransferTransaction(chainID, gasPrice, gasLimit, states)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.ontSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.ontSdk.SendTransaction(tx)
}

func (this *Ong) NewTransferFromTransaction(chainID uint64, gasPrice, gasLimit uint64, sender, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &ont.TransferFrom{
		Sender: sender,
		From:   from,
		To:     to,
		Value:  amount,
	}
	return this.native.NewNativeInvokeTransaction(
		chainID,
		gasPrice,
		gasLimit,
		ONG_CONTRACT_VERSION,
		ONG_CONTRACT_ADDRESS,
		ont.TRANSFERFROM_NAME,
		[]interface{}{state},
	)
}

func (this *Ong) TransferFrom(chainID uint64, gasPrice, gasLimit uint64, sender *Account, from, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferFromTransaction(chainID, gasPrice, gasLimit, sender.Address, from, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.ontSdk.SignToTransaction(tx, sender)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.ontSdk.SendTransaction(tx)
}

func (this *Ong) NewApproveTransaction(chainID uint64, gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &ont.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.native.NewNativeInvokeTransaction(
		chainID,
		gasPrice,
		gasLimit,
		ONG_CONTRACT_VERSION,
		ONG_CONTRACT_ADDRESS,
		ont.APPROVE_NAME,
		[]interface{}{state},
	)
}

func (this *Ong) Approve(chainID uint64, gasPrice, gasLimit uint64, from *Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewApproveTransaction(chainID, gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.ontSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.ontSdk.SendTransaction(tx)
}

func (this *Ong) Allowance(chainID uint64, from, to common.Address) (uint64, error) {
	type allowanceStruct struct {
		From common.Address
		To   common.Address
	}
	preResult, err := this.native.PreExecInvokeNativeContract(
		chainID,
		ONG_CONTRACT_ADDRESS,
		ONG_CONTRACT_VERSION,
		ont.ALLOWANCE_NAME,
		[]interface{}{&allowanceStruct{From: from, To: to}},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Ong) Symbol(chainID uint64) (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		chainID,
		ONG_CONTRACT_ADDRESS,
		ONG_CONTRACT_VERSION,
		ont.SYMBOL_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Ong) BalanceOf(chainID uint64, address common.Address) (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		chainID,
		ONG_CONTRACT_ADDRESS,
		ONG_CONTRACT_VERSION,
		ont.BALANCEOF_NAME,
		[]interface{}{address[:]},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Ong) Name(chainID uint64) (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		chainID,
		ONG_CONTRACT_ADDRESS,
		ONG_CONTRACT_VERSION,
		ont.NAME_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Ong) Decimals(chainID uint64) (byte, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		chainID,
		ONG_CONTRACT_ADDRESS,
		ONG_CONTRACT_VERSION,
		ont.DECIMALS_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	decimals, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return byte(decimals.Uint64()), nil
}

func (this *Ong) TotalSupply(chainID uint64) (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		chainID,
		ONG_CONTRACT_ADDRESS,
		ONG_CONTRACT_VERSION,
		ont.TOTAL_SUPPLY_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}
