package tests

import (
	"testing"

	"github.com/taufik-rama/maiden/internal/config"
)

func TestConfig(t *testing.T) {
	{
		c, err := config.New("definitions/empty.yaml")
		if err != nil {
			t.Fatal(err)
		}
		if c.Fixtures != nil || c.Services != nil {
			t.Fail()
		}
	}
}
