package did

import (
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

type DIDItem interface {
	util.Byter
	util.IsValider
	Currency() currencytypes.CurrencyID
}

var MigrateDIDItems uint = 100

var (
	MigrateDIDFactHint = hint.MustNewHint("mitum-did-migrate-did-operation-fact-v0.0.1")
	MigrateDIDHint     = hint.MustNewHint("mitum-did-migrate-did-operation-v0.0.1")
)

type MigrateDIDFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []MigrateDIDItem
}

func NewMigrateDIDFact(
	token []byte, sender mitumbase.Address, items []MigrateDIDItem) MigrateDIDFact {
	bf := mitumbase.NewBaseFact(MigrateDIDFactHint, token)
	fact := MigrateDIDFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MigrateDIDFact) IsValid(b []byte) error {
	if n := len(fact.items); n < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items")))
	} else if n > int(MigrateDIDItems) {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("items, %d over max, %d", n, MigrateDIDItems)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if it.contract.Equal(fact.sender) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		k := fmt.Sprintf("%s-%s", it.contract, it.pubKey)

		if _, found := founds[k]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("pubKey %v for contract account %v", it.PubKeyReformed(), it.Contract())))
		}

		founds[k] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact MigrateDIDFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact MigrateDIDFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MigrateDIDFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact MigrateDIDFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact MigrateDIDFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact MigrateDIDFact) Items() []MigrateDIDItem {
	return fact.items
}

func (fact MigrateDIDFact) Addresses() ([]mitumbase.Address, error) {
	var as []mitumbase.Address

	adrMap := make(map[string]struct{})
	for i := range fact.items {
		for j := range fact.items[i].Addresses() {
			if _, found := adrMap[fact.items[i].Addresses()[j].String()]; !found {
				adrMap[fact.items[i].Addresses()[j].String()] = struct{}{}
				as = append(as, fact.items[i].Addresses()[j])
			}
		}
	}
	as = append(as, fact.sender)

	return as, nil
}

type MigrateDID struct {
	common.BaseOperation
}

func NewMigrateDID(fact MigrateDIDFact) (MigrateDID, error) {
	return MigrateDID{BaseOperation: common.NewBaseOperation(MigrateDIDHint, fact)}, nil
}
