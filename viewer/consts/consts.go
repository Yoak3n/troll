package consts

import (
	"os"
	"path"

	"github.com/Yoak3n/gulu/util"
)

var TrollPath = "troll"

func init() {
	userConfigPath, _ := os.UserConfigDir()
	TrollPath = path.Join(userConfigPath, "troll")
	_ = util.CreateDirNotExists(TrollPath)
}
