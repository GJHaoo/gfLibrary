package gen

import (
	"github.com/iWinston/gf-cli/library/mlog"
)

// doGenApi implements the "gen api" command.
func doGenApi() {
	name, description, systemName := getArgs()

	// genFileForce(apiTemplateMap["default"], "./gen/api", name+"_api.go", name, description)
	// genFileForce(apiTemplateMap["core"], "./gen/api", name+".go", name, description)
	genFile(apiTemplateMap["index"], "./app/system/"+systemName+"/api", name+"_api.go", name, description, systemName)

	mlog.Print("gen api done!")
}
