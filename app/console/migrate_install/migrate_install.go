package migrate_install

import (
	"fmt"
	"github.com/largezhou/gin_starter/app/model"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"io/fs"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
)

var db = model.DB

func New() *cli.Command {
	return &cli.Command{
		Name:      "migrate:install",
		Usage:     "使用迁移文件",
		UsageText: "migrate:install",
		Category:  "migrate",
		Action: func(c *cli.Context) error {
			if err := createMigrationTableIfNeeded(); err != nil {
				return err
			}

			filenames, err := getMigrationFilenames()
			if err != nil {
				return err
			}

			migrationMap, maxBatch := getExistsMigrationInfo()
			var newFilenames []string
			for _, filename := range filenames {
				if _, ok := migrationMap[filename]; !ok {
					newFilenames = append(newFilenames, filename)
				}
			}

			if len(newFilenames) == 0 {
				fmt.Println("没有新的迁移")
				return nil
			}

			var upRegexp, _ = regexp.Compile("(-- __UP__[\\s\\S]*?)-- __DOWN__")
			err = db.Transaction(func(tx *gorm.DB) error {
				var newMigrations []*migration
				for _, filename := range newFilenames {
					newMigrations = append(newMigrations, &migration{
						File:  filename,
						Batch: maxBatch + 1,
					})

					file, err := os.Open("./app/migration/" + filename)
					if err != nil {
						return nil
					}

					fileContent, err := ioutil.ReadAll(file)
					if err != nil {
						return err
					}

					matches := upRegexp.FindSubmatch(fileContent)
					if len(matches) != 2 {
						return fmt.Errorf("%s 未找到升级 SQL", filename)
					}

					fmt.Printf("正在迁移 %s\n", filename)
					fmt.Println(string(matches[1]))
					if result := tx.Exec(string(matches[1])); result.Error != nil {
						return result.Error
					}
				}

				if result := tx.Create(newMigrations); result.Error != nil {
					return result.Error
				}

				return nil
			})

			if err != nil {
				return err
			}

			return nil
		},
	}
}

type migration struct {
	model.Model
	File  string `gorm:"type:varchar(200);unique;not null"`
	Batch uint   `gorm:"type:integer;not null"`
}

func createMigrationTableIfNeeded() error {
	if db.Migrator().HasTable(&migration{}) {
		return nil
	}

	if err := db.Migrator().CreateTable(&migration{}); err != nil {
		return err
	}

	return nil
}

func getMigrationFilenames() ([]string, error) {
	var files []fs.FileInfo
	var filenames []string
	files, err := ioutil.ReadDir("./app/migration/")
	if err != nil {
		return filenames, err
	}

	for _, fileInfo := range files {
		name := strings.ToLower(fileInfo.Name())
		if strings.HasSuffix(name, ".sql") {
			filenames = append(filenames, name)
		}
	}

	sort.Strings(filenames)

	return filenames, nil
}

func getExistsMigrationInfo() (map[string]migration, uint) {
	msMap := make(map[string]migration)
	var maxBatch uint = 0

	var ms []migration
	db.Model(&migration{}).Select("file, batch").Find(&ms)
	for _, m := range ms {
		if m.Batch > maxBatch {
			maxBatch = m.Batch
		}
		msMap[m.File] = m
	}
	return msMap, maxBatch
}
