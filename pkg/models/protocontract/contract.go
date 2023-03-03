package protocontract

import (
	"github.com/jhump/protoreflect/desc"

	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/utils/fsutils"
)

type Base struct {
	Package    string
	GoPackage  string
	HeaderPath string

	Messages     []string
	ImportsPaths []string

	Services []Service
}

// Contract represents contract files with scope.
//
//	`HeaderPath` stores main contract path.
//	`[]ImportsPaths` stores paths of contracts which was imported in the main contract file.
type Contract struct {
	Base

	Desc *desc.FileDescriptor
}

type MessageType struct {
	Name      string
	Package   string
	GoPackage string

	IsExternal bool
}

type Method struct {
	Name    string
	Package string

	InType  MessageType
	OutType MessageType

	MethodType types.ProtoMethodType
}

type Service struct {
	Name      string
	GoPackage string

	Methods []Method
}

func (c Contract) HasMethods() bool {
	for _, service := range c.Services {
		if len(service.Methods) > 0 {
			return true
		}
	}

	return false
}

func (c Contract) IsProtoContract() bool {
	return fsutils.GetFileExt(c.HeaderPath) == "proto"
}

// SetOfContracts is just alias for several contracts.
type SetOfContracts []Contract

// WithMethodsOnly returns only contracts with methods.
func (set SetOfContracts) WithMethodsOnly() SetOfContracts {
	var filtered SetOfContracts

	for _, contract := range set {
		if contract.HasMethods() {
			filtered = append(filtered, contract)
		}
	}

	return filtered
}

func (set SetOfContracts) HasContractsWithMethods() bool {
	for _, contract := range set {
		if contract.HasMethods() {
			return true
		}
	}

	return false
}

func (set SetOfContracts) HasContractsWithProto() bool {
	for _, contract := range set {
		if contract.IsProtoContract() {
			return true
		}
	}

	return false
}
