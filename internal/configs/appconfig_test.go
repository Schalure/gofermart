package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createAppConfigFile(t *testing.T) {

	teseCases := []struct {
		name string
		appConfig AppConfig
		want struct {
			data string
		}
	}{
		{
			name: "simple test",
			appConfig: AppConfig{
				Env: Debug,
			},
			want: struct{data string}{
				data: "Environment: debug\n",
			},
		},
	}


	for _, test := range teseCases {
		t.Run(test.name,func(t *testing.T) {

			err := test.appConfig.createAppConfigFile(defaultAppConfigFilePath)
			require.NoErrorf(t, err, "stop test with error: ", err)

			data, err := os.ReadFile(defaultAppConfigFilePath)
			require.NoErrorf(t, err, "stop test with error: ", err)

			assert.Equal(t, test.want.data, string(data))

			err = os.Remove(defaultAppConfigFilePath)
			require.NoErrorf(t, err, "stop test with error: ", err)
		})
	}
}

func Test_getAppConfigFromFile(t *testing.T) {

	teseCases := []struct {
		name string
		data string
		want struct {
			appConfig AppConfig
		}
	}{
		{
			name: "simple test",
			data: "Environment: prod\n",
			want: struct{appConfig AppConfig}{
				appConfig: AppConfig{
					Env: Prod,
				},
			},
		},
	}


	for _, test := range teseCases {
		t.Run(test.name,func(t *testing.T) {

			file, err := os.CreateTemp("", "storage*.json")
			require.NoError(t, err, "stop test with error: ", err)

			_, err = file.Write([]byte(test.data))
			require.NoError(t, err, "stop test with error: ", err)
			file.Close()

			appConfig := AppConfig {
				Env: Debug,
			}

			err = appConfig.getAppConfigFromFile(file.Name())
			require.NoError(t, err, "stop test with error: ", err)

			assert.Equal(t, test.want.appConfig, appConfig)

			err = os.Remove(file.Name())
			require.NoError(t, err, "stop test with error: ", err)
		})
	}
}