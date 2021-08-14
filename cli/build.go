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
	"fmt"
	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/DataDrake/static-cling/config"
	"github.com/DataDrake/static-cling/render"
	"github.com/DataDrake/static-cling/templates"
	log "github.com/DataDrake/waterlog"
)

func init() {
	cmd.Register(&Build)
}

// Build generates an updated static site tree
var Build = cmd.Sub{
	Name:  "build",
	Alias: "up",
	Short: "Generate the site",
	Flags: &BuildFlags{
		Src:   ".",
		Build: "build",
	},
	Run: BuildRun,
}

// BuildFlags are flags used by the "build" sub-command
type BuildFlags struct {
	Src   string `short:"S" long:"source" desc:"source of project files (default '.')"`
	Build string `short:"B" long:"build"  desc:"name of the buid dir, relative to source (default 'build')"`
}

// BuildRun carries out the "build" sub-command
func BuildRun(r *cmd.Root, s *cmd.Sub) {
	// gFlags := r.Flags.(*GlobalFlags)
	flags := s.Flags.(*BuildFlags)
	log.Infoln("Loading configuration")
	conf, err := config.Load(config.Path(flags.Src))
	if err != nil {
		log.Fatalf("Failed to load configuration: %q\n", err)
	}
	fmt.Printf("%#v\n", conf)
	log.Goodln("Config Loaded.")

	log.Infoln("Loading templates")
	tmpls, err := templates.Load(flags.Src)
	if err != nil {
		log.Fatalf("Failed to load templates: %q\n", err)
	}
	fmt.Printf("%#v\n", tmpls)
	log.Goodln("Templates Loaded.")
	log.Infoln("Setting up the rendering process")
	site, err := render.NewSite(&conf, tmpls)
	if err != nil {
		log.Fatalf("Failed to set up rendering: %q\n", err)
	}
	fmt.Printf("%#v\n", site)
}
