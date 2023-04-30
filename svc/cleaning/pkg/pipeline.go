package pkg

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
)

func JSONToCSV(r io.Reader, w io.Writer) error {
	resp := Response{}
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return err
	}

	csvw := csv.NewWriter(w)
	defer csvw.Flush()

	header := []string{"season", "time_on_ice", "assists", "goals", "pim", "shots", "games", "hits", "power_play_goals", "power_play_points", "power_play_time_on_ice", "even_time_on_ice", "penalty_minutes", "face_off_pct", "shot_pct", "game_winning_goals", "over_time_goals", "short_handed_goals", "short_handed_points", "short_handed_time_on_ice", "blocked", "plus_minus", "points", "shifts", "team_name", "league_name"}
	if err := csvw.Write(header); err != nil {
		return err
	}

	for _, stat := range resp.Stats {
		for _, split := range stat.Splits {
			row := []string{
				split.Season,
				split.Stat.TimeOnIce,
				optionalIntToString(split.Stat.Assists),
				optionalIntToString(split.Stat.Goals),
				optionalIntToString(split.Stat.Pim),
				optionalIntToString(split.Stat.Shots),
				optionalIntToString(split.Stat.Games),
				optionalIntToString(split.Stat.Hits),
				optionalIntToString(split.Stat.PowerPlayGoals),
				optionalIntToString(split.Stat.PowerPlayPoints),
				split.Stat.PowerPlayTimeOnIce,
				split.Stat.EvenTimeOnIce,
				split.Stat.PenaltyMinutes,
				optionalFloatToString(split.Stat.FaceOffPct),
				optionalFloatToString(split.Stat.ShotPct),
				optionalIntToString(split.Stat.GameWinningGoals),
				optionalIntToString(split.Stat.OverTimeGoals),
				optionalIntToString(split.Stat.ShortHandedGoals),
				optionalIntToString(split.Stat.ShortHandedPoints),
				split.Stat.ShortHandedTimeOnIce,
				optionalIntToString(split.Stat.Blocked),
				optionalIntToString(split.Stat.PlusMinus),
				optionalIntToString(split.Stat.Points),
				optionalIntToString(split.Stat.Shifts),
				split.Team.Name,
				split.League.Name,
			}
			if err := csvw.Write(row); err != nil {
				return err
			}
		}
	}
	return nil
}

func optionalIntToString(in *int) string {
	if in != nil {
		return strconv.Itoa(*in)
	}
	return ""
}

func optionalFloatToString(in *float64) string {
	if in != nil {
		return strconv.FormatFloat(*in, 'f', 2, 64)
	}
	return ""
}

type Response struct {
	Stats []Stats `json:"stats"`
}

type Type struct {
	DisplayName string `json:"displayName"`
	GameType    any    `json:"gameType"`
}

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Stat struct {
	TimeOnIce            string   `json:"timeOnIce,omitempty"`
	Assists              *int     `json:"assists,omitempty"`
	Goals                *int     `json:"goals,omitempty"`
	Pim                  *int     `json:"pim,omitempty"`
	Shots                *int     `json:"shots,omitempty"`
	Games                *int     `json:"games,omitempty"`
	Hits                 *int     `json:"hits,omitempty"`
	PowerPlayGoals       *int     `json:"powerPlayGoals,omitempty"`
	PowerPlayPoints      *int     `json:"powerPlayPoints,omitempty"`
	PowerPlayTimeOnIce   string   `json:"powerPlayTimeOnIce,omitempty"`
	EvenTimeOnIce        string   `json:"evenTimeOnIce,omitempty"`
	PenaltyMinutes       string   `json:"penaltyMinutes,omitempty"`
	FaceOffPct           *float64 `json:"faceOffPct,omitempty"`
	ShotPct              *float64 `json:"shotPct,omitempty"`
	GameWinningGoals     *int     `json:"gameWinningGoals,omitempty"`
	OverTimeGoals        *int     `json:"overTimeGoals,omitempty"`
	ShortHandedGoals     *int     `json:"shortHandedGoals,omitempty"`
	ShortHandedPoints    *int     `json:"shortHandedPoints,omitempty"`
	ShortHandedTimeOnIce string   `json:"shortHandedTimeOnIce,omitempty"`
	Blocked              *int     `json:"blocked,omitempty"`
	PlusMinus            *int     `json:"plusMinus,omitempty"`
	Points               *int     `json:"points,omitempty"`
	Shifts               *int     `json:"shifts,omitempty"`
}

type Splits struct {
	Season string `json:"season"`
	Stat   Stat   `json:"stat,omitempty"`
	Team   Team   `json:"team,omitempty"`
	League League `json:"league,omitempty"`
}

type Stats struct {
	Type   Type     `json:"type"`
	Splits []Splits `json:"splits"`
}
