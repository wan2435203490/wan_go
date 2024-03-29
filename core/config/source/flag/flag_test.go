package flag

import (
	"encoding/json"
	"flag"
	"testing"
)

var (
	dbuser = flag.String("database-db_user", "default", "db db_user")
	dbhost = flag.String("database-host", "", "db host")
	dbpw   = flag.String("database-password", "", "db pw")
)

func initTestFlags() {
	flag.Set("database-host", "localhost")
	flag.Set("database-password", "some-password")
	flag.Parse()
}

func TestFlagsrc_Read(t *testing.T) {
	initTestFlags()
	source := NewSource()
	c, err := source.Read()
	if err != nil {
		t.Error(err)
	}

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Error(err)
	}

	actualDB := actual["database"].(map[string]interface{})
	if actualDB["host"] != *dbhost {
		t.Errorf("expected %v got %v", *dbhost, actualDB["host"])
	}

	if actualDB["password"] != *dbpw {
		t.Errorf("expected %v got %v", *dbpw, actualDB["password"])
	}

	// unset flags should not be loaded
	if actualDB["db_user"] != nil {
		t.Errorf("expected %v got %v", nil, actualDB["db_user"])
	}
}

func TestFlagsrc_ReadAll(t *testing.T) {
	initTestFlags()
	source := NewSource(IncludeUnset(true))
	c, err := source.Read()
	if err != nil {
		t.Error(err)
	}

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Error(err)
	}

	actualDB := actual["database"].(map[string]interface{})

	// unset flag defaults should be loaded
	if actualDB["db_user"] != *dbuser {
		t.Errorf("expected %v got %v", *dbuser, actualDB["db_user"])
	}
}
