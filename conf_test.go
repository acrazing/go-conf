// @since 2024-09-06 14:23:45
// @author acrazing <joking.young@gmail.com>

package conf_test

import (
	"encoding/json"
	"github.com/acrazing/go-conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DB struct {
	Url string `json:"url" toml:"url" yaml:"url"`
}

type NestedStore struct {
	Path string `json:"path" yaml:"path"`
}

func init() {
	conf.Register[DB]("db")
	conf.Register[NestedStore]("nested.store")
}

const url = "https://www.google.com/"

func testLoad(t *testing.T, ext string, isErr bool) *conf.Registry {
	r, err := conf.Load("data/config." + ext)
	if isErr {
		assert.Error(t, err)
		return nil
	}
	assert.NoError(t, err)
	db, err := conf.Get[DB](r)
	assert.NoError(t, err)
	assert.Equal(t, url, db.Url)
	_, err = conf.Get[conf.Registry](r)
	assert.Error(t, err)
	return r
}

func testNested(t *testing.T, ext string, isErr bool) *conf.Registry {
	r := testLoad(t, ext, false)
	store, err := conf.Get[NestedStore](r)
	if isErr {
		assert.Error(t, err)
		return r
	}
	assert.Equal(t, "/", store.Path)
	return r
}

func TestAll(t *testing.T) {
	testNested(t, "json", false)
	testNested(t, "toml", true)
	testLoad(t, "yaml", false)
	testLoad(t, "yml", false)
	testLoad(t, "ini", true)
	testLoad(t, "json5", true)
	testLoad(t, "bad", true)
}

func TestDefault(t *testing.T) {
	d := conf.Default()
	b, err := json.MarshalIndent(d, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, `{
  "db": {
    "url": ""
  },
  "nested": {
    "store": {
      "path": ""
    }
  }
}`, string(b))
}

func ExampleGet() {
	type ServerConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	// Register first
	conf.Register[ServerConfig]("server")

	// Init registry
	r, err := conf.Load("config.json")
	if err != nil {
		panic(err.Error())
	}

	// Get config from registry
	serverConfig, err := conf.Get[ServerConfig](r)
	if err != nil {
		panic(err)
	}
	println(serverConfig.Host, serverConfig.Port)
}

func ExampleDefault() {
	type ServerConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	// Register first
	conf.Register[ServerConfig]("server")

	d := conf.Default()
	txt, _ := json.MarshalIndent(d, "", "  ")
	println(txt)
}
