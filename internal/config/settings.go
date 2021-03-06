package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/photoprism/photoprism/pkg/fs"
	"github.com/photoprism/photoprism/pkg/txt"
	"gopkg.in/yaml.v2"
)

// DisableSettings returns true if the user is not allowed to change settings.
func (c *Config) DisableSettings() bool {
	return c.params.DisableSettings
}

type MapsSettings struct {
	Animate int    `json:"animate" yaml:"animate"`
	Style   string `json:"style" yaml:"style"`
}

type LibrarySettings struct {
	ConvertRaw     bool `json:"raw" yaml:"raw"`
	CreateThumbs   bool `json:"thumbs" yaml:"thumbs"`
	CompleteRescan bool `json:"rescan" yaml:"rescan"`
	MoveImported   bool `json:"move" yaml:"move"`
	HidePrivate    bool `json:"private" yaml:"private"`
	RequireReview  bool `json:"review" yaml:"review"`
	GroupRelated   bool `json:"group" yaml:"group"`
}

type FeatureSettings struct {
	Upload   bool `json:"upload" yaml:"upload"`
	Import   bool `json:"import" yaml:"import"`
	Labels   bool `json:"labels" yaml:"labels"`
	Places   bool `json:"places" yaml:"places"`
	Archive  bool `json:"archive" yaml:"archive"`
	Download bool `json:"download" yaml:"download"`
	Edit     bool `json:"edit" yaml:"edit"`
	Share    bool `json:"share" yaml:"share"`
}

// Settings contains Web UI settings
type Settings struct {
	Theme    string          `json:"theme" yaml:"theme"`
	Language string          `json:"language" yaml:"language"`
	Maps     MapsSettings    `json:"maps" yaml:"maps"`
	Features FeatureSettings `json:"features" yaml:"features"`
	Library  LibrarySettings `json:"library" yaml:"library"`
}

// NewSettings returns a empty Settings
func NewSettings() *Settings {
	return &Settings{
		Theme:    "default",
		Language: "en",
		Maps: MapsSettings{
			Animate: 0,
			Style:   "streets",
		},
		Features: FeatureSettings{
			Upload:   true,
			Import:   true,
			Labels:   true,
			Places:   true,
			Archive:  true,
			Download: true,
			Edit:     true,
			Share:    true,
		},
		Library: LibrarySettings{
			CompleteRescan: false,
			ConvertRaw:     false,
			CreateThumbs:   false,
			MoveImported:   false,
			GroupRelated:   true,
			RequireReview:  true,
			HidePrivate:    true,
		},
	}
}

// Propagate updates settings in other packages as needed.
func (s *Settings) Propagate() {

}

// Load uses a yaml config file to initiate the configuration entity.
func (s *Settings) Load(fileName string) error {
	if !fs.FileExists(fileName) {
		return fmt.Errorf("settings file not found: %s", txt.Quote(fileName))
	}

	yamlConfig, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlConfig, s); err != nil {
		return err
	}

	s.Propagate()

	return nil
}

// Save uses a yaml config file to initiate the configuration entity.
func (s *Settings) Save(fileName string) error {
	data, err := yaml.Marshal(s)

	if err != nil {
		return err
	}

	s.Propagate()

	if err := ioutil.WriteFile(fileName, data, os.ModePerm); err != nil {
		return err
	}

	s.Propagate()

	return nil
}

// initSettings initializes user settings from a config file.
func (c *Config) initSettings() {
	c.settings = NewSettings()
	p := c.SettingsFile()

	if err := c.settings.Load(p); err != nil {
		log.Error(err)
	}

	c.settings.Propagate()
}

// Settings returns the current user settings.
func (c *Config) Settings() *Settings {
	return c.settings
}
