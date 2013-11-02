package apptail

import (
	"encoding/json"
	"path/filepath"
)

func GetDockerAppEnv(rootPath string) (map[string]string, error) {
	data, err := ReadFileLimit(filepath.Join(rootPath, "/app/etc/droplet.env.json"), 50*1000)
	if err != nil {
		return nil, err
	}

	env := map[string]string{}

	err = json.Unmarshal(data, &env)
	return env, err
}
