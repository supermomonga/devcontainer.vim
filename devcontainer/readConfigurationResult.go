package devcontainer

import (
	"encoding/json"
)

type ReadConfigurationError struct {
	msg string
}

func (e *ReadConfigurationError) Error() string {
	return e.msg
}

// `devcontainers read-configuration` コマンドの実行結果スキーマ
//
//	Example:	{
//							"configuration":{
//								"name":"development environment",
//								"dockerComposeFile":[
//									"../docker-compose.yaml"
//								],
//								"service":"app",
//								"workspaceFolder":"/work",
//								"remoteUser":"root",
//								"configFilePath": {
//									"$mid":1,
//									"fsPath":"/home/mikoto/project/oasiz-mqtt-client/.devcontainer/devcontainer.json",
//									"path":"/home/mikoto/project/oasiz-mqtt-client/.devcontainer/devcontainer.json",
//									"scheme":"vscode-fileHost"
//								}
//							},
//							"workspace": {
//								"workspaceFolder":"/work"
//							}
//						}
type ReadConfigurationCommandResult struct {
	Configuration Configuration `json:"configuration"`
}

type Configuration struct {
	ConfigFilePath ConfigFilePath `json:"configFilePath"`
}

type ConfigFilePath struct {
	FsPath string `json:"fsPath"`
}

func GetConfigFilePath(readConfigurationCommandResult string) (string, error) {
	result, err := UnmarshalReadConfigurationCommandResult([]byte(readConfigurationCommandResult))
	if err != nil {
		return "", &ReadConfigurationError{msg: "`devcontainer read-configuration` の出力パースに失敗しました。`.devcontainer.json が存在することと、 docker エンジンが起動していることを確認してください。"}
	}

	return result.Configuration.ConfigFilePath.FsPath, nil
}

func UnmarshalReadConfigurationCommandResult(data []byte) (ReadConfigurationCommandResult, error) {
	var result ReadConfigurationCommandResult

	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
