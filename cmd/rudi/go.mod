module go.xrstf.de/rudi/cmd/rudi

go 1.18

require (
	github.com/TylerBrock/colorjson v0.0.0-20200706003622-8a50f05110d2
	github.com/chzyer/readline v1.5.1
	github.com/spf13/pflag v1.0.5
	go.xrstf.de/go-term-markdown v0.0.0-20231119170546-73a1852b91cc
	go.xrstf.de/rudi v0.1.1-0.20231203234653-fb45c79f7482
	go.xrstf.de/rudi-contrib/semver v0.1.0
	go.xrstf.de/rudi-contrib/yaml v0.1.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/MichaelMure/go-term-text v0.3.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/gomarkdown/markdown v0.0.0-20231115200524-a660076da3fd // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)

replace go.xrstf.de/rudi => ../../

replace github.com/TylerBrock/colorjson => github.com/xrstf/colorjson v0.0.0-20231123184920-5ea6fecf578f
