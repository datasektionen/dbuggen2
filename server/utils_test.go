package server

import "testing"

func TestDetermineOrder(t *testing.T) {
	t.Run("valid positive array", func(t *testing.T) {
		orderString := "1,2,3,4,5"
		expected := []int{1, 2, 3, 4, 5}
		got, err := determineOrder(orderString)
		if err != nil {
			t.Errorf("encountered error %v", err)
		}
		for i, g := range got {
			if g != expected[i] {
				t.Errorf("expected array %v, but got %v", expected, got)
			}
		}
	})

	t.Run("valid mixed array", func(t *testing.T) {
		orderString := "-1,2,5,6543,-4"
		expected := []int{-1, 2, 5, 6543, -4}
		got, err := determineOrder(orderString)
		if err != nil {
			t.Errorf("encountered error %v", err)
		}
		for i, g := range got {
			if g != expected[i] {
				t.Errorf("expected array %v, but got %v", expected, got)
			}
		}
	})

	t.Run("invalid array with spaces", func(t *testing.T) {
		orderString := "-1,d2, 3"
		got, err := determineOrder(orderString)
		if err == nil {
			t.Errorf("should have encountered an error, got %v instead", got)
		}
	})
}
