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
	"html/template"
	"path/filepath"
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
		SrcDir:    ".",
		OutputDir: "build",
	},
	Run: BuildRun,
}

// BuildFlags are flags used by the "build" sub-command
type BuildFlags struct {
	SrcDir    string `short:"S" long:"srcdir" desc:"source of project files (default '.')"`
	OutputDir string `short:"O" long:"outputdir" desc:"name of the output dir, relative to source (default 'build')"`
}

// BuildRun carries out the "build" sub-command
func BuildRun(r *cmd.Root, s *cmd.Sub) {
	// gFlags := r.Flags.(*GlobalFlags)
	flags := s.Flags.(*BuildFlags)
	buildDir := filepath.Join(flags.SrcDir, flags.OutputDir)
	util.CreateDir(buildDir)
	assetsDir := filepath.Join(flags.SrcDir, "assets")
	util.CopyDir(assetsDir, buildDir)
	templateDir := filepath.Join(flags.SrcDir, "templates", "*")
	templates := template.Must(template.ParseGlob(templateDir))
	println(templates.DefinedTemplates())
}
