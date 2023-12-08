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
		assert   func(*testing.T, Tested)
	}{}
	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			a := NewTested(tt.payload.s)
			tt.assert(t, a)
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
		assert   func(*testing.T, string)
	}{}
	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			a := Jamban(tt.payload.s)
			tt.assert(t, a)
		})
	}

}
