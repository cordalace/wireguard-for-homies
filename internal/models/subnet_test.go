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

func TestSubnet_ToJSON(t *testing.T) {
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
				ID:   uuid.MustParse("1a4c0f86-07fa-4fdc-9f97-349925edb975"),
				CIDR: ipNetMustParse("192.168.0.1/24"),
			},
			want:    `{"id":"1a4c0f86-07fa-4fdc-9f97-349925edb975","cidr":"192.168.0.0/24"}`,
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
				data: []byte(`{"id":"1a4c0f86-07fa-4fdc-9f97-349925edb975","cidr":"192.168.0.0/24"}`),
			},
			want: &Subnet{
				ID:   uuid.MustParse("1a4c0f86-07fa-4fdc-9f97-349925edb975"),
				CIDR: ipNetMustParse("192.168.0.1/24"),
			},
			wantErr: false,
		},
		{
			name: "ipv6",
			args: args{
				data: []byte(`{"id":"1a4c0f86-07fa-4fdc-9f97-349925edb975","cidr":"2002:0:0:1234::/64"}`),
			},
			want: &Subnet{
				ID:   uuid.MustParse("1a4c0f86-07fa-4fdc-9f97-349925edb975"),
				CIDR: ipNetMustParse("2002:0:0:1234::/64"),
			},
			wantErr: false,
		},
		{
			name: "bad id",
			args: args{
				data: []byte(`{"id":"bad uuid","cidr":"2002:0:0:1234::/64"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad id json type",
			args: args{
				data: []byte(`{"id":42,"cidr":"2002:0:0:1234::/64"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad cidr",
			args: args{
				data: []byte(`{"id":"1a4c0f86-07fa-4fdc-9f97-349925edb975","cidr":"bad cidr"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad cidr json type",
			args: args{
				data: []byte(`{"id":"1a4c0f86-07fa-4fdc-9f97-349925edb975","cidr":42}`),
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
