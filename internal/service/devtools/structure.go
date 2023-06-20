package devtools

import (
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/sourcebuild"
)

func expandSourceBuildEnvVarsParams(eVars []interface{}) ([]*sourcebuild.ProjectEnvEnvVars, error) {
	envVars := make([]*sourcebuild.ProjectEnvEnvVars, 0, len(eVars))

	for _, v := range eVars {
		env := new(sourcebuild.ProjectEnvEnvVars)
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "key":
				env.Key = ncloud.String(value.(string))
			case "value":
				env.Value = ncloud.String(value.(string))
			}
		}
		envVars = append(envVars, env)
	}

	return envVars, nil
}
