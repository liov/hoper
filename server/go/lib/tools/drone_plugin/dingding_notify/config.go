package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	fmt.Println("开始了")
	app := &cli.App{
		Name:  "notify",
		Usage: "通知",
		Action: func(c *cli.Context) error {
			config := GetConfig(c)
			return Notify(config)
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
				Name:    "repo",
				Usage:   "git repo",
				EnvVars: []string{"DRONE_REPO"},
			},
			&cli.StringFlag{
				Name:    "commit",
				Usage:   "git commit",
				EnvVars: []string{"DRONE_COMMIT"},
			},
			&cli.StringFlag{
				Name:    "commit_author_name",
				Usage:   "git commit author name",
				EnvVars: []string{"DRONE_COMMIT_AUTHOR_NAME"},
			},
			&cli.StringFlag{
				Name:    "commit_author",
				Usage:   "git commit author",
				EnvVars: []string{"DRONE_COMMIT_AUTHOR"},
			},
			&cli.StringFlag{
				Name:    "commit_link",
				Usage:   "git commit link",
				EnvVars: []string{"DRONE_COMMIT_LINK"},
			},
			&cli.StringFlag{
				Name:    "commit_ref",
				Usage:   "git commit ref",
				EnvVars: []string{"DRONE_COMMIT_REF"},
			},
			&cli.StringFlag{
				Name:    "commit_message",
				Usage:   "git commit message",
				EnvVars: []string{"DRONE_COMMIT_MESSAGE"},
			},
			&cli.StringFlag{
				Name:    "commit_branch",
				Usage:   "git commit branch",
				EnvVars: []string{"DRONE_COMMIT_BRANCH"},
			},
			&cli.StringFlag{
				Name:    "commit_tag",
				Usage:   "git commit tag",
				EnvVars: []string{"DRONE_TAG"},
			},
			&cli.StringFlag{
				Name:    "drone_build_link",
				Usage:   "drone build link",
				EnvVars: []string{"DRONE_BUILD_LINK"},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Repo          string
	CommitAuthor  string
	Commit        string
	CommitTag     string
	CommitRef     string
	CommitMessage string
	CommitBranch  string
	DingToken     string
	DingSecret    string
	BuildLink     string
}

func GetConfig(c *cli.Context) *Config {
	return &Config{
		Repo:          c.String("repo"),
		CommitAuthor:  c.String("commit_author_name"),
		Commit:        c.String("commit"),
		CommitTag:     c.String("commit_tag"),
		CommitRef:     c.String("commit_ref"),
		CommitMessage: c.String("commit_message"),
		CommitBranch:  c.String("commit_branch"),
		DingToken:     c.String("ding_token"),
		DingSecret:    c.String("ding_secret"),
		BuildLink:     c.String("drone_build_link"),
	}
}
