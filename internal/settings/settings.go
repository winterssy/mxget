package settings

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/winterssy/gjson"
	"github.com/winterssy/glog"
	"github.com/winterssy/mxget/pkg/provider"
)

const (
	dir      = "./downloads"
	platform = "nc"
)

var (
	Cfg   = &Config{}
	Limit int
	Tag   bool
	Lyric bool
	Force bool
)

type (
	Config struct {
		Dir      string `json:"dir"`
		Platform string `json:"platform"`

		// 预留字段，其它设置项
		others   map[string]interface{} `json:"-"`
		filePath string                 `json:"-"`
	}
)

func Init() {
	err := Cfg.setup()
	if err != nil {
		_ = Cfg.Save()
		glog.Fatalf("Initialize config failed, reset to defaults: %v", err)
	}
}

func (c *Config) setup() error {
	if c.setupConfigFile() != nil {
		return c.initConfigFile()
	}

	err := c.loadConfigFile()
	if err != nil {
		return err
	}

	err = c.check()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) setupConfigFile() error {
	var cfgDir string
	xdgDir := os.Getenv("XDG_CONFIG_HOME")
	if xdgDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			cfgDir = "."
		} else {
			cfgDir = filepath.Join(home, ".config", "mxget")
		}
	} else {
		cfgDir = filepath.Join(xdgDir, "mxget")
	}

	if cfgDir == "." || os.MkdirAll(cfgDir, 0755) != nil {
		c.filePath = ".mxget.json"
	} else {
		c.filePath = filepath.Join(cfgDir, "mxget.json")
	}

	_, err := os.Stat(c.filePath)
	return err
}

func (c *Config) initConfigFile() error {
	c.Dir = dir
	c.Platform = platform
	return c.Save()
}

func (c *Config) loadConfigFile() error {
	b, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return err
	}

	return gjson.Unmarshal(b, c)
}

func (c *Config) check() error {
	if provider.GetDesc(c.Platform) == "unknown" {
		rawPlatform := c.Platform
		c.Platform = platform
		return fmt.Errorf("unexpected music platform: %q", rawPlatform)
	}

	err := os.MkdirAll(c.Dir, 0755)
	if err != nil {
		c.Dir = dir
		return fmt.Errorf("cant't make download dir: %s", err.Error())
	}

	return nil
}

func (c *Config) Save() error {
	b, err := gjson.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.filePath, b, 0644)
}

func (c *Config) Reset() {
	_ = c.initConfigFile()
}
