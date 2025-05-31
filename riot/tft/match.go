package tft

import (
	"fmt"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/internal"
	log "github.com/sirupsen/logrus"
	"time"
)

// MatchClient provides methods for match endpoints of the League of Legends TFT API.
type MatchClient struct {
	c *internal.Client
}

// MatchListFilter providing filtering options for GetMatchesByPUUID
type MatchListFilter struct {
	// The matchlist started storing timestamps on June 16th, 2021.
	// Any matches played before June 16th, 2021 won't be included in the results if the startTime filter is set.
	StartTime, EndTime time.Time
}

func (f *MatchListFilter) buildParams() string {
	if f == nil {
		return ""
	}

	params := ""

	if f.StartTime != (time.Time{}) {
		params += fmt.Sprintf("&startTime=%d", f.StartTime.Unix())
	}

	if f.EndTime != (time.Time{}) {
		params += fmt.Sprintf("&endTime=%d", f.EndTime.Unix())
	}

	return params
}

// GetMatchesByPUUID returns a list of match ids by PUUID
func (mc *MatchClient) GetMatchesByPUUID(puuid string, start, count int, filters ...*MatchListFilter) ([]string, error) {
	logger := mc.logger().WithField("method", "GetMatchesByPUUID")
	mc.c.Region = api.Region(api.RegionToRoute[mc.c.Region])
	url := fmt.Sprintf(endpointMatchesByPUUID, puuid, start, count)
	for _, f := range filters {
		url += f.buildParams()
	}
	var out []string
	if err := mc.c.GetInto(url, &out); err != nil {
		logger.Debug(err)
		return nil, err
	}
	return out, nil
}

// GetMatchByMatchID returns a match by matchID
func (mc *MatchClient) GetMatchByMatchID(matchId string) (*Match, error) {
	logger := mc.logger().WithField("method", "GetMatchByMatchID")
	mc.c.Region = api.Region(api.RegionToRoute[mc.c.Region])
	url := fmt.Sprintf(endpointMatchByMatchID, matchId)
	var out *Match
	if err := mc.c.GetInto(url, &out); err != nil {
		logger.Debug(err)
		return nil, err
	}
	return out, nil
}

func (mc *MatchClient) logger() log.FieldLogger {
	return mc.c.Logger().WithField("category", "match")
}
