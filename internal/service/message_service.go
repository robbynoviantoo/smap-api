package service

import (
	"context"
	"smap-api/internal/model"
	"smap-api/internal/repository"
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SaveMessage(ctx context.Context, msg *model.Message) error {
	return s.repo.SaveMessage(ctx, msg)
}

func (s *MessageService) GetChatHistory(ctx context.Context, userID1, userID2 uint, limit, offset int) ([]model.Message, error) {
	return s.repo.GetChatHistory(ctx, userID1, userID2, limit, offset)
}
