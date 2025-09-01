package config

import (
	"testing"
)

func TestInit(t *testing.T) {
	rightConfig := Configuration{
		Auth: Auth{
			Cookie: "test",
		},
		Topic: "test",
	}
	c := Init()
	t.Attr("config", "read to Init")
	if c != nil && *c != rightConfig {
		t.Error("Init fail")
	}
}
