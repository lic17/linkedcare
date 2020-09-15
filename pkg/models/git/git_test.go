package git

import (
	"testing"
)

func TestGitReadVerifyWithBasicAuth(t *testing.T) {
	shouldSuccess := []map[string]string{
		{
			"username": "",
			"password": "",
			"remote":   "https://github.com/linkedcare/linkedcare",
		},
	}
	shouldFailed := []map[string]string{
		{
			"username": "",
			"password": "",
			"remote":   "https://github.com/linkedcare/linkedcare12222",
		},
		{
			"username": "",
			"password": "",
			"remote":   "git@github.com:linkedcare/linkedcare.git",
		},
		{
			"username": "runzexia",
			"password": "",
			"remote":   "git@github.com:linkedcare/linkedcare.git",
		},
		{
			"username": "",
			"password": "",
			"remote":   "git@fdsfs41342`@@@2414!!!!github.com:linkedcare/linkedcare.git",
		},
	}
	for _, item := range shouldSuccess {
		err := gitReadVerifyWithBasicAuth(item["username"], item["password"], item["remote"])
		if err != nil {

			t.Errorf("should could access repo [%s] with %s:%s, %v", item["username"], item["password"], item["remote"], err)
		}
	}

	for _, item := range shouldFailed {
		err := gitReadVerifyWithBasicAuth(item["username"], item["password"], item["remote"])
		if err == nil {
			t.Errorf("should could access repo [%s] with %s:%s ", item["username"], item["password"], item["remote"])
		}
	}
}
