package queue

import (
	"github.com/konsumgandalf/mpga-protoype-david/pkg/nexus/common/structs"
	"testing"
	"time"
)

func TestPriorityQueue(t *testing.T) {
	mockPriorityQueue := NewDeadlineQueue()

	startTime := time.Now()

	mockItemList := []*DeadlineItem{
		{
			BaseNexusItem: structs.BaseNexusItem{ID: "1", Index: 0},
			Deadline:      startTime,
		},
		{
			BaseNexusItem: structs.BaseNexusItem{ID: "2", Index: 1},
			Deadline:      startTime.Add(20 * time.Second),
		},
	}
	firstItem := mockItemList[0]

	// Test Push
	for _, item := range mockItemList {
		mockPriorityQueue.Push(item)
	}
	if mockPriorityQueue.Len() != 2 {
		t.Errorf("Expected length 2, got %d", mockPriorityQueue.Len())
	}

	// Test Peek
	if mockPriorityQueue.Peek() != firstItem {
		t.Errorf("Expected to peek item1, got different item")
	}

	// Test Update
	newDeadline1 := mockPriorityQueue.Peek().Deadline.Add(40 * time.Minute)
	mockPriorityQueue.Update(mockPriorityQueue.Peek(), newDeadline1)
	if !mockPriorityQueue.Peek().Deadline.Equal(newDeadline1) {
		t.Log("Correctly updated deadline of item1, now item2 is the peek item")
	} else {
		t.Errorf("Expected to peek item2, but got item1")
	}

	newDeadline2 := mockPriorityQueue.Peek().Deadline.Add(40 * time.Minute)
	mockPriorityQueue.Update(mockPriorityQueue.Peek(), newDeadline2)

	// Test Pop
	popped := mockPriorityQueue.Pop()
	if popped != firstItem {
		t.Errorf("Expected to pop item1, got different item")
	}
	if mockPriorityQueue.Len() != 1 {
		t.Errorf("Expected length 1, got %d", mockPriorityQueue.Len())
	}

	// Test Remove
	mockPriorityQueue.Remove("2")
	if mockPriorityQueue.Len() != 0 {
		t.Errorf("Expected length 0, got %d", mockPriorityQueue.Len())
	}
}

func TestDeadlineImpl(t *testing.T) {
	mockDeadlineHeap := &deadlineHeap{}

	// Test Len
	if mockDeadlineHeap.Len() != 0 {
		t.Errorf("Expected length 0, got %d", mockDeadlineHeap.Len())
	}

	// Test Push
	startTime := time.Now()

	mockItemList := []*DeadlineItem{
		{
			BaseNexusItem: structs.BaseNexusItem{ID: "1", Index: 0},
			Deadline:      startTime,
		},
		{
			BaseNexusItem: structs.BaseNexusItem{ID: "2", Index: 1},
			Deadline:      startTime.Add(20 * time.Second),
		},
	}
	firstItem := mockItemList[0]

	// Test Push
	for _, item := range mockItemList {
		mockDeadlineHeap.Push(item)
	}
	if mockDeadlineHeap.Len() != 2 {
		t.Errorf("Expected length 2, got %d", mockDeadlineHeap.Len())
	}

	// Test Less
	if !mockDeadlineHeap.Less(0, 1) {
		t.Errorf("Expected item1 to be less than item2")
	}

	// Test Swap
	mockDeadlineHeap.Swap(0, 1)
	if mockDeadlineHeap.Less(0, 1) {
		t.Errorf("Expected item2 to be less than item1 after swap")
	}

	// Test Pop
	popped := mockDeadlineHeap.Pop().(*DeadlineItem)
	if popped != firstItem {
		t.Errorf("Expected to pop item1, got different item")
	}
	if mockDeadlineHeap.Len() != 1 {
		t.Errorf("Expected length 1, got %d", mockDeadlineHeap.Len())
	}
}
