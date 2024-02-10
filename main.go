package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
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

	for _, file := range files {
		name := (start + "/" + file.Name())[len(root)+1:]

		log.Info().Str("start", start).Str("name", name).Str("root", root).Msg("processing object")

		if file.IsDir() {
			log.Info().Str("start", start).Str("file.Name", file.Name()).Str("root", root).Msg("processing folder")

			if err := os.MkdirAll(*fOut+"/"+name, 0755); err != nil {
				return nil, err
			}

			// support for nested folders
			f2, err := processFiles(start+"/"+file.Name(), root)
			if err != nil {
				return nil, err
			}

			f = append(f, f2...)
		} else {
			inPath := start + "/" + file.Name()
			scriptPath := start + "/" + name
			htmlPath := *fOut + "/" + name + ".html"

			content, matter, err := readFile(inPath)
			if err != nil {
				return nil, err
			}

			iterator, err := ShellLexer.Tokenise(nil, content)
			if err != nil {
				return nil, err
			}

			buf := bytes.Buffer{}
			if err := HTMLFormatter.Format(&buf, CatppuccinMacchiato, iterator); err != nil {
				log.Fatal().Err(err).Msg("error writing css")
			}

			if err := renderFile(name, buf.String(), htmlPath, matter); err != nil {
				return nil, err
			}
			log.Info().Str("name", name).Str("inPath", inPath).Str("htmlPath", htmlPath).Str("scriptPath", scriptPath).Str("root", root).Str("start", start).Str("file.Name", file.Name()).Str("file", file.Name()).Str("start", start).Str("root", root).Str("htmlPath", htmlPath).Str("scriptPath", scriptPath).Str("inPath", inPath).Str("name", name).Str("file.Name", file.Name()).Str("start", start).Str("root", root).Msg("rendered")
			//if err := copyFile(inPath, scriptPath); err != nil {
			//	return nil, err
			//}

			f = append(f, internal.Script{
				Name:   name,
				Matter: matter,
			})
		}
	}

	return f, nil
}

func writeStringToFile(string, dst string) error {
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = d.WriteString(string)
	return err
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
