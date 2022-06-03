package console

import (
	"github.com/largezhou/gin_starter/app/console/make_migration"
	"github.com/largezhou/gin_starter/app/console/migrate_install"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	make_migration.New(),
	migrate_install.New(),
}
