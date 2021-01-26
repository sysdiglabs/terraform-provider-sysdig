package secure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) GetPolicyAssignmentBundleByName(ctx context.Context, name string) (*PolicyAssignmentBundle, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.policyAssignmentBundleURL(name), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to find bundle %s: %s", name, err.Error())
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return nil, err
	}

	var policyAssignmentBundle PolicyAssignmentBundle
	err = json.Unmarshal(body, &policyAssignmentBundle)
	if err != nil {
		return nil, err
	}

	return &policyAssignmentBundle, nil
}

func (client *sysdigSecureClient) PutPolicyAssignmentBundle(ctx context.Context, bundle PolicyAssignmentBundle) (*PolicyAssignmentBundle, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.policyAssignmentBundleURL(bundle.Id), bundle.ToJSON())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s: %s", bundle.Id, err.Error())
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return nil, err
	}

	return &bundle, nil
}

func (client *sysdigSecureClient) policyAssignmentBundleURL(name string) string {
	return fmt.Sprintf("%s/api/scanning/v1/mappings?bundleId=%s", client.URL, name)
}
