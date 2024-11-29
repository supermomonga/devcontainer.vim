//go:build windows

package tools

const PortForwarderFileName = "port-forwarder.exe"

// devcontainer-cli のダウンロード URL
const downloadURLPortForwarderCliPattern = "https://github.com/mikoto2000/port-forwarder/releases/download/{{ .TagName }}/port-forwarder-linux-amd64"

