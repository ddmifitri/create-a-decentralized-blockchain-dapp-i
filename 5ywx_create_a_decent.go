package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Decent struct {
	client  *ethclient.Client
	account *accounts.Account
}

func NewDecent(endpoint, privateKey string) (*Decent, error) {
 client, err := ethclient.Dial(endpoint)
 if err != nil {
  return nil, err
 }

 account, err := accounts.HexToECDSA(privateKey)
 if err != nil {
  return nil, err
 }

 return &Decent{client: client, account: account}, nil
}

func (d *Decent) DeployContract(contractName string, bytecode []byte) (common.Address, error) {
 auth, err := d.client.PendingNonceAt(context.Background(), d.account.Address)
 if err != nil {
  return common.Address{}, err
 }

 gasLimit := uint64(3000000)
 gasPrice, err := d.client.SuggestGasPrice(context.Background())
 if err != nil {
  return common.Address{}, err
 }

 tx := types.NewContractCreation(auth_nonce, gasLimit, gasPrice, bytecode)
 signedTx, err := types.SignTx(tx, types.NewEIP155Signer(d.client.ChainID()), d.account.PrivateKey)
 if err != nil {
  return common.Address{}, err
 }

 err = d.client.SendTransaction(context.Background(), signedTx)
 if err != nil {
  return common.Address{}, err
 }

 fmt.Printf("Contract %s deployed to address %s\n", contractName, signedTx.ContractAddress())
 return signedTx.ContractAddress(), nil
}

func (d *Decent) CallContract(contractAddress common.Address, functionName string, args ...interface{}) ([]byte, error) {
 input, err := d.encodeInput(functionName, args...)
 if err != nil {
  return nil, err
 }

 msg := ethereum.CallMsg{
  To:   contractAddress,
  Data: input,
 }

 result, err := d.client.CallContract(context.Background(), msg, nil)
 if err != nil {
  return nil, err
 }

 return result, nil
}

func main() {
 decent, err := NewDecent("https://mainnet.infura.io/v3/YOUR_PROJECT_ID", "0x Your_Private_Key")
 if err != nil {
  log.Fatal(err)
 }

 contractAddress, err := decent.DeployContract("MyContract", []byte{0x60, 0x60})
 if err != nil {
  log.Fatal(err)
 }

 result, err := decent.CallContract(contractAddress, "-myFunction", uint256(42))
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("Result:", result)
}