package model

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	res := &CommentResponse{}
	fp, err := os.Open("..\\..\\test\\2025-08-31T234200.200.json")
	if err != nil {
		t.Error(err)
	}
	defer fp.Close()
	buf, _ := io.ReadAll(fp)

	err = json.Unmarshal(buf, &res)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", res)
}
