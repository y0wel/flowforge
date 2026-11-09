package steps

import (
	"context"
	"testing"
)

func TestValidateStep_Execute(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		step    *ValidateStep
		wantErr bool
	}{
		{
			name: "all required fields present",
			data: map[string]any{"vendor": "Acme", "total": 42.5},
			step: &ValidateStep{
				Name:     "validate",
				Required: []string{"vendor", "total"},
				Numeric:  []string{"total"},
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			data: map[string]any{"vendor": "Acme"},
			step: &ValidateStep{
				Name:     "validate",
				Required: []string{"vendor", "total"},
			},
			wantErr: true,
		},
		{
			name: "required field is empty string",
			data: map[string]any{"vendor": ""},
			step: &ValidateStep{
				Name:     "validate",
				Required: []string{"vendor"},
			},
			wantErr: true,
		},
		{
			name: "numeric field has wrong type",
			data: map[string]any{"total": "not a number"},
			step: &ValidateStep{
				Name:    "validate",
				Numeric: []string{"total"},
			},
			wantErr: true,
		},
		{
			name: "numeric field missing is not an error by itself",
			data: map[string]any{},
			step: &ValidateStep{
				Name:    "validate",
				Numeric: []string{"total"},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := NewContext()
			data.Data = tc.data

			err := tc.step.Execute(context.Background(), data)

			if tc.wantErr && err == nil {
				t.Error("expected an error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
