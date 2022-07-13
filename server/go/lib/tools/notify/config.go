package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cli.NewApp()
	app := &cli.App{
		Name:  "notify",
		Usage: "通知",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ding_token",
				Usage:   "dingding webhook url",
				EnvVars: []string{"PLUGIN_DING_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "ding_secret",
				Usage:   "dingding secret",
				EnvVars: []string{"PLUGIN_DING_SECRET"},
			},
			&cli.StringFlag{
				Name:    "drone_commit_sha",
				Usage:   "drone_commit_sha",
				EnvVars: []string{"DRONE_COMMIT_SHA"},
			},
			&cli.StringFlag{
				Name:    "commit_message",
				Usage:   "git commit message",
				EnvVars: []string{"DRONE_COMMIT_MESSAGE"},
			},
			&cli.StringFlag{
				Name:    "current_branch",
				Usage:   "current git branch",
				EnvVars: []string{"PLUGIN_CURRENT_BRANCH"},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	//读取当前环境类型
	DroneServerDomain string
	//读取当前提交的commit版本
	CommitSHA     string
	CommitMessage string
	DingToken     string
	DingSecret    string
	CurrentBranch string
}

func GetConfig(c *cli.Context) *Config {
	return &Config{
		DroneServerDomain: c.String("drone_server_domain"),
		CommitSHA:         c.String("drone_commit_sha"),
		CommitMessage:     c.String("commit_message"),
		DingToken:         c.String("ding_token"),
		DingSecret:        c.String("ding_secret"),
		CurrentBranch:     c.String("current_branch"),
	}
}
