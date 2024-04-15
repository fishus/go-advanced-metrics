package app

import (
	"os"
	"text/template"
)

func PrintBuildInfo(ver, date, commit string) {
	type buildInfo struct {
		Version string
		Date    string
		Commit  string
	}

	ver = defaultValue(ver, "N/A")
	date = defaultValue(date, "N/A")
	commit = defaultValue(commit, "N/A")

	info := buildInfo{
		Version: ver,
		Date:    date,
		Commit:  commit,
	}

	const tpl = `Build version: {{.Version}}
Build date: {{.Date}}
Build commit: {{.Commit}}
`

	t := template.Must(template.New("list").Parse(tpl))

	err := t.Execute(os.Stdout, info)
	if err != nil {
		panic(err)
	}
}

func defaultValue(v, defaultValue string) string {
	if v == "" {
		return defaultValue
	}
	return v
}
