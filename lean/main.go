package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
)

const banner = `
 _                        ______ _                 _
| |                      / _____) |               | |
| |      ____ ____ ____ | /     | | ___  _   _  _ | |
| |     / _  ) _  |  _ \| |     | |/ _ \| | | |/ || |
| |____( (/ ( ( | | | | | \_____| | |_| | |_| ( (_| |
|_______)____)_||_|_| |_|\______)_|\___/ \____|\____|

`

var (
	version          = "0.0.1"
	commandBuildFrom = "go tool"
	isDeployFromGit  = false
)

func thirdPartyCommand(c *cli.Context, _cmdName string) {
	cmdName := "lean-" + _cmdName

	// executeble not found:
	execPath, err := exec.LookPath(cmdName)
	if e, ok := err.(*exec.Error); ok {
		if e.Err == exec.ErrNotFound {
			cli.ShowAppHelp(c)
			return
		}
	}
	cmd := exec.Command(execPath, c.Args()[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// disable the log prefix
	log.SetFlags(0)

	go func() {
		_ = checkUpdate()
	}()

	// add banner text to help text
	cli.AppHelpTemplate = banner + cli.AppHelpTemplate
	cli.SubcommandHelpTemplate = banner + cli.SubcommandHelpTemplate

	app := cli.NewApp()
	app.Name = "lean"
	app.Version = version
	app.Usage = "Command line to manage and deploy LeanCloud apps"

	app.CommandNotFound = thirdPartyCommand

	app.Commands = []cli.Command{
		{
			Name:   "login",
			Usage:  "登录 LeanCloud 账户。",
			Action: loginAction,
		},
		{
			Name:   "info",
			Usage:  "查看当前登录用户以及应用信息。",
			Action: infoAction,
		},
		{
			Name:   "up",
			Usage:  "本地启动云引擎应用。",
			Action: upAction,
		},
		{
			Name:   "init",
			Usage:  "初始化云引擎项目。",
			Action: initAction,
		},
		{
			Name:   "switch",
			Usage:  "切换当前项目关联的 LeanCloud 应用。",
			Action: switchAction,
		},
		{
			Name:   "deploy",
			Usage:  "部署云引擎项目到服务器",
			Action: deployAction,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "g",
					Usage:       "从 git 部署项目",
					Destination: &isDeployFromGit,
				},
			},
		},
	}

	app.Run(os.Args)
}
