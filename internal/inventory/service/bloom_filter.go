package service

import (
	"strconv"

	"github.com/bits-and-blooms/bloom/v3"
)

const (
	bloomN = 100000
	bloomP = 0.01
)

type productBloomFilter struct {
	filter *bloom.BloomFilter
}

func newProductBloomFilter() *productBloomFilter {
	return &productBloomFilter{
		filter: bloom.NewWithEstimates(bloomN, bloomP),
	}
}

func (b *productBloomFilter) add(id int64) {
	b.filter.AddString(strconv.FormatInt(id, 10))
}

func (b *productBloomFilter) addAll(ids []int64) {
	for _, id := range ids {
		b.add(id)
	}
}

func (b *productBloomFilter) mayExist(id int64) bool {
	return b.filter.TestString(strconv.FormatInt(id, 10))
}

func (b *productBloomFilter) clear() {
	b.filter = bloom.NewWithEstimates(bloomN, bloomP)
}