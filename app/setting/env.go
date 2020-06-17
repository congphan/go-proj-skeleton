package setting

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"go-prj-skeleton/app/strutil"
)

type envSettings struct {
	Env string `envconfig:"env" default:"development"`

	PrintEnvs string `envconfig:"gohelpers_print_envs" default:""`

	// PostgreSql
	PostgreHost           string `envconfig:"postgre_host" default:"db"`
	PostgrePort           string `envconfig:"postgre_port" default:"5432"`
	PostgreUser           string `envconfig:"postgre_user" default:"admin"`
	PostgrePassword       string `envconfig:"postgre_password" default:"moneyforward@123"`
	PostgreDatabaseName   string `envconfig:"postgre_database_name" default:"postgres"`
	PostgreMaxConnections int    `envconfig:"postgre_max_connections" default:"16"`
}

// ProjectEnvSettings is the singeton hold all the env vars
var ProjectEnvSettings *envSettings

//environment variables must be in following format
//SETTING_POSTGRE_HOST
//SETTING_POSTGRE_USER
func (settings *envSettings) readEnvironmentVariables() error {

	err := envconfig.Process("setting", settings)
	if err != nil {
		return err
	}
	return nil
}

// EnvSettingsInit reads settings from env vars
func EnvSettingsInit(filterKeys []string) {
	ProjectEnvSettings = &envSettings{}

	if err := ProjectEnvSettings.readEnvironmentVariables(); err != nil {
		log.WithError(err).Errorln("error while read env vars for project settings")
	}

	if ProjectEnvSettings.PrintEnvs != "" {
		filterKeys = strutil.CleanEmpty(strings.Split(ProjectEnvSettings.PrintEnvs, ","))
	}

	log.Println("setting env vars:")
	PrintEnvs("setting", ProjectEnvSettings, filterKeys)
}

// PrintEnvs ...
func PrintEnvs(prefix string, spec interface{}, filterKeys []string) {
	noFilter := strutil.Include(filterKeys, "all")

	s := reflect.ValueOf(spec).Elem()
	envInfos := []string{}

	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		alt := typeOfSpec.Field(i).Tag.Get("envconfig")
		fieldName := typeOfSpec.Field(i).Name
		if alt != "" {
			fieldName = alt
		}

		key := strings.ToUpper(fmt.Sprintf("%s_%s", prefix, fieldName))
		value := f.Interface()
		if noFilter {
			if strings.Contains(key, "SECRET") || strings.Contains(key, "PASSWORD") { // not to print secrets if not asked to
				continue
			}

			envInfos = append(envInfos, fmt.Sprintf("%s: `%v`", key, value))
			continue
		}

		if strutil.Include(filterKeys, key) {
			envInfos = append(envInfos, fmt.Sprintf("%s: `%v`", key, value))
			continue
		}
	}

	fmt.Println(strings.Join(envInfos, "\n"))
}
