/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fab

import (
	reqContext "context"

	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// ProposalProcessor simulates transaction proposal, so that a client can submit the result for ordering.
type ProposalProcessor interface {
	ProcessTransactionProposal(reqContext.Context, ProcessProposalRequest) (*TransactionProposalResponse, error)
}

// TxnHeaderOptions contains options for creating a Transaction Header
type TxnHeaderOptions struct {
	Nonce   []byte
	Creator []byte
}

// TxnHeaderOpt is a Transaction Header option
type TxnHeaderOpt func(*TxnHeaderOptions)

// WithNonce specifies the nonce to use when creating the Transaction Header
func WithNonce(nonce []byte) TxnHeaderOpt {
	return func(options *TxnHeaderOptions) {
		options.Nonce = nonce
	}
}

// WithCreator specifies the creator to use when creating the Transaction Header
func WithCreator(creator []byte) TxnHeaderOpt {
	return func(options *TxnHeaderOptions) {
		options.Creator = creator
	}
}

// ProposalSender provides the ability for a transaction proposal to be created and sent.
type ProposalSender interface {
	CreateTransactionHeader(opts ...TxnHeaderOpt) (TransactionHeader, error)
	SendTransactionProposal(*TransactionProposal, []ProposalProcessor) ([]*TransactionProposalResponse, error)
}

// TransactionID provides the identifier of a Fabric transaction proposal.
type TransactionID string

// EmptyTransactionID represents a non-existing transaction (usually due to error).
const EmptyTransactionID = TransactionID("")

// SystemChannel is the Fabric channel for managaing resources.
const SystemChannel = ""

// TransactionHeader provides a handle to transaction metadata.
type TransactionHeader interface {
	TransactionID() TransactionID
	Creator() []byte
	Nonce() []byte
	ChannelID() string
}

// ChaincodeInvokeRequest contains the parameters for sending a transaction proposal.
// nolint: maligned
type ChaincodeInvokeRequest struct {
	ChaincodeID  string
	Lang         pb.ChaincodeSpec_Type
	TransientMap map[string][]byte
	Fcn          string
	Args         [][]byte
	IsInit       bool
}

//RequestType shows which PAC transaction to create
type RequestType int32

const (
	//PrepareTxRequest is set if client wants generate & submit PrepareTx
	PrepareTxRequest RequestType = 1
	//DecideTxRequest is set if client wants generate & submit DecideTx
	DecideTxRequest RequestType = 2
	//AbortTxRequest is set if client wants generate & submit AbortTx
	AbortTxRequest RequestType = 3
)

//ClientData contains additional transaction data for private atomic commmit
type ClientData struct {
	//RequestedTransaction shows which PAC transaction to create
	RequestedTransaction RequestType
	//ValidationData contains marshalled specific transaction data that is
	//expected from the client due to RequestedTransaction type is set
	ValidationData []byte
}

// TransactionProposal contains a marashalled transaction proposal.
type TransactionProposal struct {
	TxnID TransactionID
	*pb.Proposal
	PACClientData ClientData
}

// ProcessProposalRequest requests simulation of a proposed transaction from transaction processors.
type ProcessProposalRequest struct {
	SignedProposal *pb.SignedProposal
}

// TransactionProposalResponse respresents the result of transaction proposal processing.
type TransactionProposalResponse struct {
	Endorser string
	// Status is the EndorserStatus
	Status int32
	// ChaincodeStatus is the status returned by Chaincode
	ChaincodeStatus int32
	*pb.ProposalResponse
}
