package conf

import (
	"encoding/json"
	"github.com/opsgenie/oec/git"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	fpath "path/filepath"
	"strings"
	"time"
)

const unknownFileExtErrMessage = "Unknown configuration file extension[%s]. Only \".json\" and \".yml(.yaml)\" types are allowed."

func checkFileExtension(filepath string) error {

	extension := fpath.Ext(strings.ToLower(filepath))

	switch extension {
	case ".json", ".yml", ".yaml":
		return nil
	default:
		return errors.Errorf(unknownFileExtErrMessage, extension)
	}
}

func readFile(filepath string) (*Configuration, error) {

	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	configuration := &Configuration{}
	extension := fpath.Ext(strings.ToLower(filepath))

	switch extension {
	case ".json":
		return configuration, json.Unmarshal(file, configuration)
	case ".yml", ".yaml":
		return configuration, yaml.Unmarshal(file, configuration)
	default:
		return nil, errors.Errorf(unknownFileExtErrMessage, extension)
	}
}

func addHomeDirPrefix(filepath string) string {
	if filepath == "" {
		return filepath
	}

	tildePrefix := "~" + string(os.PathSeparator)

	if strings.HasPrefix(filepath, tildePrefix) {
		homedir, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatalf("Unable to determine user home directory: %s", err)
		}
		return fpath.Join(homedir, strings.TrimPrefix(filepath, tildePrefix))
	}

	return fpath.Clean(filepath)
}

func addHomeDirPrefixToActionMappings(mappings ActionMappings) {
	for index, action := range mappings {
		if action.SourceType == LocalSourceType {
			action.Filepath = addHomeDirPrefix(action.Filepath)
		}
		if action.SourceType == GitSourceType {
			action.GitOptions.PrivateKeyFilepath = addHomeDirPrefix(action.GitOptions.PrivateKeyFilepath)
		}
		action.Stdout = addHomeDirPrefix(action.Stdout)
		action.Stderr = addHomeDirPrefix(action.Stderr)
		mappings[index] = action
	}
}

func AddRepositoryPathToGitActionFilepaths(mappings ActionMappings, repositories git.Repositories) {
	for index, action := range mappings {
		if action.SourceType == GitSourceType {
			repository, err := repositories.Get(action.GitOptions.Url)
			if err != nil {
				continue
			}
			action.Filepath = fpath.Join(repository.Path, action.Filepath)
			mappings[index] = action
		}
	}
}

func PrepareLogFormat() logrus.Formatter {
	formatType := strings.ToLower(os.Getenv("OEC_LOG_FORMAT_TYPE"))
	switch formatType {
	case "text":
		return &logrus.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		}
	case "json":
		return &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	case "colored":
		fallthrough
	default:
		return &logrus.TextFormatter{
			ForceColors:            true,
			FullTimestamp:          true,
			TimestampFormat:        time.RFC3339Nano,
			DisableLevelTruncation: true,
		}
	}
}

func copyActionMappings(mappings ActionMappings) ActionMappings {

	copyActionMappings := make(map[ActionName]MappedAction, len(mappings))
	for k, v := range mappings {
		copyActionMappings[k] = v
	}
	return copyActionMappings
}
