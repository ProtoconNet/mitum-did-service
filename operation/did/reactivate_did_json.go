package did

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type ReactivateDIDFactJSONMarshaler struct {
	mitumbase.BaseFactJSONMarshaler
	Sender   mitumbase.Address `json:"sender"`
	Contract mitumbase.Address `json:"contract"`
	DID      string            `json:"did"`
	Currency types.CurrencyID  `json:"currency"`
}

func (fact ReactivateDIDFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(ReactivateDIDFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		DID:                   fact.did,
		Currency:              fact.currency,
	})
}

type ReactivateDIDFactJSONUnmarshaler struct {
	mitumbase.BaseFactJSONUnmarshaler
	Sender   string `json:"sender"`
	Contract string `json:"contract"`
	DID      string `json:"did"`
	Currency string `json:"currency"`
}

func (fact *ReactivateDIDFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u ReactivateDIDFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, u.Sender, u.Contract, u.DID, u.Currency); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op ReactivateDID) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		op.BaseOperation.JSONMarshaler(),
	)
}

func (op *ReactivateDID) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
