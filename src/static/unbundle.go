package static

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Shopify/themekit/src/cmdutil"
	"github.com/Shopify/themekit/src/colors"
)

var zipData string

// Register will set the zip data that can be decompressed
func Register(data string) {
	zipData = data
}

// Unbundle will saftely put all files in place without overwriting files that
// already exist
func Unbundle(ctx *cmdutil.Ctx) error {
	zipReader, err := zip.NewReader(strings.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}

	paths := []string{}
	files := make(map[string]*zip.File, len(zipReader.File))
	for _, zipFile := range zipReader.File {
		files[zipFile.Name] = zipFile
		paths = append(paths, zipFile.Name)
	}

	sort.Slice(paths, func(i, j int) bool { return strings.Compare(paths[i], paths[j]) < 0 })

	dirs := map[string]bool{}
	for _, path := range paths {
		if err := prepareDir(ctx, path, dirs); err != nil {
			return err
		}
		contents, err := files[path].Open()
		if err != nil {
			return err
		}
		if err := writeFile(ctx, path, contents); err != nil {
			return err
		}
		contents.Close()
	}
	return nil
}

func prepareDir(ctx *cmdutil.Ctx, path string, checked map[string]bool) error {
	dir := filepath.Dir(path)
	if _, ok := checked[dir]; ok {
		return nil
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		ctx.Log.Printf("%s directory %s.\n", colors.Green("Created"), dir)
	} else {
		ctx.Log.Printf("%s directory %s.\n", colors.Yellow("Found"), dir)
	}
	checked[dir] = true
	return nil
}

func writeFile(ctx *cmdutil.Ctx, path string, contents io.Reader) error {
	if _, err := os.Stat(path); err == nil {
		ctx.Log.Printf("\t%s file %s, already exists.\n", colors.Yellow("Skipped"), path)
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	defer file.Sync()

	_, err = io.Copy(file, contents)
	if err == nil {
		ctx.Log.Printf("\t%s file %s.\n", colors.Green("Created"), file.Name())
	}
	return nil
}
