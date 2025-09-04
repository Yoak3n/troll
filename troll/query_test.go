package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Yoak3n/troll/scanner/model/dto"
)

func TestGenWordCountMap(t *testing.T) {
	videos := make([]dto.VideoDataOutput, 0)
	dir, _ := os.ReadDir("../data/cache/华为")
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		fp, err := os.ReadFile("../data/cache/华为/" + f.Name())
		if err != nil {
			t.Error(err)
		}
		data := &dto.VideoDataOutput{}
		err = json.Unmarshal(fp, &data)
		if err != nil {
			t.Error(err)
		}
		videos = append(videos, *data)
	}

	meta, w := genWordCountMap(videos)
	n := getTopN(w, 10)
	for _, v := range n {
		fmt.Printf("%.f %s\t", v.Value, v.Key)
		fmt.Printf("\ttf:%.2f idf:%.2f count:%d\n", meta[v.Key].tf, meta[v.Key].idf, meta[v.Key].count)
	}

}
