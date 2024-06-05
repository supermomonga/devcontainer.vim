package tools

import (
	"strings"
	"text/template"

	"github.com/mikoto2000/devcontainer.vim/util"
)

// devcontainer/cli のツール情報
var DEVCONTAINER Tool = Tool{
	FileName: DEVCONTAINER_FILE_NAME,
	CalculateDownloadUrl: func() string {
		latestTagName, err := util.GetLatestReleaseFromGitHub("mikoto2000", "devcontainers-cli")
		if err != nil {
			panic(err)
		}

		pattern := "pattern"
		tmpl, err := template.New(pattern).Parse(VIM_DOWNLOAD_URL_PATTERN)
		if err != nil {
			panic(err)
		}

		tmplParams := map[string]string{"TagName": latestTagName}
		var downloadUrl strings.Builder
		err = tmpl.Execute(&downloadUrl, tmplParams)
		if err != nil {
			panic(err)
		}
		return downloadUrl.String()
	},
	installFunc: func(downloadUrl string, installDir string, fileName string, override bool) (string, error) {
		return simpleInstall(downloadUrl, installDir, fileName, override)
	},
}