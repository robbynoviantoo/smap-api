package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"smap-api/internal/model"
	"smap-api/internal/repository"
	"time"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) GetAll(ctx context.Context) ([]model.Event, error) {
	return s.repo.GetAll(ctx)
}

func (s *EventService) Create(ctx context.Context, req *model.EventCreateRequest) (*model.Event, error) {
	if req.Title == "" {
		return nil, errors.New("title wajib diisi")
	}
	if req.Start == "" {
		return nil, errors.New("start wajib diisi")
	}

	start, err := parseDateTime(req.Start)
	if err != nil {
		return nil, fmt.Errorf("format start tidak valid: %w", err)
	}

	var end *time.Time
	if req.End != "" {
		t, err := parseDateTime(req.End)
		if err != nil {
			return nil, fmt.Errorf("format end tidak valid: %w", err)
		}
		end = &t
	}

	allDay := true
	if req.AllDay != nil {
		allDay = *req.AllDay
	}

	e := &model.Event{
		Title:  req.Title,
		Start:  &start,
		End:    end,
		AllDay: allDay,
		Color:  req.Color,
	}

	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *EventService) Update(ctx context.Context, id uint, req *model.EventUpdateRequest) (*model.Event, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("event not found")
	}

	start, err := parseDateTime(req.Start)
	if err != nil {
		return nil, fmt.Errorf("format start tidak valid: %w", err)
	}
	existing.Start = &start

	if req.End != "" {
		t, err := parseDateTime(req.End)
		if err != nil {
			return nil, fmt.Errorf("format end tidak valid: %w", err)
		}
		existing.End = &t
	} else {
		existing.End = nil
	}

	if req.AllDay != nil {
		existing.AllDay = *req.AllDay
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Color != "" {
		existing.Color = req.Color
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *EventService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("event not found")
	}
	return s.repo.Delete(ctx, id)
}

// SyncHolidays mengambil hari libur nasional Indonesia dari date.nager.at (gratis, tanpa API key)
// dan menyimpan yang belum ada ke database.
func (s *EventService) SyncHolidays(ctx context.Context) (int, error) {
	year := time.Now().Year()
	// date.nager.at: gratis, global, tidak perlu API key
	// Response: [{"date":"2026-01-01","localName":"Tahun Baru Masehi","name":"New Year's Day",...}]
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/ID", year)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return 0, fmt.Errorf("gagal menghubungi API hari libur: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API hari libur mengembalikan status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Struktur response dari date.nager.at
	var raw []struct {
		Date      string `json:"date"`       // "2026-01-01"
		LocalName string `json:"localName"`  // nama dalam bahasa lokal (Indonesia)
		Name      string `json:"name"`       // nama dalam bahasa Inggris (fallback)
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return 0, fmt.Errorf("gagal parse response: %w", err)
	}

	var events []model.Event
	for _, h := range raw {
		t, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue // skip invalid date
		}
		// Gunakan localName (Bahasa Indonesia) jika ada, fallback ke name
		title := h.LocalName
		if title == "" {
			title = h.Name
		}
		events = append(events, model.Event{
			Title:  title,
			Start:  &t,
			End:    &t,
			AllDay: true,
			Color:  "Danger",
		})
	}

	return s.repo.BulkCreateHolidays(ctx, events)
}

// parseDateTime mencoba beberapa format umum.
func parseDateTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("tidak dapat parse tanggal: %q", s)
}
