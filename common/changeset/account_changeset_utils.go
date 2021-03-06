package changeset

import (
	"bytes"
	"sort"

	"github.com/EVRICE/tgeth_alpha/common/dbutils"
	"github.com/EVRICE/tgeth_alpha/ethdb"
)

func findInAccountChangeSet(c ethdb.CursorDupSort, blockNumber uint64, key []byte, keyLen int) ([]byte, error) {
	fromDBFormat := FromDBFormat(keyLen)
	k, v, err := c.SeekBothRange(dbutils.EncodeBlockNumber(blockNumber), key)
	if err != nil {
		return nil, err
	}
	if k == nil {
		return nil, nil
	}
	_, k, v = fromDBFormat(k, v)
	if !bytes.HasPrefix(k, key) {
		return nil, nil
	}
	return v, nil
}

func encodeAccounts2(blockN uint64, s *ChangeSet, f func(k, v []byte) error) error {
	sort.Sort(s)
	newK := dbutils.EncodeBlockNumber(blockN)
	for _, cs := range s.Changes {
		newV := make([]byte, len(cs.Key)+len(cs.Value))
		copy(newV, cs.Key)
		copy(newV[len(cs.Key):], cs.Value)
		if err := f(newK, newV); err != nil {
			return err
		}
	}
	return nil
}
