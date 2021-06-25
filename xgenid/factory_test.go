package xgenid

import (
	_ "github.com/duchiporexia/goutils/xlog"
	"testing"
)

func TestId(t *testing.T) {
	serverId = 1
	for i := 0; i <= 2; i++ {
		generateSingleId(t)
	}
}

func generateSingleId(t *testing.T) {
	rounds := 4096 * 1024
	unique := map[int64]struct{}{}

	for i := 0; i < rounds; i++ {
		id, _ := Id()
		unique[id] = struct{}{}
	}
	if v := len(unique); v != rounds {
		t.Errorf("expected %v; got %v\n", rounds, v)
		return
	}
}

func TestIdN(t *testing.T) {
	serverId = 1
	for i := 0; i <= 2; i++ {
		generateBatchIds(t)
	}
}

func generateBatchIds(t *testing.T) {
	n := 4096
	rounds := 1024
	allIds := make([][]int64, 0, rounds)

	for i := 0; i < rounds; i++ {
		ids, _ := IdN(n)
		allIds = append(allIds, ids)
	}

	unique := map[int64]struct{}{}
	for _, ids := range allIds {
		if v := len(ids); v != n {
			t.Errorf("expected %v; got %v\n", n, v)
		}
		for _, id := range ids {
			unique[id] = struct{}{}
		}
	}

	expected := n * rounds
	if v := len(unique); v != expected {
		t.Errorf("expected %v; got %v\n", expected, v)
		return
	}
}
