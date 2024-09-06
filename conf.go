// @since 2024-09-06 13:31:22
// @author acrazing <joking.young@gmail.com>

package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

type node struct {
	path     string
	branches []string
	value    any
}

var configMap = map[any]*node{}

// Register is used for register a struct to a path in a config file.
// It should be called in init function.
func Register[T any](path string) {
	branches := strings.Split(path, ".")
	var key *T = nil
	var value T
	if n, ok := configMap[key]; ok {
		panic(fmt.Sprintf("type %T is already in use with path %s", value, n.path))
	}
	if n, ok := configMap[path]; ok {
		panic(fmt.Sprintf("path %s is already in use type %T", n.path, n.value))
	}
	n := &node{
		path:     path,
		branches: branches,
		value:    value,
	}
	configMap[key] = n
	configMap[path] = n
}

func mapJson[T any](input any) (*T, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	data := new(T)
	err = json.Unmarshal(raw, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Registry is the root config registry
type Registry struct {
	node *simplejson.Json
}

// Get a configuration node from registry
func Get[T any](registry *Registry) (*T, error) {
	var key *T
	var value T
	n := configMap[key]
	if n == nil {
		return nil, errors.New(fmt.Sprintf("invalid config type %T", value))
	}
	node := registry.node.GetPath(n.branches...)
	if node.Interface() == nil {
		return nil, errors.New(fmt.Sprintf("mis config path %s for type %T", n.path, value))
	}
	return mapJson[T](node)
}

// Load config from a file
func Load(filename string) (*Registry, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var data any
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".json":
		err = json.Unmarshal(raw, &data)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(raw, &data)
	case ".toml":
		err = toml.Unmarshal(raw, &data)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported config file extension: %s", ext))
	}
	if err != nil {
		return nil, err
	}
	return New(data), nil
}

// New create a registry with specified config data
func New(data any) *Registry {
	node := simplejson.New()
	node.SetPath(nil, data)
	return &Registry{node: node}
}

// Default returns registered config nodes default value structure
func Default() any {
	node := simplejson.New()
	for k, v := range configMap {
		if _, ok := k.(string); ok {
			continue
		}
		node.SetPath(v.branches, &v.value)
	}
	return node.Interface()
}
