package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/adrg/frontmatter"
	"github.com/sascha-andres/reuse"
	"github.com/sascha-andres/reuse/flag"
	"gopkg.in/yaml.v2"
)

var (
	directory string
	banner    string
)

func init() {
	log.SetPrefix("[frontmatter] ")
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.LUTC)

	flag.StringVar(&directory, "d", "", "Directory to search for files")
	flag.StringVar(&banner, "b", "", "Banner to add to files")
}

func main() {
	flag.Parse()

	if len(directory) == 0 {
		log.Fatal("Directory not set")
	}

	if len(banner) == 0 {
		log.Fatal("Banner not set")
	}

	err := addBanner(directory, banner)
	if err != nil {
		log.Fatal(err)
	}
}

func addBanner(directory, banner string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	for i := range files {
		if files[i].IsDir() {
			continue
		}
		err := addBannerToFile(directory, files[i].Name(), banner)
		if err != nil {
			return err
		}
	}
	return nil
}

func addBannerToFile(directory, filename, banner string) error {
	mdFile := path.Join(directory, filename)

	read, err := os.OpenFile(mdFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer reuse.DiscardError(read.Close)
	var deserializedFrontmatter map[string]any
	content, err := frontmatter.Parse(read, &deserializedFrontmatter)
	if err != nil {
		return nil
	}
	err = read.Close()
	if err != nil {
		return err
	}
	if _, ok := deserializedFrontmatter["banner"]; ok {
		return nil
	}

	log.Printf("Adding banner to %s", filename)
	deserializedFrontmatter["banner"] = banner

	b, err := yaml.Marshal(deserializedFrontmatter)
	if err != nil {
		return err
	}

	return os.WriteFile(mdFile, []byte(fmt.Sprintf("---\n%s---\n%s", string(b), string(content))), 0644)
}
