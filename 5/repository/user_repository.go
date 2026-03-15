package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"practice5/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type UserFilter struct {
	ID            *int
	Name          *string
	Email         *string
	Gender        *string
	BirthDateFrom *time.Time
	BirthDateTo   *time.Time
}

func (r *UserRepository) GetPaginatedUsers(
	page, pageSize int,
	filter UserFilter,
	orderBy string,
	orderDir string,
) (model.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	allowedOrderColumns := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"gender":     "gender",
		"birth_date": "birth_date",
	}

	column, ok := allowedOrderColumns[strings.ToLower(orderBy)]
	if !ok || column == "" {
		column = "id"
	}

	dir := strings.ToUpper(orderDir)
	if dir != "DESC" {
		dir = "ASC"
	}

	whereParts := make([]string, 0)
	args := make([]interface{}, 0)

	addFilter := func(cond string, value interface{}) {
		whereParts = append(whereParts, cond)
		args = append(args, value)
	}

	if filter.ID != nil {
		addFilter("id = $"+fmt.Sprint(len(args)+1), *filter.ID)
	}
	if filter.Name != nil && *filter.Name != "" {
		addFilter("name ILIKE $"+fmt.Sprint(len(args)+1), "%"+*filter.Name+"%")
	}
	if filter.Email != nil && *filter.Email != "" {
		addFilter("email ILIKE $"+fmt.Sprint(len(args)+1), "%"+*filter.Email+"%")
	}
	if filter.Gender != nil && *filter.Gender != "" {
		addFilter("gender = $"+fmt.Sprint(len(args)+1), *filter.Gender)
	}
	if filter.BirthDateFrom != nil {
		addFilter("birth_date >= $"+fmt.Sprint(len(args)+1), *filter.BirthDateFrom)
	}
	if filter.BirthDateTo != nil {
		addFilter("birth_date <= $"+fmt.Sprint(len(args)+1), *filter.BirthDateTo)
	}

	whereSQL := ""
	if len(whereParts) > 0 {
		whereSQL = "WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM users " + whereSQL

	var totalCount int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return model.PaginatedResponse{}, fmt.Errorf("count users: %w", err)
	}

	query := fmt.Sprintf(
		"SELECT id, name, email, gender, birth_date FROM users %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereSQL,
		column,
		dir,
		len(args)+1,
		len(args)+2,
	)

	argsWithPagination := append(append([]interface{}{}, args...), pageSize, offset)

	rows, err := r.db.Query(query, argsWithPagination...)
	if err != nil {
		return model.PaginatedResponse{}, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return model.PaginatedResponse{}, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return model.PaginatedResponse{}, fmt.Errorf("rows err: %w", err)
	}

	return model.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *UserRepository) GetCommonFriends(userID1, userID2 int) ([]model.User, error) {
	if userID1 == userID2 {
		return []model.User{}, nil
	}

	query := `
SELECT u.id, u.name, u.email, u.gender, u.birth_date
FROM users u
JOIN user_friends f1 ON f1.friend_id = u.id AND f1.user_id = $1
JOIN user_friends f2 ON f2.friend_id = u.id AND f2.user_id = $2
ORDER BY u.id;
`

	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, fmt.Errorf("select common friends: %w", err)
	}
	defer rows.Close()

	var result []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, fmt.Errorf("scan common friend: %w", err)
		}
		result = append(result, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return result, nil
}
