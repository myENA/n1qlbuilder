package n1qlbuilder

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	b := Query().
		Select("a", "b", "c").
		From("mybucket").
		Where("a>b").
		And("c<d").
		Or("e is null").
		Limit(1)

	t.Logf("I built this query: %s", b)
}

func TestBuilderUpdate(t *testing.T) {
	b := Query().
		Update("trust-backup").
		Set("profile=$3").
		Where("type=$1").
		And("profile=$2").
		Returning("`trust-backup`.credentials.access_key")
	t.Logf("I built this query: %s", b)
}

