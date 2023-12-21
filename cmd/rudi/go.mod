module go.xrstf.de/rudi/cmd/rudi

go 1.18

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/TylerBrock/colorjson v0.0.0-20200706003622-8a50f05110d2
	github.com/chzyer/readline v1.5.1
	github.com/muesli/termenv v0.15.2
	github.com/spf13/pflag v1.0.5
	go.xrstf.de/rudi v0.4.0
	go.xrstf.de/rudi-contrib/semver v0.1.4
	go.xrstf.de/rudi-contrib/uuid v0.1.4
	go.xrstf.de/rudi-contrib/yaml v0.1.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)

replace go.xrstf.de/rudi => ../../

replace github.com/TylerBrock/colorjson => github.com/xrstf/colorjson v0.0.0-20231123184920-5ea6fecf578f
