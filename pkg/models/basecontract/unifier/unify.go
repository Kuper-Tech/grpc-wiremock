package unifier

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/afero"

	contract "github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/unifier/openapi"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/models/basecontract/unifier/proto"
	"github.com/SberMarket-Tech/grpc-wiremock/pkg/sourcer/types"
)

type unifier interface {
	Unify(context.Context, contract.SetOfContracts, string) error
}

type compiler interface {
	CompileToOpenAPI(context.Context, contract.Contract, string, io.Writer) error
}

type baseUnifier struct {
	fs afero.Fs

	compiler compiler

	logs io.Writer
}

func NewUnifier(fs afero.Fs, compiler compiler, logs io.Writer) *baseUnifier {
	return &baseUnifier{fs: fs, logs: logs, compiler: compiler}
}

func (u *baseUnifier) Unify(ctx context.Context, contracts contract.SetOfContracts, path string) error {
	specificUnifier, err := u.getUnifier(u.fs, contracts, u.compiler)
	if err != nil {
		return fmt.Errorf("get unifier: %w", err)
	}

	if err := specificUnifier.Unify(ctx, contracts, path); err != nil {
		return fmt.Errorf("unify: %w", err)
	}

	return nil
}

func (u *baseUnifier) getUnifier(fs afero.Fs, contracts contract.SetOfContracts, compiler compiler) (unifier, error) {
	switch contracts.FileType() {
	case types.ProtoType:
		return proto.NewUnifier(fs, compiler, u.logs), nil

	case types.OpenAPIType:
		return openapi.NewUnifier(fs), nil
	}

	return nil, fmt.Errorf("no available unifier")
}
