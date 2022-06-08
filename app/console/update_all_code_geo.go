package console

import (
	"github.com/largezhou/lz_tools_backend/app/service"
	"github.com/urfave/cli/v2"
)

var codeService = service.NewCodeService()

func NewUpdateAllCodeGeoCommand() *cli.Command {
	return &cli.Command{
		Name:  "updateAllCodeGeo",
		Usage: "更新所有场所码的坐标到 redis",
		Action: func(c *cli.Context) error {
			return codeService.UpdateAllCodeGeo(ctx)
		},
	}
}
