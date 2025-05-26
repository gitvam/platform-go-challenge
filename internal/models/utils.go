package models

import (
	"encoding/json"
	"errors"
)

func DecodeAsset(raw map[string]interface{}) (Asset, error) {
	t, ok := raw["type"].(string)
	if !ok {
		return nil, errors.New("missing or invalid asset type")
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var a Asset
	switch AssetType(t) {
	case AssetTypeChart:
		var chart Chart
		if err := json.Unmarshal(data, &chart); err != nil {
			return nil, err
		}
		a = &chart
	case AssetTypeInsight:
		var insight Insight
		if err := json.Unmarshal(data, &insight); err != nil {
			return nil, err
		}
		a = &insight
	case AssetTypeAudience:
		var audience Audience
		if err := json.Unmarshal(data, &audience); err != nil {
			return nil, err
		}
		a = &audience
	default:
		return nil, errors.New("unknown asset type")
	}

	return a, nil
}
