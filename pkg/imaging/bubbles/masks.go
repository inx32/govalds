package bubbles

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"path"
	"slices"

	"forge.pi.home.arpa/govalds/bot/pkg/imaging"
	"forge.pi.home.arpa/govalds/bot/pkg/imaging/bubble"

	"github.com/disgoorg/disgo/discord"
)

type Properties struct {
	Alpha       *string `json:"alpha"`
	Overlay     *string `json:"overlay"`
	Description string  `json:"description"`
}

type Bubble struct {
	Mask       bubble.Mask
	Properties *Properties
}

//go:embed all:assets
var DefaultAssets embed.FS

var Bubbles map[string]Bubble

func loadImage(fsys fs.FS, p string) (image.Image, error) {
	st, err := fs.Stat(fsys, p)
	if err != nil {
		return nil, fmt.Errorf("stat image file %s: %w", p, err)
	}

	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("stat image file %s: expected regular file", p)
	}

	fd, err := fsys.Open(p)
	if err != nil {
		return nil, fmt.Errorf("open image file %s: %w", p, err)
	}
	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, fmt.Errorf("decode image file %s: %w", p, err)
	}
	return img, nil
}

func loadMask(properties *Properties, fsys fs.FS, directory string) (bubble.Mask, error) {
	var overlay *image.NRGBA
	var alpha *image.Alpha

	if properties.Overlay != nil {
		p := path.Join(directory, *properties.Overlay)
		img, err := loadImage(fsys, p)
		if err != nil {
			return nil, fmt.Errorf("load overlay: %w", err)
		}
		overlay = imaging.ToNRGBA(img)
	}

	if properties.Alpha != nil {
		p := path.Join(directory, *properties.Alpha)
		img, err := loadImage(fsys, p)
		if err != nil {
			return nil, fmt.Errorf("load alpha: %w", err)
		}
		alpha = imaging.ToAlpha(img)
	}

	return bubble.NewMask(overlay, alpha)
}

func loadProperties(fsys fs.FS, directory string) (*Properties, error) {
	path := path.Join(directory, "properties.json")
	fd, err := fsys.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open properties file %s: %w", path, err)
	}
	defer fd.Close()

	decoder := json.NewDecoder(fd)
	decoder.DisallowUnknownFields()

	properties := &Properties{}
	err = decoder.Decode(properties)
	if err != nil {
		return nil, fmt.Errorf("decode properties file %s: %w", path, err)
	}

	return properties, nil
}

func Load(fsys fs.FS, root string) (map[string]Bubble, error) {
	bubbles := make(map[string]Bubble)
	entries, err := fs.ReadDir(fsys, root)
	if err != nil {
		return nil, fmt.Errorf("read directory %s: %w", root, err)
	}

	for _, entry := range entries {
		p := path.Join(root, entry.Name())
		if !entry.IsDir() {
			return nil, fmt.Errorf("stat %s: expected directory", p)
		}

		properties, err := loadProperties(fsys, p)
		if err != nil {
			return nil, fmt.Errorf("load properties for %s: %w", p, err)
		}

		mask, err := loadMask(properties, fsys, p)
		if err != nil {
			return nil, fmt.Errorf("load mask for %s: %w", p, err)
		}

		bubbles[entry.Name()] = Bubble{
			Mask:       mask,
			Properties: properties,
		}
	}

	if _, ok := bubbles["default"]; !ok {
		return nil, errors.New("no bubble with name \"default\"")
	}
	return bubbles, nil
}

func init() {
	bubbles, err := Load(DefaultAssets, "assets")
	if err != nil {
		panic(fmt.Errorf("unable to load bubbles: %w", err))
	}
	Bubbles = bubbles
}

func GetChoices() []discord.ApplicationCommandOptionChoiceString {
	var choices []discord.ApplicationCommandOptionChoiceString

	var keys []string
	for k := range Bubbles {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		bubble := Bubbles[k]
		choices = append(choices, discord.ApplicationCommandOptionChoiceString{
			Name:  bubble.Properties.Description,
			Value: k,
		})
	}

	return choices
}
