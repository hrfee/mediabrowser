package mediabrowser

import (
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	expected := "2021-01-27T03:16:36"
	t1, _ := time.Parse(TimeFormat, expected)
	t1Bytes, _ := (&Time{t1}).MarshalJSON()
	if string(t1Bytes) != expected {
		t.Errorf("output != expected: %s != %s", string(t1Bytes), expected)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := map[string]time.Time{}
	tests["2021-01-27T03:16:36.28538ZZ"], _ = time.Parse(TimeFormat, "2021-01-27T03:16:36")
	tests["\"2021-01-27T03:16:36.28538ZZ\""], _ = time.Parse(TimeFormat, "2021-01-27T03:16:36")
	tests["\"2021-01-27T03:16:36.28538Z\""], _ = time.Parse(TimeFormat, "2021-01-27T03:16:36")
	tests["\"2021-01-09T20:58:41.5907920+00:00\""], _ = time.Parse(TimeFormat, "2021-01-09T20:58:41")
	for in, expected := range tests {
		parsed := Time{}
		err := parsed.UnmarshalJSON([]byte(in))
		if err != nil {
			t.Errorf("%s failed to parse with error: %v", in, err)
			continue
		}
		if parsed.Time != expected {
			t.Errorf("%s parsed incorrectly to %v, should have been %v", in, parsed.Time, expected)
		}
	}
}

func benchUnmarshalJSON(tests []string, b *testing.B) {
	t := Time{}
	for n := 0; n < b.N; n++ {
		for _, s := range tests {
			t.UnmarshalJSON([]byte(s))
		}
	}
}

func BenchmarkUnmarshalJSON1(b *testing.B) {
	benchUnmarshalJSON([]string{
		"\"2021-01-27T03:16:36.28538Z\"",
	}, b)
}

func BenchmarkUnmarshalJSON4(b *testing.B) {
	benchUnmarshalJSON([]string{
		"2021-01-27T03:16:36.28538ZZ",
		"\"2021-01-27T03:16:36.28538ZZ\"",
		"\"2021-01-27T03:16:36.28538Z\"",
		"\"2021-01-09T20:58:41.5907920+00:00\"",
	}, b)
}
