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
	"github.com/DataDrake/static-cling/util"
	log "github.com/DataDrake/waterlog"
	"os"
	"path/filepath"
)

func init() {
	cmd.Register(&Init)
}

// Init populates a new blank project tree
var Init = cmd.Sub{
	Name:  "init",
	Alias: "new",
	Short: "Create a new static-cling project",
	Flags: &InitFlags{
		DestDir: ".",
	},
	Run: InitRun,
}

// InitFlags are flags used by the "init" sub-command
type InitFlags struct {
	DestDir string `short:"D" long:"destdir" desc:"destination for new project (default is '.')"`
}

// InitRun carries out the "init" sub-command
func InitRun(r *cmd.Root, s *cmd.Sub) {
	// gFlags := r.Flags.(*GlobalFlags)
	flags := s.Flags.(*InitFlags)
	// args := s.Args.(*InitArgs)
	log.Infof("Creating project directory '%s'\n", flags.DestDir)
	if err := os.MkdirAll(flags.DestDir, 0755); err != nil {
		log.Fatalf("Failed to create project directory '%s', reason: %s\n", flags.DestDir, err)
	}
	log.Goodln("Done.")
	util.CreateDir(filepath.Join(flags.DestDir, "assets"))
	util.CreateDir(filepath.Join(flags.DestDir, "config"))
	util.CreateDir(filepath.Join(flags.DestDir, "content"))
	util.CreateDir(filepath.Join(flags.DestDir, "content", "blog"))
	util.CreateDir(filepath.Join(flags.DestDir, "templates"))
}
