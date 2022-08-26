package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

func (configor *Configor) getENVPrefix(config interface{}) string {
	if configor.Config.ENVPrefix == "" {
		if prefix := os.Getenv("CONFIGOR_ENV_PREFIX"); prefix != "" {
			return prefix
		}
		return "Configor"
	}
	return configor.Config.ENVPrefix
}

func getConfigurationFileWithENVPrefix(file, env string) (string, time.Time, error) {
	var (
		envFile string
		extname = path.Ext(file)
	)

	if extname == "" {
		envFile = fmt.Sprintf("%v.%v", file, env)
	} else {
		envFile = fmt.Sprintf("%v.%v%v", strings.TrimSuffix(file, extname), env, extname)
	}

	if fileInfo, err := os.Stat(envFile); err == nil && fileInfo.Mode().IsRegular() {
		return envFile, fileInfo.ModTime(), nil
	}
	return "", time.Now(), fmt.Errorf("failed to find file %v", file)
}

func (configor *Configor) getConfigurationFiles(watchMode bool, files ...string) ([]string, map[string]time.Time) {
	var resultKeys []string
	var results = map[string]time.Time{}

	if !watchMode && (configor.Config.Debug || configor.Config.Verbose) {
		fmt.Printf("Current environment: '%v'\n", configor.GetEnvironment())
	}

	for i := len(files) - 1; i >= 0; i-- {
		foundFile := false
		file := files[i]

		// check configuration
		if fileInfo, err := os.Stat(file); err == nil && fileInfo.Mode().IsRegular() {
			foundFile = true
			resultKeys = append(resultKeys, file)
			results[file] = fileInfo.ModTime()
		}

		// check configuration with env
		if file, modTime, err := getConfigurationFileWithENVPrefix(file, configor.GetEnvironment()); err == nil {
			foundFile = true
			resultKeys = append(resultKeys, file)
			results[file] = modTime
		}

		// check example configuration
		if !foundFile {
			if example, modTime, err := getConfigurationFileWithENVPrefix(file, "example"); err == nil {
				if !watchMode && !configor.Silent {
					fmt.Printf("Failed to find configuration %v, using example file %v\n", file, example)
				}
				resultKeys = append(resultKeys, example)
				results[example] = modTime
			} else if !configor.Silent {
				fmt.Printf("Failed to find configuration %v\n", file)
			}
		}
	}
	return resultKeys, results
}

func processFile(config interface{}, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	switch {
	case strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml"):
		return yaml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".toml"):
		return toml.Unmarshal(data, config)
	case strings.HasSuffix(file, ".json"):
		return unmarshalJSON(data, config)
	default:
		if err := toml.Unmarshal(data, config); err == nil {
			return nil
		}

		if err := unmarshalJSON(data, config); err == nil {
			return nil
		} else if strings.Contains(err.Error(), "json: unknown field") {
			return err
		}

		yamlError := yaml.Unmarshal(data, config)

		if yamlError == nil {
			return nil
		} else if yErr, ok := yamlError.(*yaml.TypeError); ok {
			return yErr
		}

		return errors.New("failed to decode config")
	}
}

// unmarshalJSON unmarshals the given data into the config interface.
// If the errorOnUnmatchedKeys boolean is true, an error will be returned if there
// are keys in the data that do not match fields in the config interface.
func unmarshalJSON(data []byte, config interface{}) error {
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(config)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func getPrefixForStruct(prefixes []string, fieldStruct *reflect.StructField) []string {
	if fieldStruct.Anonymous && fieldStruct.Tag.Get("anonymous") == "true" {
		return prefixes
	}
	return append(prefixes, fieldStruct.Name)
}

func (configor *Configor) processDefaults(config interface{}) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank {
			// Set default configuration if blank
			if value := fieldStruct.Tag.Get("default"); value != "" {
				if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
					return err
				}
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		switch field.Kind() {
		case reflect.Struct:
			if err := configor.processDefaults(field.Addr().Interface()); err != nil {
				return err
			}
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
					if err := configor.processDefaults(field.Index(i).Addr().Interface()); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (configor *Configor) processTags(config interface{}, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			envNames    []string
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
			envName     = fieldStruct.Tag.Get("env") // read configuration from shell env
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		if envName == "" {
			envNames = append(envNames, strings.Join(append(prefixes, fieldStruct.Name), "_"))                  // Configor_DB_Name
			envNames = append(envNames, strings.ToUpper(strings.Join(append(prefixes, fieldStruct.Name), "_"))) // CONFIGOR_DB_NAME
		} else {
			envNames = []string{envName}
		}

		if configor.Config.Verbose {
			fmt.Printf("Trying to load struct `%v`'s field `%v` from env %v\n", configType.Name(), fieldStruct.Name, strings.Join(envNames, ", "))
		}

		// Load From Shell ENV
		for _, env := range envNames {
			if value := os.Getenv(env); value != "" {
				if configor.Config.Debug || configor.Config.Verbose {
					fmt.Printf("Loading configuration for struct `%v`'s field `%v` from env %v...\n", configType.Name(), fieldStruct.Name, env)
				}

				switch reflect.Indirect(field).Kind() {
				case reflect.Bool:
					switch strings.ToLower(value) {
					case "", "0", "f", "false":
						field.Set(reflect.ValueOf(false))
					default:
						field.Set(reflect.ValueOf(true))
					}
				case reflect.String:
					field.Set(reflect.ValueOf(value))
				default:
					if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
						return err
					}
				}
				break
			}
		}

		if isBlank := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()); isBlank && fieldStruct.Tag.Get("required") == "true" {
			// return error if it is required but blank
			return errors.New(fieldStruct.Name + " is required, but blank")
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if err := configor.processTags(field.Addr().Interface(), getPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Slice {
			if arrLen := field.Len(); arrLen > 0 {
				for i := 0; i < arrLen; i++ {
					if reflect.Indirect(field.Index(i)).Kind() == reflect.Struct {
						if err := configor.processTags(field.Index(i).Addr().Interface(), append(getPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(i))...); err != nil {
							return err
						}
					}
				}
			} else {
				// load slice from env
				newVal := reflect.New(field.Type().Elem()).Elem()
				if newVal.Kind() == reflect.Struct {
					idx := 0
					for {
						newVal = reflect.New(field.Type().Elem()).Elem()
						if err := configor.processTags(newVal.Addr().Interface(), append(getPrefixForStruct(prefixes, &fieldStruct), fmt.Sprint(idx))...); err != nil {
							return err
						} else if reflect.DeepEqual(newVal.Interface(), reflect.New(field.Type().Elem()).Elem().Interface()) {
							break
						} else {
							idx++
							field.Set(reflect.Append(field, newVal))
						}
					}
				}
			}
		}
	}
	return nil
}

func (configor *Configor) load(config interface{}, watchMode bool, files ...string) (err error, changed bool) {
	defer func() {
		if configor.Config.Debug || configor.Config.Verbose {
			if err != nil {
				fmt.Printf("Failed to load configuration from %v, got %v\n", files, err)
			}

			fmt.Printf("Configuration:\n  %#v\n", config)
		}
	}()

	configFiles, configModTimeMap := configor.getConfigurationFiles(watchMode, files...)

	if watchMode {
		if len(configModTimeMap) == len(configor.configModTimes) {
			var changed bool
			for f, t := range configModTimeMap {
				if v, ok := configor.configModTimes[f]; !ok || t.After(v) {
					changed = true
				}
			}

			if !changed {
				return nil, false
			}
		}
	}

	// process defaults
	configor.processDefaults(config)

	for _, file := range configFiles {
		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Loading configurations from file '%v'...\n", file)
		}
		if err = processFile(config, file); err != nil {
			return err, true
		}
	}
	configor.configModTimes = configModTimeMap

	if prefix := configor.getENVPrefix(config); prefix == "-" {
		err = configor.processTags(config)
	} else {
		err = configor.processTags(config, prefix)
	}

	return err, true
}

func (configor *Configor) handle(handle func([]byte), watchMode bool, files ...string) (err error, changed bool) {
	defer func() {
		if configor.Config.Debug || configor.Config.Verbose {
			if err != nil {
				fmt.Printf("Failed to load configuration from %v, got %v\n", files, err)
			}
		}
	}()

	configFiles, configModTimeMap := configor.getConfigurationFiles(watchMode, files...)

	if watchMode {
		if len(configModTimeMap) == len(configor.configModTimes) {
			var changed bool
			for f, t := range configModTimeMap {
				if v, ok := configor.configModTimes[f]; !ok || t.After(v) {
					changed = true
				}
			}

			if !changed {
				return nil, false
			}
		}
	}

	for _, file := range configFiles {
		if configor.Config.Debug || configor.Config.Verbose {
			fmt.Printf("Loading configurations from file '%v'...\n", file)
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err, true
		}
		handle(data)
	}
	configor.configModTimes = configModTimeMap

	return err, true
}
