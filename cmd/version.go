// Copyright © 2022 ezbuy & LITB Team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CommitHash string

func version(commit string) string {
	if commit == "" {
		return fmt.Sprintf("ezorm v%d.%d.%d", vMajor, vMinor, vPatch)
	}
	return fmt.Sprintf("ezorm v%d.%d.%d-%s", vMajor, vMinor, vPatch, commit)
}

const (
	vMajor = 2
	vMinor = 6
	vPatch = 4
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "EzOrm 版本信息",
	Long:  `EzOrm 版本信息`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Fprintln(os.Stdout, version(CommitHash))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
