package main

import (
	"reflect"
	"testing"
)

func TestMinCoins(t *testing.T) {
	tests := []struct {
		name          string
		amount        int
		denominations []int
		expected      int
	}{
		{"Zero amount", 0, []int{1, 5, 10, 25, 50}, 0},
		{"Only pennies needed", 4, []int{1, 5, 10, 25, 50}, 4},
		{"Exact denomination", 5, []int{1, 5, 10, 25, 50}, 1},
		{"Simple combination", 11, []int{1, 5, 10, 25, 50}, 2},
		{"Example 1", 87, []int{1, 5, 10, 25, 50}, 5},
		{"Example 2", 42, []int{1, 5, 10, 25, 50}, 5},
		{"Larger amount", 99, []int{1, 5, 10, 25, 50}, 8},
		{"Cannot make amount", 3, []int{5, 10, 25}, -1},
		{"Custom denominations", 30, []int{1, 5, 10}, 3},
		{"Larger denominations", 63, []int{1, 5, 10, 25}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinCoins(tt.amount, tt.denominations)
			if result != tt.expected {
				t.Errorf("MinCoins(%d, %v) = %d, expected %d",
					tt.amount, tt.denominations, result, tt.expected)
			}
		})
	}
}

func TestCoinCombination(t *testing.T) {
	tests := []struct {
		name          string
		amount        int
		denominations []int
		expected      map[int]int
	}{
		{"Zero amount", 0, []int{1, 5, 10, 25, 50}, map[int]int{}},
		{"Only pennies needed", 4, []int{1, 5, 10, 25, 50}, map[int]int{1: 4}},
		{"Exact denomination", 5, []int{1, 5, 10, 25, 50}, map[int]int{5: 1}},
		{"Simple combination", 11, []int{1, 5, 10, 25, 50}, map[int]int{1: 1, 10: 1}},
		{"Example 1", 87, []int{1, 5, 10, 25, 50}, map[int]int{1: 2, 10: 1, 25: 1, 50: 1}},
		{"Example 2", 42, []int{1, 5, 10, 25, 50}, map[int]int{1: 2, 5: 1, 10: 1, 25: 1}},
		{"Larger amount", 99, []int{1, 5, 10, 25, 50}, map[int]int{1: 4, 10: 2, 25: 1, 50: 1}},
		{"Cannot make amount", 3, []int{5, 10, 25}, map[int]int{}},
		{"Custom denominations", 30, []int{1, 6, 10}, map[int]int{10: 3}},
		{"Larger denominations", 63, []int{1, 5, 10, 25}, map[int]int{1: 3, 10: 1, 25: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CoinCombination(tt.amount, tt.denominations)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CoinCombination(%d, %v) = %v, expected %v",
					tt.amount, tt.denominations, result, tt.expected)
			}
		})
	}
}

func TestExampleCases(t *testing.T) {
	// Standard U.S. denominations
	denominations := []int{1, 5, 10, 25, 50}

	// Example 1
	if result := MinCoins(87, denominations); result != 5 {
		t.Errorf("MinCoins(87, %v) = %d, expected 5", denominations, result)
	}

	// Example 2
	expectedCombination1 := map[int]int{1: 2, 10: 1, 25: 1, 50: 1}
	if result := CoinCombination(87, denominations); !reflect.DeepEqual(result, expectedCombination1) {
		t.Errorf("CoinCombination(87, %v) = %v, expected %v",
			denominations, result, expectedCombination1)
	}

	// Example 3
	if result := MinCoins(42, denominations); result != 5 {
		t.Errorf("MinCoins(42, %v) = %d, expected 5", denominations, result)
	}

	// Example 4
	expectedCombination2 := map[int]int{1: 2, 5: 1, 10: 1, 25: 1}
	if result := CoinCombination(42, denominations); !reflect.DeepEqual(result, expectedCombination2) {
		t.Errorf("CoinCombination(42, %v) = %v, expected %v",
			denominations, result, expectedCombination2)
	}
}
