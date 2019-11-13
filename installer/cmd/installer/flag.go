/*

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
	"errors"

	"github.com/spf13/pflag"
)

type flags struct {
	configFile        *string
	configFileChanged bool
	dryRun            *bool
	dryRunChanged     bool
}

func parseFlags() *flags {
	flags := &flags{
		configFile: pflag.String("config", "", "Config file containing an AddonInstallerConfiguration"),
		dryRun:     pflag.Bool("dry-run", false, "If true, only print what would happen without actually installing any addons"),
	}

	hideKlogFlags()
	pflag.ErrHelp = errors.New("")

	pflag.Parse()
	flags.configFileChanged = pflag.CommandLine.Changed("config")
	flags.dryRunChanged = pflag.CommandLine.Changed("dry-run")

	return flags
}

func hideKlogFlags() {
	pflag.CommandLine.MarkHidden("alsologtostderr")
	pflag.CommandLine.MarkHidden("log_backtrace_at")
	pflag.CommandLine.MarkHidden("log_dir")
	pflag.CommandLine.MarkHidden("logtostderr")
	pflag.CommandLine.MarkHidden("stderrthreshold")
	pflag.CommandLine.MarkHidden("v")
	pflag.CommandLine.MarkHidden("vmodule")
}
