package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	acceptPostureRiskCreatePath = "%s/api/cspm/v1/compliance/risk-acceptances"
	acceptPostureRiskGetPath    = "%s/api/cspm/v1/compliance/risk-acceptances/%s"
	acceptPostureRiskDelete     = "%s/api/cspm/v1/compliance/violations/revoke"
	acceptPostureRiskUpdate     = "%s/api/cspm/v1/compliance/risk-acceptances/%s"
)

type PostureAcceptRiskInterface interface {
	Base
	SaveAcceptPostureRisk(ctx context.Context, p *AccepetPostureRiskRequest) (*AcceptPostureRiskResponse, string, error)
	GetAcceptancePostureRisk(ctx context.Context, id string) (*AcceptPostureRiskResponse, string, error)
	DeleteAcceptancePostureRisk(ctx context.Context, p *DeleteAcceptPostureRisk) error
	UpdateAcceptancePostureRisk(ctx context.Context, p *UpdateAccepetPostureRiskRequest) (*AcceptPostureRisk, string, error)
}

func (c *Client) SaveAcceptPostureRisk(ctx context.Context, p *AccepetPostureRiskRequest) (risk *AcceptPostureRiskResponse, errString string, err error) {
	payload, err := Marshal(p)
	if err != nil {
		return nil, "", err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.getPostureControlURL(acceptPostureRiskCreatePath), payload)
	if err != nil {
		return nil, "", err
	}

	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}
	resp, err := Unmarshal[AcceptPostureRiskResponse](response.Body)
	if err != nil {
		return nil, "", err
	}

	return &resp, "", nil
}

func (c *Client) GetAcceptancePostureRisk(ctx context.Context, id string) (risk *AcceptPostureRiskResponse, errString string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(acceptPostureRiskGetPath, c.config.url, id), nil)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}

	wrapper, err := Unmarshal[AcceptPostureRiskResponse](response.Body)
	if err != nil {
		return nil, "", err
	}
	return &wrapper, "", nil
}

func (c *Client) DeleteAcceptancePostureRisk(ctx context.Context, p *DeleteAcceptPostureRisk) (err error) {
	payload, err := Marshal(p)
	if err != nil {
		return err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, fmt.Sprintf(acceptPostureRiskDelete, c.config.url), payload)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) UpdateAcceptancePostureRisk(ctx context.Context, p *UpdateAccepetPostureRiskRequest) (risk *AcceptPostureRisk, errString string, err error) {
	payload, err := Marshal(p)
	if err != nil {
		return nil, "", err
	}
	response, err := c.requester.Request(ctx, http.MethodPatch, fmt.Sprintf(acceptPostureRiskUpdate, c.config.url, p.AcceptanceID), payload)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		errStatus, err := c.ErrorAndStatusFromResponse(response)
		return nil, errStatus, err
	}
	resp, err := Unmarshal[AcceptPostureRiskResponse](response.Body)
	if err != nil {
		return nil, "", err
	}

	return &resp.Data, "", nil
}
