package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/mimir/pkg/distributor"
	"github.com/pkg/errors"
)

func (c *MimirClient) Backfill(ctx context.Context, source string, tenantID int, logger log.Logger) error {
	bfReq := distributor.BackfillRequest{
		TenantID: tenantID,
	}
	payload, err := json.Marshal(bfReq)
	if err != nil {
		return errors.Wrap(err, "failed to JSON encode backfill request")
	}

	res, err := c.doRequest("/api/v1/backfill", http.MethodPost, payload)
	if err != nil {
		return errors.Wrap(err, "backfill request failed")
	}
	defer res.Body.Close()

	var bfResp struct {
		Token string `json:"token"`
	}
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&bfResp); err != nil {
		return errors.Wrap(err, "failed to decode backfill response")
	}

	level.Info(logger).Log("msg", "Backfill request successful", "token", bfResp.Token)

	return nil
}
