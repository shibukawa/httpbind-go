package parser

import "sort"

func sortStrings(ss []string) {
	sort.Strings(ss)
}

func sortRoutes(routes []Route) {
	sort.SliceStable(routes, func(i, j int) bool {
		a, b := routes[i], routes[j]
		if a.Method != b.Method {
			return a.Method < b.Method
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Handler.Name != b.Handler.Name {
			return a.Handler.Name < b.Handler.Name
		}
		return a.Handler.Form < b.Handler.Form
	})
}
