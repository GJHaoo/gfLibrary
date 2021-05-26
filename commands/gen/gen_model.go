package gen

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/iWinston/gf-cli/library/mlog"
)

// doGenModel implements the "gen model" command.
func doGenModel() {
	parser, err := gcmd.Parse(g.MapStrBool{
		"n, name":       true,
		"d,description": true,
	})
	if err != nil {
		mlog.Fatal(err)
	}
	name := parser.GetOpt("name")
	description := parser.GetOpt("description")

	genFile(modelTemplate, "./app/model", name+"_model.go", name, description)

	mlog.Print("gen model done!")
}
