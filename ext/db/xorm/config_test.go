package xorm

import "testing"

func TestConfig(t *testing.T) {
	t.Logf("%#v", dbConfigs)
	t.Logf("%#v", defaultConfig)
}
