package console

import (
	"github.com/largezhou/lz_tools_backend/app/console/make_migration"
	"github.com/largezhou/lz_tools_backend/app/console/migrate_install"
	"github.com/largezhou/lz_tools_backend/app/console/migrate_rollback"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	make_migration.New(),
	migrate_install.New(),
	migrate_rollback.New(),
}
