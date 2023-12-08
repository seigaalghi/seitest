package tested

import (
	"testing"

	"github.com/seigaalghi/seitest/interfaces"
)

func TestNewTested(t *testing.T) {
	type payload struct {
		s interfaces.Sample
	}
	tests := []struct {
		scenario string
		payload  payload
		assert   func(t *testing.T, res0 Tested)
	}{
		// Put Your Scenario Here
	}
	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			res0 := NewTested(tt.payload.s)
			tt.assert(t, res0)
		})
	}

}

func TestJamban(t *testing.T) {
	type payload struct {
		s string
	}
	tests := []struct {
		scenario string
		payload  payload
		assert   func(t *testing.T, res0 string)
	}{
		// Put Your Scenario Here
	}
	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			res0 := Jamban(tt.payload.s)
			tt.assert(t, res0)
		})
	}

}
