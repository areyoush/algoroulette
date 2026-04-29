package repository

import (
	"database/sql"
	"fmt"

	"github.com/areyoush/algoroulette/internal/model"
)

type QuestionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// GetRandom returns a random question visible to the user (global OR owned by them),
// joined with their personal status/bookmark/notes for that question.
func (r *QuestionRepository) GetRandom(userID int, topic, difficulty string) (*model.Question, error) {
	query := `
		SELECT
			q.id, q.title, q.topic, q.difficulty, q.slug, q.description, q.is_global, q.owner_id,
			uqs.status, COALESCE(uqs.bookmarked, FALSE), uqs.notes
		FROM questions q
		LEFT JOIN user_question_status uqs
			ON uqs.question_id = q.id AND uqs.user_id = $1
		WHERE (q.is_global = TRUE OR q.owner_id = $1)
	`
	args := []any{userID}
	i := 2

	if topic != "" {
		query += fmt.Sprintf(" AND q.topic = $%d", i)
		args = append(args, topic)
		i++
	}

	if difficulty != "" {
		query += fmt.Sprintf(" AND q.difficulty = $%d", i)
		args = append(args, difficulty)
		i++
	}

	query += " ORDER BY RANDOM() LIMIT 1"

	row := r.db.QueryRow(query, args...)
	q := &model.Question{}
	err := row.Scan(
		&q.ID, &q.Title, &q.Topic, &q.Difficulty, &q.Slug, &q.Description, &q.IsGlobal, &q.OwnerID,
		&q.Status, &q.Bookmarked, &q.Notes,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return q, err
}

// Insert adds a user-owned question (not global).
func (r *QuestionRepository) Insert(userID int, q *model.Question) error {
	return r.db.QueryRow(
		`INSERT INTO questions (title, topic, difficulty, slug, is_global, owner_id)
		 VALUES ($1, $2, $3, $4, FALSE, $5)
		 RETURNING id`,
		q.Title, q.Topic, q.Difficulty, q.Slug, userID,
	).Scan(&q.ID)
}

// InsertBatch adds multiple user-owned questions.
func (r *QuestionRepository) InsertBatch(userID int, questions []model.Question) error {
	query := `INSERT INTO questions (title, topic, difficulty, slug, is_global, owner_id) VALUES `
	args := []any{}

	for i, q := range questions {
		base := i * 5
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, FALSE, $%d),", base+1, base+2, base+3, base+4, base+5)
		args = append(args, q.Title, q.Topic, q.Difficulty, q.Slug, userID)
	}

	query = query[:len(query)-1] // remove trailing comma
	_, err := r.db.Exec(query, args...)
	return err
}

// InsertGlobalBatch adds global questions (admin use only — called internally, not via HTTP).
func (r *QuestionRepository) InsertGlobalBatch(questions []model.Question) error {
	query := `INSERT INTO questions (title, topic, difficulty, slug, is_global, owner_id) VALUES `
	args := []any{}

	for i, q := range questions {
		base := i * 4
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, TRUE, NULL),", base+1, base+2, base+3, base+4)
		args = append(args, q.Title, q.Topic, q.Difficulty, q.Slug)
	}

	query = query[:len(query)-1]
	_, err := r.db.Exec(query, args...)
	return err
}

// DeleteAllForUser deletes only the questions owned by a specific user.
func (r *QuestionRepository) DeleteAllForUser(userID int) error {
	_, err := r.db.Exec("DELETE FROM questions WHERE owner_id = $1", userID)
	return err
}

// UpsertStatus sets or clears solved/skipped for a user on a question.
func (r *QuestionRepository) UpsertStatus(userID, questionID int, status *string) error {
	_, err := r.db.Exec(`
		INSERT INTO user_question_status (user_id, question_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, question_id)
		DO UPDATE SET status = EXCLUDED.status
	`, userID, questionID, status)
	return err
}

// UpsertBookmark sets bookmarked for a user on a question.
func (r *QuestionRepository) UpsertBookmark(userID, questionID int, bookmarked bool) error {
	_, err := r.db.Exec(`
		INSERT INTO user_question_status (user_id, question_id, bookmarked)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, question_id)
		DO UPDATE SET bookmarked = EXCLUDED.bookmarked
	`, userID, questionID, bookmarked)
	return err
}

// UpsertNotes sets notes for a user on a question.
func (r *QuestionRepository) UpsertNotes(userID, questionID int, notes *string) error {
	_, err := r.db.Exec(`
		INSERT INTO user_question_status (user_id, question_id, notes)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, question_id)
		DO UPDATE SET notes = EXCLUDED.notes
	`, userID, questionID, notes)
	return err
}