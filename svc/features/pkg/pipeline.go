package pkg

import (
	"fmt"
	"io"
	"strconv"

	"github.com/gocarina/gocsv"
)

type In struct {
	Season               string `csv:"season"`
	TimeOnIce            string `csv:"time_on_ice"`
	Assists              string `csv:"assists"`
	Goals                string `csv:"goals"`
	Pim                  string `csv:"pim"`
	Shots                string `csv:"shots"`
	Games                string `csv:"games"`
	Hits                 string `csv:"hits"`
	PowerPlayGoals       string `csv:"power_play_goals"`
	PowerPlayPoints      string `csv:"power_play_points"`
	PowerPlayTimeOnIce   string `csv:"power_play_time_on_ice"`
	EvenTimeOnIce        string `csv:"even_time_on_ice"`
	FaceOffPct           string `csv:"face_off_pct"`
	ShotPct              string `csv:"shot_pct"`
	GameWinningGoals     string `csv:"game_winning_goals"`
	OverTimeGoals        string `csv:"over_time_goals"`
	ShortHandedGoals     string `csv:"short_handed_goals"`
	ShortHandedPoints    string `csv:"short_handed_points"`
	ShortHandedTimeOnIce string `csv:"short_handed_time_on_ice"`
	Blocked              string `csv:"blocked"`
	PlusMinus            string `csv:"plus_minus"`
	Points               string `csv:"points"`
	Shifts               string `csv:"shifts"`
	TeamName             string `csv:"team_name"`
	LeagueName           string `csv:"league_name"`
}

// All time on ice is translated from mmm:ss to seconds
type Out struct {
	Season                    string `csv:"season"`
	SeasonSinceDraft          string `csv:"season_since_draft"`
	TimeOnIcePerGame          string `csv:"time_on_ice_per_game"`
	PointsPerGame             string `csv:"points_per_game"`
	GoalsPerGame              string `csv:"goals_per_game"`
	AssistsPerGame            string `csv:"assists_per_game"`
	PimsPerGame               string `csv:"pims_per_game"`
	ShotsPerGame              string `csv:"shots_per_game"`
	HitsPerGame               string `csv:"hits_per_game"`
	PowerPlayTimeOnIcePerGame string `csv:"power_play_time_on_ice_per_game"`
	PowerPlayGoalsPerGame     string `csv:"power_play_goals_per_game"`
	PowerPlayAssistsPerGame   string `csv:"power_play_assists_per_game"`
	PowerPlayPointsPerGame    string `csv:"power_play_points_per_game"`
	EvenTimeOnIcePerGame      string `csv:"even_time_on_ice_per_game"`
	// EvenGoalsPerGame          string `csv:"even_goals_per_game"`
	// EvenAssistsPerGame        string `csv:"even_assists_per_game"`
	// EvenPointsPerGame         string `csv:"even_points_per_game"`
	FaceOffPct                string `csv:"face_off_pct"`
	ShotPct                   string `csv:"shot_pct"`
	GameWinningGoalsPerGame   string `csv:"game_winning_goals_per_game"`
	OverTimeGoalsPerGame      string `csv:"over_time_goals_per_game"`
	ShortHandedTimeOnIce      string `csv:"short_handed_time_on_ice"`
	ShortHandedGoalsPerGame   string `csv:"short_handed_goals_per_game"`
	ShortHandedAssistsPerGame string `csv:"short_handed_assists_per_game"`
	ShortHandedPointsPerGame  string `csv:"short_handed_points_per_game"`
	BlockedPerGame            string `csv:"blocked_per_game"`
	PlusMinus                 string `csv:"plus_minus"`
	ShiftsPerGame             string `csv:"shifts_per_game"`
}

// FeatureExtract parses a csv file and encodes it into a feature vector
// features include things like goals per game, rather than just goals and games played
func FeatureExtract(r io.Reader, w io.Writer) error {
	in := []In{}

	if err := gocsv.Unmarshal(r, &in); err != nil {
		return fmt.Errorf("error unmarshalling csv: %w", err)
	}

	seasonDrafted := "20142015" // todo compute this

	outs := []Out{}
	for _, i := range in {
		if i.LeagueName != "National Hockey League" {
			continue
		}

		out := Out{
			Season:                    i.Season,
			SeasonSinceDraft:          seasonSinceDraft(i.Season, seasonDrafted),
			TimeOnIcePerGame:          convertTimeOnIce(i.TimeOnIce, i.Games),
			PointsPerGame:             convertNumberPerGame(i.Points, i.Games),
			GoalsPerGame:              convertNumberPerGame(i.Goals, i.Games),
			AssistsPerGame:            convertNumberPerGame(i.Assists, i.Games),
			PimsPerGame:               convertNumberPerGame(i.Pim, i.Games),
			ShotsPerGame:              convertNumberPerGame(i.Shots, i.Games),
			HitsPerGame:               convertNumberPerGame(i.Hits, i.Games),
			PowerPlayTimeOnIcePerGame: convertTimeOnIce(i.PowerPlayTimeOnIce, i.Games),
			PowerPlayGoalsPerGame:     convertNumberPerGame(i.PowerPlayGoals, i.Games),
			PowerPlayAssistsPerGame:   convertNumberPerGame(stringMinus(i.PowerPlayPoints, i.PowerPlayGoals), i.Games),
			PowerPlayPointsPerGame:    convertNumberPerGame(i.PowerPlayPoints, i.Games),
			EvenTimeOnIcePerGame:      convertTimeOnIce(i.EvenTimeOnIce, i.Games),
			// EvenGoalsPerGame:          "todo",
			// EvenAssistsPerGame:        "todo",
			// EvenPointsPerGame:         "todo",
			FaceOffPct:                i.FaceOffPct,
			ShotPct:                   i.ShotPct,
			GameWinningGoalsPerGame:   convertNumberPerGame(i.GameWinningGoals, i.Games),
			OverTimeGoalsPerGame:      convertNumberPerGame(i.OverTimeGoals, i.Games),
			ShortHandedTimeOnIce:      convertTimeOnIce(i.ShortHandedTimeOnIce, i.Games),
			ShortHandedGoalsPerGame:   convertNumberPerGame(i.ShortHandedGoals, i.Games),
			ShortHandedAssistsPerGame: convertNumberPerGame(stringMinus(i.ShortHandedPoints, i.ShortHandedGoals), i.Games),
			ShortHandedPointsPerGame:  convertNumberPerGame(i.ShortHandedPoints, i.Games),
			BlockedPerGame:            convertNumberPerGame(i.Blocked, i.Games),
			PlusMinus:                 i.PlusMinus,
			ShiftsPerGame:             convertNumberPerGame(i.Shifts, i.Games),
		}
		outs = append(outs, out)
	}

	c, err := gocsv.MarshalString(&outs)
	if err != nil {
		return fmt.Errorf("error marshalling csv: %w", err)
	}
	_, err = w.Write([]byte(c))
	if err != nil {
		return fmt.Errorf("error writing csv: %w", err)
	}

	return nil
}

func convertNumberPerGame(a, b string) string {
	return fmt.Sprintf("%.2f", convertStringToFloat(a)/convertStringToFloat(b))
}

// converts timeonice stats with games played to per game stats in seconds
func convertTimeOnIce(timeOnIce string, games string) string {
	toi := convertTimeToSeconds(timeOnIce)
	gf := convertStringToFloat(games)
	return fmt.Sprintf("%.2f", (toi / gf))
}

func convertTimeToSeconds(timeOnIce string) float64 {
	var minutes, seconds float64
	fmt.Sscanf(timeOnIce, "%f:%f", &minutes, &seconds)
	return minutes*60 + seconds
}

func convertStringToFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func stringMinus(a, b string) string {
	return fmt.Sprintf("%d", mustAtoi(a)-mustAtoi(b))
}

func mustAtoi(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func seasonSinceDraft(season, seasonDrafted string) string {
	if seasonDrafted == "" {
		return ""
	}
	// season represented as 20222023, we need first 4 digits
	season = season[:4]
	sinceDraft, err := strconv.Atoi(season) // 2019
	if err != nil {
		panic(err)
	}
	seasonDrafted = seasonDrafted[:4]
	draftYear, err := strconv.Atoi(seasonDrafted) // 2015
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", (sinceDraft-draftYear)+1)
}
