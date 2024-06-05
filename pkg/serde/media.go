package serde

import (
	"mime"
	"sort"
	"strconv"
	"strings"

	"github.com/google/jsonapi"
)

type MediaType string

const (
	JSONAPI MediaType = jsonapi.MediaType
	JSON    MediaType = "application/json"
	Text    MediaType = "text/plain"
	YAML    MediaType = "application/yaml"
	Default MediaType = JSON
)

func MediaTypeFor(value string) MediaType {
	mediaTypes := sortedByQualityValues(strings.Split(value, ","))

	for _, mediaType := range mediaTypes {
		switch mediaType {
		case jsonapi.MediaType:
			return JSONAPI
		case "application/json":
			return JSON
		case "text/plain":
			return Text
		default:
			continue
		}
	}
	return Default
}

func sortedByQualityValues(mimeTypes []string) []string {
	weights := map[string]float64{}
	valid := []string{}

	for i, mimeType := range mimeTypes {
		if mtype, params, err := mime.ParseMediaType(mimeType); err == nil {
			valid = append(valid, mtype)
			if _, ok := weights[mtype]; !ok {
				weights[mtype] = 1.0 - (0.1 * float64(i))
			}
			if quality, ok := params["q"]; ok {
				if val, err := strconv.ParseFloat(quality, 64); err == nil {
					weights[mtype] = val
				}
			}
		}
	}
	sort.Slice(valid, func(x, y int) bool {
		xWeight := weights[valid[x]]
		yWeight := weights[valid[y]]
		return xWeight > yWeight
	})
	return valid
}
