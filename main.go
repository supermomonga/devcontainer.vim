// TODO: vimFilePath の名前変更(hostVimFilePath?)

package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/mikoto2000/devcontainer.vim/devcontainer"
	"github.com/mikoto2000/devcontainer.vim/docker"
	"github.com/mikoto2000/devcontainer.vim/tools"
	"github.com/mikoto2000/devcontainer.vim/util"
)

const FLAG_NAME_LICENSE = "license"
const FLAG_NAME_HELP_LONG = "help"
const FLAG_NAME_HELP_SHORT = "h"
const FLAG_NAME_VERSION_LONG = "version"
const SPLIT_ARG_MARK = "--"

//go:embed LICENSE
var license string

//go:embed NOTICE
var notice string

const APP_NAME = "devcontainer.vim"

func main() {
	// コマンドラインオプションのパース

	// devcontainer.vim 用のディレクトリ作成
	// 1. ユーザーコンフィグ用ディレクトリ
	//    `os.UserConfigDir` + `devcontainer.vim`
	// 2. ユーザーキャッシュ用ディレクトリ
	//    `os.UserCacheDir` + `devcontainer.vim`
	util.CreateDirectory(os.UserConfigDir, APP_NAME)
	appCacheDir, binDir, appConfigDir := util.CreateDirectory(os.UserCacheDir, APP_NAME)

	devcontainerVimArgProcess := (&cli.App{
		Name:                   "devcontainer.vim",
		Usage:                  "devcontainer for vim.",
		Version:                "0.0.1",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:               FLAG_NAME_LICENSE,
				Aliases:            []string{"l"},
				Value:              false,
				DisableDefaultText: true,
				Usage:              "show licensesa.",
			},
		},
		Action: func(cCtx *cli.Context) error {
			// ライセンスフラグが立っていればライセンスを表示して終
			if cCtx.Bool(FLAG_NAME_LICENSE) {
				fmt.Println(license)
				fmt.Println()
				fmt.Println(notice)
				os.Exit(0)
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:            "run",
				Usage:           "Run container use `docker run`",
				UsageText:       "devcontainer.vim run [DOCKER_OPTIONS...] [DOCKER_ARGS...]",
				HideHelp:        true,
				SkipFlagParsing: true,
				Action: func(cCtx *cli.Context) error {
					// `docker run` でコンテナを立てる

					// 必要なファイルのダウンロード
					vimPath, err := tools.VIM.Install(binDir)
					if err != nil {
						panic(err)
					}

					// Requirements のチェック
					// 1. docker
					isExistsDocker := util.IsExistsCommand("docker")
					if !isExistsDocker {
						fmt.Fprintf(os.Stderr, "docker not found.")
						os.Exit(1)
					}

					// コンテナ起動
					docker.Run(cCtx.Args().Slice(), vimPath)

					return nil
				},
			},
			{
				Name:            "start",
				Usage:           "Run `devcontainer up` and `devcontainer exec`",
				UsageText:       "devcontainer.vim start [DEVCONTAINER_OPTIONS...] WORKSPACE_FOLDER",
				HideHelp:        true,
				SkipFlagParsing: true,
				Action: func(cCtx *cli.Context) error {
					// devcontainer でコンテナを立てる

					// 必要なファイルのダウンロード
					vimPath, err := tools.VIM.Install(binDir)
					if err != nil {
						panic(err)
					}

					devcontainerFilePath, err := tools.DEVCONTAINER.Install(binDir)
					if err != nil {
						panic(err)
					}

					// コマンドライン引数の末尾は `--workspace-folder` の値として使う
					args := cCtx.Args().Slice()
					workspaceFolder := args[len(args)-1]
					configFilePath, err := createConfigFile(devcontainerFilePath, workspaceFolder, appConfigDir)
					if err != nil {
						panic(err)
					}

					// devcontainer を用いたコンテナ立ち上げ
					devcontainer.ExecuteDevcontainer(args, devcontainerFilePath, vimPath, configFilePath)

					return nil
				},
			},
			{
				Name:            "down",
				Usage:           "Stop and remove devcontainers.",
				UsageText:       "devcontainer.vim down WORKSPACE_FOLDER",
				HideHelp:        true,
				SkipFlagParsing: true,
				Action: func(cCtx *cli.Context) error {
					// devcontainer でコンテナを立てる

					// 必要なファイルのダウンロード
					devcontainerPath, err := tools.DEVCONTAINER.Install(binDir)
					if err != nil {
						panic(err)
					}

					// devcontainer を用いたコンテナ終了
					devcontainer.Down(cCtx.Args().Slice(), devcontainerPath)

					// 設定ファイルを削除
					// コマンドライン引数の末尾は `--workspace-folder` の値として使う
					args := cCtx.Args().Slice()
					workspaceFolder := args[len(args)-1]
					configDir := util.GetConfigDir(appCacheDir, workspaceFolder)

					fmt.Printf("Remove configuration file: `%s`\n", configDir)
					os.RemoveAll(configDir)

					return nil
				},
			},
		},
	})

	// アプリ実行
	err := devcontainerVimArgProcess.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

// devcontainer.vim 起動時に使用する設定ファイルを作成する
// 設定ファイルは、 devcontainer.vim のキャッシュ内の `config` ディレクトリに、
// ワークスペースフォルダのパスを md5 ハッシュ化した名前のディレクトリに格納する.
func createConfigFile(devcontainerFilePath string, workspaceFolder string, appConfigDir string) (string, error) {
	// devcontainer の設定ファイルパス取得
	configFilePath, err := devcontainer.GetConfigurationFilePath(devcontainerFilePath, workspaceFolder)
	if err != nil {
		return "", err
	}

	// devcontainer.vim 用の追加設定ファイルを探す
	configurationFileName := configFilePath[:len(configFilePath)-len(filepath.Ext(configFilePath))]
	additionalConfigurationFilePath := configurationFileName + ".vim.json"

	// 設定管理フォルダに JSON を配置
	mergedConfigFilePath, err := util.CreateConfigFileForDevcontainerVim(appConfigDir, workspaceFolder, configFilePath, additionalConfigurationFilePath)

	fmt.Printf("Use configuration file: `%s`", mergedConfigFilePath)

	return mergedConfigFilePath, err
}
