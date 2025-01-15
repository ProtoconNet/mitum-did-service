package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-did-registry/operation/did"
	"github.com/ProtoconNet/mitum-did-registry/state"
	"github.com/ProtoconNet/mitum-did-registry/types"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit

	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.DataHint, Instance: types.Data{}},
	{Hint: types.DocumentHint, Instance: types.Document{}},

	{Hint: did.CreateDIDHint, Instance: did.CreateDID{}},
	{Hint: did.MigrateDIDHint, Instance: did.MigrateDID{}},
	{Hint: did.ReactivateDIDHint, Instance: did.ReactivateDID{}},
	{Hint: did.DeactivateDIDHint, Instance: did.DeactivateDID{}},
	{Hint: did.RegisterModelHint, Instance: did.RegisterModel{}},
	{Hint: did.MigrateDIDItemHint, Instance: did.MigrateDIDItem{}},
	{Hint: state.DataStateValueHint, Instance: state.DataStateValue{}},
	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.DocumentStateValueHint, Instance: state.DocumentStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: did.CreateDIDFactHint, Instance: did.CreateDIDFact{}},
	{Hint: did.MigrateDIDFactHint, Instance: did.MigrateDIDFact{}},
	{Hint: did.ReactivateDIDFactHint, Instance: did.ReactivateDIDFact{}},
	{Hint: did.DeactivateDIDFactHint, Instance: did.DeactivateDIDFact{}},
	{Hint: did.RegisterModelFactHint, Instance: did.RegisterModelFact{}},
}

func init() {
	Hinters = append(Hinters, currencycmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, currencycmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, AddedSupportedHinters...)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
