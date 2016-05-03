package timeago

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_FromDuration(t *testing.T) {
	cases := map[string]string{
		"0s":            "less than a minute",
		"29s":           "less than a minute",
		"30s":           "1 minute",
		"1m29s":         "1 minute",
		"1m30s":         "2 minutes",
		"44m29s":        "44 minutes",
		"44m30s":        "about 1 hour",
		"89m29s":        "about 1 hour",
		"89m30s":        "about 2 hours",
		"23h59m29s":     "about 24 hours",
		"23h59m30s":     "1 day",
		"41h59m29s":     "1 day",
		"41h59m30s":     "2 days",
		"2D12h":         "3 days",
		"29D23h59m29s":  "30 days",
		"29D23h59m30s":  "about 1 month",
		"44D23h59m29s":  "about 1 month",
		"44D23h59m30s":  "about 2 months",
		"59D23h59m29s":  "about 2 months",
		"59D23h59m30s":  "2 months",
		"364D23h59m29s": "12 months",
		"364D23h59m30s": "about 1 year",
		"1Y89D":         "about 1 year",
		"1Y6M":          "over 1 year",
		"1Y276D":        "almost 2 years",
		"2Y89D":         "about 2 years",
		"2Y90D":         "over 2 years",
		"2Y269D":        "over 2 years",
		"2Y9M":          "almost 3 years",
		"4Y276D":        "almost 5 years",
		"5Y89D":         "about 5 years",
		"5Y3M":          "over 5 years",
		"5Y8M29D":       "over 5 years",
		"5Y9M":          "almost 6 years",
		"9Y9M":          "almost 10 years",
		"10Y89D":        "about 10 years",
		"10Y3M":         "over 10 years",
		"10Y8M29D":      "over 10 years",
		"10Y9M":         "almost 11 years",
	}

	for input, expected := range cases {
		d := parseDuration(t, input)

		if v := FromDuration(d); v != expected {
			t.Fatalf("for %#v, expected %#v, but got %#v", input, expected, v)
		}
	}
}

func Test_FromTime(t *testing.T) {
	expected := "less than a minute ago"
	if v := FromTime(time.Now()); v != expected {
		t.Fatalf("expected %#v, but got %#v", expected, v)
	}

	futureCases := map[string]string{
		"29s":   "less than a minute from now",
		"45m":   "about 1 hour from now",
		"2D12h": "3 days from now",
	}

	for input, expected := range futureCases {
		d := parseDuration(t, input)
		if v := FromTime(time.Now().Add(d)); v != expected {
			t.Fatalf("for %#v, expected %#v, but got %#v", input, expected, v)
		}
	}

	pastCases := map[string]string{
		"29s":   "less than a minute ago",
		"45m":   "about 1 hour ago",
		"2D12h": "3 days ago",
	}

	for input, expected := range pastCases {
		d := parseDuration(t, input)
		if v := FromTime(time.Now().Add(-d)); v != expected {
			t.Fatalf("for %#v, expected %#v, but got %#v", input, expected, v)
		}
	}
}

func parseDuration(t *testing.T, s string) time.Duration {
	s, years := splitDuration(t, s, "Y")
	s, months := splitDuration(t, s, "M")
	s, days := splitDuration(t, s, "D")

	minutes := years*year + months*month + days*day

	add, err := time.ParseDuration(fmt.Sprintf("%0dm", minutes))
	if err != nil {
		t.Fatal(err)
	}

	if s == "" {
		return add
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		t.Fatal(err)
	}

	return d + add
}

func splitDuration(t *testing.T, s string, delim string) (string, int) {
	elems := strings.Split(s, delim)
	if len(elems) == 2 {
		i64, err := strconv.ParseInt(elems[0], 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		return elems[len(elems)-1], int(i64)
	}

	return s, 0
}
