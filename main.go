package main

import (
	"context"
	_ "embed"
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	sitemap "github.com/carlosstrand/go-sitemap"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sh.devminer.xyz/internal"
	"sh.devminer.xyz/internal/web"
)

//go:embed static/styles.css
var CSS string

var (
	CatppuccinMacchiato = styles.Get("catppuccin-macchiato")
	ShellLexer          = lexers.Get("shell")
	HTMLFormatter       = html.New()
)

var (
	fBaseUrl = flag.String("base-url", "https://sh.devminer.xyz", "Base URL for the website")
	fIn      = flag.String("in", "./scripts", "Folder to read the input files")
	fOut     = flag.String("out", "./out", "Folder to put the output files")
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Parse()

	if err := os.MkdirAll(*fOut, 0755); err != nil {
		log.Fatal().Err(err).Msg("error creating output folder")
	}

	files, err := processFiles(*fIn, *fIn)
	if err != nil {
		log.Fatal().Err(err).Msg("error processing files")
	}

	f, err := os.Create(*fOut + "/index.html")
	if err != nil {
		log.Fatal().Err(err).Msg("error creating index.html")
	}
	defer f.Close()

	if err := web.Index(files).Render(context.Background(), f); err != nil {
		log.Fatal().Err(err).Msg("error rendering index.html")
	}

	if err := writeStringToFile(CSS, *fOut+"/styles.css"); err != nil {
		log.Fatal().Err(err).Msg("error writing styles.css")
	}

	if err := generateSitemap(files, *fOut); err != nil {
		log.Fatal().Err(err).Msg("error generating sitemap")
	}

	log.Info().Msg("exiting")
}

func readFile(path string) (string, internal.Matter, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", internal.Matter{}, err
	}
	defer f.Close()

	matter := internal.Matter{}
	raw, err := frontmatter.Parse(f, &matter)
	if err != nil {
		return "", internal.Matter{}, err
	}

	log.Info().Str("path", path).Str("matter.Description", matter.Description).Strs("matter.Tags", matter.Tags).Msg("read file")

	return string(raw), matter, nil
}

func renderFile(name, content, outPath string, matter internal.Matter) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return web.Script(*fBaseUrl, name, content, matter).Render(context.Background(), f)
}

func processFiles(start, root string) ([]internal.Script, error) {
	log.Info().Str("start", start).Str("root", root).Msg("processing")

	files, err := os.ReadDir(start)
	if err != nil {
		return nil, err
	}

	f := make([]internal.Script, 0, len(files))
	folders := make([]fs.DirEntry, 0)

	for _, file := range files {
		name := (start + "/" + file.Name())[len(root)+1:]

		log.Info().Str("start", start).Str("name", name).Str("root", root).Msg("processing object")

		if file.IsDir() {
			folders = append(folders, file)
		} else {
			inPath := start + "/" + file.Name()
			scriptPath := *fOut + "/" + name
			htmlPath := scriptPath + ".html"

			content, matter, err := readFile(inPath)
			if err != nil {
				return nil, err
			}

			highlit, err := internal.Highlight("shell", content)
			if err != nil {
				return nil, err
			}

			if err := renderFile(name, highlit, htmlPath, matter); err != nil {
				return nil, err
			}

			if err := writeStringToFile(content, scriptPath); err != nil {
				return nil, err
			}

			log.Info().Str("name", name).Str("inPath", inPath).Str("htmlPath", htmlPath).Str("scriptPath", scriptPath).Str("root", root).Str("start", start).Str("file.Name", file.Name()).Str("file", file.Name()).Str("start", start).Str("root", root).Str("htmlPath", htmlPath).Str("scriptPath", scriptPath).Str("inPath", inPath).Str("name", name).Str("file.Name", file.Name()).Str("start", start).Str("root", root).Msg("rendered")

			f = append(f, internal.Script{
				Name:   name,
				Matter: matter,
			})
		}
	}

	for _, folder := range folders {
		name := (start + "/" + folder.Name())[len(root)+1:]

		log.Info().Str("file.Name", folder.Name()).Msg("processing folder")

		if err := os.MkdirAll(*fOut+"/"+name, 0755); err != nil {
			return nil, err
		}

		// support for nested folders
		f2, err := processFiles(start+"/"+folder.Name(), root)
		if err != nil {
			return nil, err
		}

		f = append(f, f2...)
	}

	return f, nil
}

func writeStringToFile(content, dst string) error {
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = d.WriteString(content)
	return err
}

func generateSitemap(files []internal.Script, root string) error {
	items := make([]*sitemap.SitemapItem, 0)
	for _, file := range files {
		items = append(items, &sitemap.SitemapItem{
			Loc:        *fBaseUrl + "/" + file.Name + ".html",
			LastMod:    time.Now(),
			ChangeFreq: "monthly",
			Priority:   0.5,
		})
	}

	smStr, err := sitemap.NewSitemap(items, nil).ToXMLString()
	if err != nil {
		return err
	}

	smStr = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + smStr

	return writeStringToFile(smStr, filepath.Join(root, "sitemap.xml"))
}
