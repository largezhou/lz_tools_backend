package migrate_rollback

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/largezhou/lz_tools_backend/app/app_const"
	"github.com/largezhou/lz_tools_backend/app/console/migrate_install"
	"github.com/largezhou/lz_tools_backend/app/model"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"regexp"
)

var db = model.DB.WithContext(context.WithValue(context.Background(), app_const.RequestIdKey, uuid.NewString()))

func New() *cli.Command {
	return &cli.Command{
		Name:      "migrate:rollback",
		Usage:     "回滚迁移",
		UsageText: "migrate:rollback",
		Category:  "migrate",
		Action: func(c *cli.Context) error {
			maxBatch, err := getMaxBatch()
			if err != nil {
				return err
			}

			var ms []migrate_install.Migration
			if result := db.Model(&migrate_install.Migration{}).
				Where("batch = ?", maxBatch).
				Order("file desc").
				Find(&ms); result.Error != nil {
				return nil
			}

			downRegexp, _ := regexp.Compile("(-- __DOWN__[\\s\\S]*)")
			err = db.Transaction(func(tx *gorm.DB) error {
				for _, m := range ms {
					filename := m.File
					filePath := "./app/migration/" + filename
					file, err := os.Open(filePath)
					if err != nil {
						return nil
					}

					fileContent, err := ioutil.ReadAll(file)
					if err != nil {
						return nil
					}

					matches := downRegexp.FindSubmatch(fileContent)
					if len(matches) != 2 {
						return fmt.Errorf("%s 未找到回滚 SQL", filename)
					}

					fmt.Printf("正在回滚 %s\n", filename)

					if result := tx.Exec(string(matches[1])); result.Error != nil {
						return result.Error
					}
				}

				if result := tx.Delete(ms); result.Error != nil {
					return result.Error
				}

				return nil
			})

			return err
		},
	}
}

func getMaxBatch() (uint, error) {
	var maxBatch uint
	if result := db.Model(&migrate_install.Migration{}).
		Select("max(batch)").
		Scan(&maxBatch); result.Error != nil {
		return 0, result.Error
	}
	return maxBatch, nil
}
