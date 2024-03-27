package main

import (
	"testing"
)

func Test_getRule(t *testing.T) {
	type args struct {
		rules   string
		reqData ChargeableWeightData
	}
	tests := []struct {
		name                 string
		args                 args
		wantChargeAbleWeight float64
		wantVolumeWeight     float64
		wantErr              error
	}{
		{
			name: "Default Rules",
			args: args{
				rules: DEFAULT,
				reqData: ChargeableWeightData{
					GrossWeight:  32.2,
					VolumeWeight: 2,
					ProductType:  PRODUCTA,
				},
			},
			wantChargeAbleWeight: 33,
			wantVolumeWeight:     2,
		},
		{
			name: "rules-0.31 Rules",
			args: args{
				rules: RULE_0_31,
				reqData: ChargeableWeightData{
					GrossWeight:  32.2,
					VolumeWeight: 2,
					ProductType:  PRODUCTA,
				},
			},
			wantChargeAbleWeight: 32,
			wantVolumeWeight:     2,
		},
		{
			name: "rules-0.01 Rules",
			args: args{
				rules: RULE_0_01,
				reqData: ChargeableWeightData{
					GrossWeight:  32.2,
					VolumeWeight: 2,
					ProductType:  PRODUCTA,
				},
			},
			wantChargeAbleWeight: 33,
			wantVolumeWeight:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := getRule(tt.args.rules)

			gotChargeAbleWeight, gotVolumeWeight := rule(tt.args.reqData)
			if gotChargeAbleWeight != tt.wantChargeAbleWeight {
				t.Errorf("getRule() gotChargeAbleWeight = %v, want %v", gotChargeAbleWeight, tt.wantChargeAbleWeight)
			}
			if gotVolumeWeight != tt.wantVolumeWeight {
				t.Errorf("getRule() gotVolumeWeight = %v, want %v", gotVolumeWeight, tt.wantVolumeWeight)
			}
		})
	}
}
