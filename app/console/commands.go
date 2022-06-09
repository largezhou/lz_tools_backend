package console

import (
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	NewMakeMigrationCommand(),
	NewMigrateInstallCommand(),
	NewMigrateRollbackCommand(),
	NewUpdateAllCodeGeoCommand(),
}
