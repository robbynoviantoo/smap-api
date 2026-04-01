package repository

import (
	"context"
	"database/sql"
	"smap-api/internal/model"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) SaveMessage(ctx context.Context, msg *model.Message) error {
	now := time.Now()
	query := `INSERT INTO messages (sender_id, recipient_id, content, is_read, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, msg.SenderID, msg.RecipientID, msg.Content, msg.IsRead, now, now)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	msg.ID = uint(id)
	msg.CreatedAt = now
	msg.UpdatedAt = now
	return nil
}

func (r *MessageRepository) GetChatHistory(ctx context.Context, userID1, userID2 uint, limit, offset int) ([]model.Message, error) {
	query := `
		SELECT id, sender_id, recipient_id, content, is_read, created_at, updated_at 
		FROM messages 
		WHERE (sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, userID1, userID2, userID2, userID1, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Content, &m.IsRead, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}
