/*
Copyright © 2019 ZeroLabs <zeroc0d3.team@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"

	"github.com/dimiro1/banner"
	_ "github.com/dimiro1/banner/autoload"
	"github.com/mattn/go-colorable"
	"github.com/zeroc0d3/multivpn/cmd"
)

func initLogo() {
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString(" OpenVPN Selector {{ .AnsiColor.Green }}(Running){{ .AnsiColor.Default }} ...\n\n"))
}

func main() {
	initLogo()
	cmd.Execute()
}
