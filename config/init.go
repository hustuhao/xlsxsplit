package config

import (
	"strings"

	"github.com/metafates/xlsxsplit/app"
	"github.com/metafates/xlsxsplit/filesystem"
	"github.com/metafates/xlsxsplit/where"
	"github.com/spf13/viper"
)

// Format is the format of the config file
// Available options are: json, yaml, toml
const Format = "toml"

var EnvKeyReplacer = strings.NewReplacer(".", "_")

func Init() error {
	viper.SetConfigName(app.Name)
	viper.SetConfigType(Format)
	viper.SetFs(filesystem.Api())
	viper.AddConfigPath(where.Config())
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvPrefix(app.Name)
	viper.SetEnvKeyReplacer(EnvKeyReplacer)

	setDefaults()

	err := viper.ReadInConfig()

	switch err.(type) {
	case viper.ConfigFileNotFoundError:
		// Use defaults then
		return nil
	default:
		return err
	}
}
