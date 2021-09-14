module github.com/zyedidia/flare

go 1.16

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/jwalton/gchalk v1.0.3
	github.com/zyedidia/ftdetect v0.0.0-20210226205021-01c766da7946
	github.com/zyedidia/gpeg v0.0.0-20210712001934-4e907283728d
	github.com/zyedidia/rope v0.0.0-20210616205215-37fbf22eab3a
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/zyedidia/gpeg => ../gpeg
