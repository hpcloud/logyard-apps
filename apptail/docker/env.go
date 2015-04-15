package docker

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetDockerAppEnv(rootPath string) (map[string]string, error) {
	data, err := readFileLimit(filepath.Join(rootPath, "/home/stackato/etc/droplet.env.json"), 50*1000)
	switch err := err.(type) {
	default:
		return nil, err
	case *os.PathError:
		return map[string]string{}, nil
	case nil:
		// pass
	}

	env := map[string]string{}

	err = json.Unmarshal(data, &env)
	return env, err
}
