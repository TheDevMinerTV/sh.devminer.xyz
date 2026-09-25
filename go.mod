module sh.devminer.xyz

go 1.26.0

toolchain go1.27.1

require (
	github.com/a-h/templ v0.3.1020
	github.com/adrg/frontmatter v0.2.0
	github.com/rs/zerolog v1.35.1
)

require github.com/dlclark/regexp2/v2 v2.8.0 // indirect

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/alecthomas/chroma/v2 v2.27.0
	github.com/carlosstrand/go-sitemap v0.0.0-20191230193616-37cd6896357b
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	golang.org/x/sys v0.48.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/a-h/templ v0.2.552 => github.com/thedevminertv/templ v0.0.0-20240131105554-64c9eb5c4c3d
