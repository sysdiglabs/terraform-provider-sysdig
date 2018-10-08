package secure

import (
	"errors"
	"io/ioutil"
)

func (client *sysdigSecureClient) GetUserRulesFile() (UserRulesFile, error) {
	response, _ := client.doSysdigSecureRequest("GET", client.userRulesFileURL(), nil)
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return UserRulesFile{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return UserRulesFileFromJSON(body), nil
}

func (client *sysdigSecureClient) userRulesFileURL() string {
	return client.URL + "/api/settings/falco/userRulesFile"
}

func (client *sysdigSecureClient) UpdateUserRulesFile(userRulesFileRequest UserRulesFile) (UserRulesFile, error) {
	response, _ := client.doSysdigSecureRequest("PUT", client.userRulesFileURL(), userRulesFileRequest.ToJSON())
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return UserRulesFile{}, errors.New(string(body))
	}

	defer response.Body.Close()

	return UserRulesFileFromJSON(body), nil
}
