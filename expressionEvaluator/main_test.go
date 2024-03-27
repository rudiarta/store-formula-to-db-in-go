package main

import "testing"

func Test_getRule(t *testing.T) {
	type args struct {
		data    ChargeableWeightData
		formula string
	}
	tests := []struct {
		name                 string
		args                 args
		wantChargeAbleWeight float64
		wantVolumeWeight     float64
		wantErr              bool
	}{
		{
			name: "Default Rule Formula",
			args: args{
				data: ChargeableWeightData{
					GrossWeight:  32.2,
					VolumeWeight: 2,
					ProductType:  PRODUCTA,
				},
				formula: DEFAULT,
			},
			wantChargeAbleWeight: 33,
			wantVolumeWeight:     2,
		},
		{
			name: "rules-0.31 Formula",
			args: args{
				data: ChargeableWeightData{
					GrossWeight:  32.2,
					VolumeWeight: 2,
					ProductType:  PRODUCTA,
				},
				formula: RULE_0_31,
			},
			wantChargeAbleWeight: 32,
			wantVolumeWeight:     2,
		},
		{
			name: "rules-0.01 Formula",
			args: args{
				data: ChargeableWeightData{
					GrossWeight:  32.2,
					VolumeWeight: 2,
					ProductType:  PRODUCTA,
				},
				formula: RULE_0_01,
			},
			wantChargeAbleWeight: 33,
			wantVolumeWeight:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChargeAbleWeight, gotVolumeWeight, err := getRule(tt.args.data, tt.args.formula)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotChargeAbleWeight != tt.wantChargeAbleWeight {
				t.Errorf("getRule() gotChargeAbleWeight = %v, want %v", gotChargeAbleWeight, tt.wantChargeAbleWeight)
			}
			if gotVolumeWeight != tt.wantVolumeWeight {
				t.Errorf("getRule() gotVolumeWeight = %v, want %v", gotVolumeWeight, tt.wantVolumeWeight)
			}
		})
	}
}
