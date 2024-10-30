package config


import (
	"os"
	"path/filepath"
	"strings"
	"github.com/spf13/viper"
)

type Config struct {
	viperObj   *viper.Viper
	FileName   string
	ConfigType string
}

func GetConfig(conFileName string) (conf Config, err error) {
	// var conf Config
	var viperTMP *viper.Viper
	conf.FileName = conFileName
	// filePath, fileName := filepath.Split(conFileName) //分割路径中的目录与文件
	fileExt := filepath.Ext(conFileName) //返回路径中的扩展名 如果没有点，返回空
	// filenameOnly := strings.TrimSuffix(fileName, fileExt)
	// if filePath == "" {
	// 	filePath = "."
	// }

	cTYPE := strings.ToUpper(fileExt[1:])
	switch cTYPE {
	case "YAML":
		conf.ConfigType = "yaml"
	case "TOML":
		conf.ConfigType = "toml"
	case "INI":
		conf.ConfigType = "ini"
	default:
		conf.ConfigType = "json"
	}
	if conf.ConfigType == "ini" {
		viperTMP = viper.NewWithOptions(viper.KeyDelimiter("="))
	} else {
		viperTMP = viper.New()
	}
	viperTMP.SetConfigType(conf.ConfigType) // REQUIRED if the config file does not have the extension in the name
	viperTMP.SetConfigFile(conFileName)
	// viperTMP.SetConfigName(filenameOnly)    // name of config file (without extension)
	// viperTMP.AddConfigPath(filePath)        // path to look for the config file in
	// if filePath != "." {
	// 	viperTMP.AddConfigPath(".") // optionally look for config in the working directory
	// }

	if _, err1 := os.Stat(conFileName); err1 == nil {
		err = viperTMP.ReadInConfig()
	}

	conf.viperObj = viperTMP
	return

}
func (config *Config) SetValue(key string, value any) {
	config.viperObj.Set(key, value)
	// viper.Set("host.port", 5899)
	// viper.SetDefault("ContentDir", "content")
	// viper.SetDefault("LayoutDir", "layouts")
	// viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
}

func (config *Config) GetValue(key string) any {
	return config.viperObj.Get(key)
	// GetFloat64(key string) : float64
	// GetIntSlice(key string) : []int
	// GetStringMap(key string) : map[string]any
	// GetStringMapString(key string) : map[string]string
	// GetStringSlice(key string) : []string
	// GetTime(key string) : time.Time
	// GetDuration(key string) : time.Duration
	// IsSet(key string) : bool
	// AllSettings() : map[string]any
}

func (config *Config) GetValueString(key string) string {
	return config.viperObj.GetString(key)
}
func (config *Config) GetValueBool(key string) bool {
	return config.viperObj.GetBool(key)
}
func (config *Config) GetValueInt(key string) int {
	return config.viperObj.GetInt(key)
}

func (config *Config) GetSub(key string) *viper.Viper {
	return config.viperObj.Sub(key)
}

func (config *Config) SaveConfig() (err error) {
	err = config.viperObj.WriteConfig()
	return
}

func (config *Config) Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) (err error) {
	// type config struct {
	// 	Chart struct{
	// 		Values map[string]any
	// 	}
	// }
	// var C config
	// v.Unmarshal(&C)
	err = config.viperObj.Unmarshal(rawVal, opts...)
	return
}
