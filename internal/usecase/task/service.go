package task

import (
	"context"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()
	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		DueDate:     normalized.DueDate,
		RRule:       normalized.RRule,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		DueDate:     normalized.DueDate,
		RRule:       normalized.RRule,
		UpdatedAt:   s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, ErrInvalidInput
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, ErrInvalidInput
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, ErrInvalidInput
	}

	if !input.Status.Valid() {
		return UpdateInput{}, ErrInvalidInput
	}

	return input, nil
}

func (s *Service) GetTasksWithOccurrences(ctx context.Context, from, to time.Time) ([]*taskdomain.Task, error) {
	tasks, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []*taskdomain.Task

	for _, task := range tasks {
		if task.RRule == "" {
			// Ordinary task: check if due_date is in range
			if !task.DueDate.IsZero() && task.DueDate.After(from) && task.DueDate.Before(to) {
				result = append(result, &task)
			}
			continue
		}

		dates, err := GenerateOccurrences(task.RRule, from, to)
		if err != nil {
			continue
		}

		for _, date := range dates {
			copy := task
			copy.ID = 0
			copy.DueDate = date
			copy.CreatedAt = date
			copy.UpdatedAt = date
			copy.RRule = ""
			result = append(result, &copy)
		}
	}

	return result, nil
}
