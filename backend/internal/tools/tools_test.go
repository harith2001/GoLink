package tools

import "testing"

func TestCheckCompareCount(t *testing.T) {
	if checkCompareCount(2) != nil || checkCompareCount(3) != nil {
		t.Fatal("2 and 3 are valid")
	}
	if checkCompareCount(1) == nil || checkCompareCount(4) == nil {
		t.Fatal("only 2 or 3")
	}
}

func TestTrendMonths(t *testing.T) {
	if trendMonths(nil) != 12 {
		t.Fatal("default 12")
	}
	n := 200
	if trendMonths(&n) != 120 {
		t.Fatal("cap 120")
	}
	n = 6
	if trendMonths(&n) != 6 {
		t.Fatal("keep 6")
	}
}
