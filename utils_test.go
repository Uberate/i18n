package i18n

import (
	"reflect"
	"testing"
)

func TestFillStrings(t *testing.T) {
	type args struct {
		v      []string
		length int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "test nil v", args: args{v: nil, length: 2}, want: []string{"", ""}},
		{name: "test emtpy v", args: args{v: []string{}, length: 2}, want: []string{"", ""}},
		{name: "test some value", args: args{v: []string{"test"}, length: 2}, want: []string{"test", ""}},
		{name: "test equals value", args: args{v: []string{"test", "test"}, length: 2}, want: []string{"test", "test"}},
		{name: "test bigger value", args: args{v: []string{"t", "t", "t"}, length: 2}, want: []string{"t", "t", "t"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FillStrings(tt.args.v, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringsMapperCopy(t *testing.T) {
	type args struct {
		origin       []string
		target       []string
		mapper       map[int]int
		originOffset int
		targetOffset int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "test base", args: args{
			origin:       []string{"a", "b", "c"},
			target:       []string{"", "", ""},
			mapper:       map[int]int{0: 1, 1: 2, 2: 0},
			originOffset: 0,
			targetOffset: 0,
		}, want: []string{"c", "a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringsMapperCopy(tt.args.origin, tt.args.target, tt.args.mapper, tt.args.originOffset, tt.args.targetOffset)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringsMapperCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringsMapperCopy() got = %v, want %v", got, tt.want)
			}
		})
	}
}
