package BarbDB

import (
	"testing"
)

var testdb *barbDB

func init() {
	testdb, _ = OpenDB("barbtest.db")
}

func Test_barbDB_Set(t *testing.T) {

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		db      barbDB
		args    args
		wantErr bool
	}{
		{"st", *testdb, args{"k1", "v1"}, false},
		{"nd", *testdb, args{"k2", "v2"}, false},
		{"rd", *testdb, args{"k3", "v3"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("barbDB.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_barbDB_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		db      barbDB
		args    args
		want    string
		wantErr bool
	}{
		{"get ok", *testdb, args{"k1"}, "v1", false},
		{"get nok", *testdb, args{"k1000"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("barbDB.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("barbDB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
