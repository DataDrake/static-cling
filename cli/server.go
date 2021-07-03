//
// Copyright 2021 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cli

import (
	"github.com/DataDrake/cli-ng/v2/cmd"
	"net/http"
)

func init() {
	cmd.Register(&Server)
}

// Server runs a local web server to preview changes made to the website
var Server = cmd.Sub{
	Name:  "server",
	Alias: "run",
	Short: "Start and HTTP server to serve out the files in the output directory",
	Run:   ServerRun,
}

// ServerRun carries out the "server" sub-command
func ServerRun(r *cmd.Root, s *cmd.Sub) {
	http.ListenAndServe(":8080", http.FileServer(http.Dir("build")))
}
