package gen

import (
	"github.com/iWinston/gf-cli/library/mlog"
)

// doGenService implements the "gen service" command.
func doGenService() {
	name, description, systemName := getArgs()

	genFileForce(serviceTemplateMap["default"], "./app/system/admin/service/internal", name+"_service.go", name, description, systemName)
	genFile(serviceTemplateMap["index"], "./app/system/admin/service", name+"_service.go", name, description, systemName)

	mlog.Print("gen service done!")
}
