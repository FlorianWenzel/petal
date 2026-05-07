package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var validMoods = map[string]bool{
	"bright": true,
	"tender": true,
	"calm":   true,
	"heavy":  true,
	"stormy": true,
}

type entryReq struct {
	Mood string `json:"mood"`
	Text string `json:"text"`
}

type entryOut struct {
	Mood string  `json:"mood"`
	Text *string `json:"text,omitempty"`
}

type monthEntries map[string]entryOut

func (s *server) loadMonth(w http.ResponseWriter, r *http.Request) {
	u := userFromCtx(r.Context())
	year, month, ok := parseYearMonth(r)
	if !ok {
		writeErr(w, http.StatusBadRequest, "invalid year/month")
		return
	}
	out, err := loadMonthEntries(r.Context(), s.pool, u.ID, year, month)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *server) saveDay(w http.ResponseWriter, r *http.Request) {
	u := userFromCtx(r.Context())
	year, month, ok := parseYearMonth(r)
	if !ok {
		writeErr(w, http.StatusBadRequest, "invalid year/month")
		return
	}
	day, err := strconv.Atoi(r.PathValue("day"))
	if err != nil || day < 1 || day > 31 {
		writeErr(w, http.StatusBadRequest, "invalid day")
		return
	}

	var body entryReq
	if err := decodeJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	if !validMoods[body.Mood] {
		writeErr(w, http.StatusBadRequest, "invalid mood")
		return
	}
	if len(body.Text) > 10_000 {
		writeErr(w, http.StatusBadRequest, "text too long")
		return
	}

	dayStr, ok := safeDate(year, month, day)
	if !ok {
		writeErr(w, http.StatusBadRequest, "invalid date")
		return
	}

	var textVal any
	if strings.TrimSpace(body.Text) == "" {
		textVal = nil
	} else {
		textVal = body.Text
	}

	if _, err := s.pool.Exec(r.Context(),
		`insert into entries (user_id, day, mood, text)
		 values ($1, $2::date, $3, $4)
		 on conflict (user_id, day) do update set
		   mood = excluded.mood,
		   text = excluded.text,
		   updated_at = now()`,
		u.ID, dayStr, body.Mood, textVal,
	); err != nil {
		writeErr(w, http.StatusInternalServerError, "save")
		return
	}

	out, err := loadMonthEntries(r.Context(), s.pool, u.ID, year, month)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "load")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func loadMonthEntries(ctx context.Context, pool *pgxpool.Pool, userID string, year, month int) (monthEntries, error) {
	start := fmt.Sprintf("%04d-%02d-01", year, month)
	nextYear, nextMonth := year, month+1
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	end := fmt.Sprintf("%04d-%02d-01", nextYear, nextMonth)

	rows, err := pool.Query(ctx,
		`select day, mood, text from entries
		  where user_id = $1 and day >= $2::date and day < $3::date`,
		userID, start, end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := monthEntries{}
	for rows.Next() {
		var day time.Time
		var mood string
		var text *string
		if err := rows.Scan(&day, &mood, &text); err != nil {
			return nil, err
		}
		out[strconv.Itoa(day.UTC().Day())] = entryOut{Mood: mood, Text: text}
	}
	return out, rows.Err()
}

func parseYearMonth(r *http.Request) (year, month int, ok bool) {
	y, err := strconv.Atoi(r.PathValue("year"))
	if err != nil {
		return 0, 0, false
	}
	m, err := strconv.Atoi(r.PathValue("month"))
	if err != nil || m < 1 || m > 12 {
		return 0, 0, false
	}
	return y, m, true
}

func safeDate(year, month, day int) (string, bool) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return "", false
	}
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day), true
}
