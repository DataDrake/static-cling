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

package content

import (
	"sort"
)

// Pages is a list of Pages
type Pages []*Page

// Latest gets a new list of these pages is order from newest to oldest
func (ps Pages) Latest() Pages {
	ls := latestPages(ps)
	sort.Sort(ls)
	return Pages(ls)
}

// Reverse returns the Pages in reverse order (alphanumeric by default)
func (ps Pages) Reverse() (out Pages) {
	for _, page := range ps {
		out = append(Pages{page}, out...)
	}
	return
}

// latestPages is a Pages sorted from Newest to Oldest
type latestPages Pages

// Len returns the length of the list (satisfies sort.Sort)
func (ps latestPages) Len() int {
	return len(ps)
}

// Less is true if this Page is newer (satisfies sort.Sort)
func (ps latestPages) Less(i, j int) bool {
	return ps[j].Date.After(ps[i].Date)
}

// Swap the entries of the list (satisfies sort.Sort)
func (ps latestPages) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
