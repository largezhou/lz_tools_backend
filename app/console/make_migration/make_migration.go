package make_migration

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"time"
)

var filepath string
var err error

func writeError() error {
	os.Remove(filepath)
	return err
}

func New() *cli.Command {
	return &cli.Command{
		Name:      "make:migration",
		Usage:     "创建一个数据库迁移文件",
		UsageText: "make:migration 文件名",
		Category:  "make",
		Action: func(c *cli.Context) error {
			filename := c.Args().Get(0)
			if filename == "" {
				return errors.New("缺少文件名参数")
			}
			filename = time.Now().Format("2006_01_02_15_04_05") + "_" + filename + ".sql"
			filepath = "./app/migration/" + filename

			var file *os.File
			file, err = os.Create(filepath)
			if err != nil {
				return err
			}
			defer file.Close()

			var stubFile *os.File
			stubFile, err = os.Open("./app/console/make_migration/stub/make_migration.stub")
			if err != nil {
				return writeError()
			}
			stub, err := ioutil.ReadAll(stubFile)
			if err != nil {
				return writeError()
			}

			if _, err = file.Write(stub); err != nil {
				return writeError()
			}

			fmt.Println("迁移文件创建成功: " + filename)

			return nil
		},
	}
}
