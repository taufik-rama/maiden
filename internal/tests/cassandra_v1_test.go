package tests

import (
	"testing"

	"github.com/taufik-rama/maiden/config.v1"
	"github.com/taufik-rama/maiden/storage.v1"
)

func TestCassandra(t *testing.T) {

	fixture := new(config.Fixture)
	if err := fixture.Parse("testdata/cassandra/cassandra.yaml"); err != nil {
		panic(err)
	}

	test := Cassandra{
		t: t,
		f: fixture,
	}

	test.push()
}

type Cassandra struct {
	t *testing.T
	f *config.Fixture
}

func (c Cassandra) push() {
	(storage.Cassandra{}).Push(c.f.Cassandra)
}
