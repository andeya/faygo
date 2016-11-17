package sqlx

import "testing"

func TestConfig(t *testing.T) {
	t.Logf("%#v", dbConfigs)
	t.Logf("%#v", defaultConfig)
	t.Logf("%#v", dbConfigs[DEFAULTDB_NAME])
}
