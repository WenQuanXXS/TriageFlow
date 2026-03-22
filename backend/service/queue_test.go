package service

import (
	"testing"
	"time"

	"github.com/triageflow/backend/model"
)

func TestComputeQueueOrder_PriorityOrdering(t *testing.T) {
	now := time.Now()

	urgentOrder := ComputeQueueOrder("urgent", now)
	highOrder := ComputeQueueOrder("high", now)
	normalOrder := ComputeQueueOrder("normal", now)

	if urgentOrder >= highOrder {
		t.Errorf("urgent (%d) should be < high (%d)", urgentOrder, highOrder)
	}
	if highOrder >= normalOrder {
		t.Errorf("high (%d) should be < normal (%d)", highOrder, normalOrder)
	}
}

func TestComputeQueueOrder_SamePriorityByTime(t *testing.T) {
	earlier := time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC)
	later := time.Date(2026, 3, 21, 10, 5, 0, 0, time.UTC)

	earlierOrder := ComputeQueueOrder("normal", earlier)
	laterOrder := ComputeQueueOrder("normal", later)

	if earlierOrder >= laterOrder {
		t.Errorf("earlier (%d) should be < later (%d)", earlierOrder, laterOrder)
	}
}

func TestComputeQueueOrder_HigherPriorityBeatsEarlierTime(t *testing.T) {
	// An urgent patient arriving later should still sort before a normal patient arriving earlier.
	earlier := time.Date(2026, 3, 21, 8, 0, 0, 0, time.UTC)
	later := time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC)

	urgentLater := ComputeQueueOrder("urgent", later)
	normalEarlier := ComputeQueueOrder("normal", earlier)

	if urgentLater >= normalEarlier {
		t.Errorf("urgent-later (%d) should be < normal-earlier (%d)", urgentLater, normalEarlier)
	}
}

func TestPriorityWeight(t *testing.T) {
	cases := []struct {
		priority string
		weight   int
	}{
		{"urgent", 0},
		{"high", 1},
		{"normal", 2},
		{"unknown", 3},
		{"", 3},
	}

	for _, tc := range cases {
		w := model.PriorityWeight(tc.priority)
		if w != tc.weight {
			t.Errorf("PriorityWeight(%q) = %d, want %d", tc.priority, w, tc.weight)
		}
	}
}
