package db

import "testing"

func BenchmarkXxx(b *testing.B) {
	db := InitStorage()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v, ok := db.LookupRead("odit")
		if !ok {
			db.Add("odit", 1)
			continue
		}
		switch num := v.(type) {
		case int:
			num += 1
			db.Add("odit", num)
		}
	}
}
