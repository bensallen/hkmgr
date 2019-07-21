package config

import (
	"net"
	"reflect"
	"testing"

	"github.com/bensallen/hkmgr/internal/network"
)

func TestTap_toBridge(t *testing.T) {
	type fields struct {
		Bridge  string
		IP      string
		Nat     bool
		NatIf   string
		PfRules []string
		DHCP    bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *network.Bridge
		wantErr bool
	}{
		{
			name: "bridge0",
			fields: fields{
				Bridge: "bridge0",
			},
			want: &network.Bridge{Device: "bridge0"},
		},
		{
			name: "bridge0 192.168.0.1",
			fields: fields{
				Bridge: "bridge0",
				IP:     "192.168.0.1",
			},
			want: &network.Bridge{Device: "bridge0", IP: net.ParseIP("192.168.0.1"), Netmask: net.IPMask{255, 255, 255, 0}},
		},
		{
			name: "bridge0 10.0.0.1/16",
			fields: fields{
				Bridge: "bridge0",
				IP:     "10.0.0.1/16",
			},
			want: &network.Bridge{Device: "bridge0", IP: net.ParseIP("10.0.0.1"), Netmask: net.IPMask{255, 255, 0, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &Tap{
				Bridge:  tt.fields.Bridge,
				IP:      tt.fields.IP,
				Nat:     tt.fields.Nat,
				NatIf:   tt.fields.NatIf,
				PfRules: tt.fields.PfRules,
				DHCP:    tt.fields.DHCP,
			}
			got, err := tap.toBridge()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tap.toBridge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tap.toBridge() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
