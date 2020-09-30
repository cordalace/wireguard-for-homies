package models

import (
	"net"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func ipNetMustParse(s string) *net.IPNet {
	_, ret, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return ret
}

func TestSubnetToJSON(t *testing.T) {
	type fields struct {
		ID   uuid.UUID
		CIDR *net.IPNet
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "ipv4",
			fields: fields{
				ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
				CIDR: ipNetMustParse("192.168.0.0/24"),
			},
			want:    `{"id":"2eba5e83-d0c3-46f0-bbeb-884e62e19b62","cidr":"192.168.0.0/24"}`,
			wantErr: false,
		},
		{
			name: "ipv6",
			fields: fields{
				ID:   uuid.MustParse("1a4c0f86-07fa-4fdc-9f97-349925edb975"),
				CIDR: ipNetMustParse("2002:0:0:1234::/64"),
			},
			want:    `{"id":"1a4c0f86-07fa-4fdc-9f97-349925edb975","cidr":"2002:0:0:1234::/64"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			s := &Subnet{
				ID:   tt.fields.ID,
				CIDR: tt.fields.CIDR,
			}
			got, err := s.ToJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Subnet.ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Subnet.ToJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestSubnetFromJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Subnet
		wantErr bool
	}{
		{
			name: "ipv4",
			args: args{
				data: []byte(`{"id":"506ca6b1-a2da-4396-b9ba-304c032fc401","cidr":"192.168.1.0/24"}`),
			},
			want: &Subnet{
				ID:   uuid.MustParse("506ca6b1-a2da-4396-b9ba-304c032fc401"),
				CIDR: ipNetMustParse("192.168.1.0/24"),
			},
			wantErr: false,
		},
		{
			name: "ipv6",
			args: args{
				data: []byte(`{"id":"2a782886-2661-45ce-bb6f-106b429ac76e","cidr":"2002:0:0:1235::/64"}`),
			},
			want: &Subnet{
				ID:   uuid.MustParse("2a782886-2661-45ce-bb6f-106b429ac76e"),
				CIDR: ipNetMustParse("2002:0:0:1235::/64"),
			},
			wantErr: false,
		},
		{
			name: "bad id",
			args: args{
				data: []byte(`{"id":"bad uuid","cidr":"2002:0:0:1235::/64"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad id json type",
			args: args{
				data: []byte(`{"id":42,"cidr":"2002:0:0:1235::/64"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad cidr",
			args: args{
				data: []byte(`{"id":"bfc2ab7e-a77f-437d-81f5-ae463712fc88","cidr":"bad cidr"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad cidr json type",
			args: args{
				data: []byte(`{"id":"41e0a99b-4b9c-4d84-81d4-4aba1b0e2a61","cidr":42}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubnetFromJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubnetFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubnetFromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubnetListToJSON(t *testing.T) {
	tests := []struct {
		name    string
		l       SubnetList
		want    string
		wantErr bool
	}{
		{
			name: "one subnet",
			l: SubnetList{
				&Subnet{
					ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
					CIDR: ipNetMustParse("192.168.0.0/24"),
				},
			},
			want:    `[{"id":"2eba5e83-d0c3-46f0-bbeb-884e62e19b62","cidr":"192.168.0.0/24"}]`,
			wantErr: false,
		},
		{
			name:    "empty",
			l:       SubnetList{},
			want:    `[]`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.ToJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("SubnetList.ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("SubnetList.ToJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestSubnetListFromJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []*Subnet
		wantErr bool
	}{
		{
			name: "one subnet",
			args: args{
				data: []byte(`[{"id":"2eba5e83-d0c3-46f0-bbeb-884e62e19b62","cidr":"192.168.0.0/24"}]`),
			},
			want: []*Subnet{
				{
					ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
					CIDR: ipNetMustParse("192.168.0.0/24"),
				},
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				data: []byte(`[{"id":"2eba5e83-d0c3-46f0-bbeb-884e62e19b62","cidr":"192.168.0.0/24"}]`),
			},
			want: []*Subnet{
				{
					ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
					CIDR: ipNetMustParse("192.168.0.0/24"),
				},
			},
			wantErr: false,
		},
		{
			name: "bad json type",
			args: args{
				data: []byte(`42`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad id",
			args: args{
				data: []byte(`[{"id":"non uuid","cidr":"192.168.0.0/24"}]`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad cidr",
			args: args{
				data: []byte(`[{"id":"2eba5e83-d0c3-46f0-bbeb-884e62e19b62","cidr":"bad cidr"}]`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubnetListFromJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubnetListFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubnetListFromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
