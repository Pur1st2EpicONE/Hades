package postgres

import (
	"Hades/internal/models"
	"fmt"
)

func buildQuery(options models.Options) (string, []any) {

	condition, args := buildWhere(options)

	query := `

	SELECT id, type, amount, date, category, description, created_at
	FROM items` + condition + fmt.Sprintf(`
	ORDER BY %s %s`, options.SortBy, options.Sort)

	return query, args

}

func buildWhere(options models.Options) (string, []any) {

	condition := ` WHERE TRUE`
	args := []any{}
	argIndex := 1

	if options.Type != "" {
		condition += fmt.Sprintf(` AND type = $%d`, argIndex)
		args = append(args, options.Type)
		argIndex++
	}

	if !options.From.IsZero() {
		condition += fmt.Sprintf(` AND date >= $%d`, argIndex)
		args = append(args, options.From)
		argIndex++
	}

	if !options.To.IsZero() {
		condition += fmt.Sprintf(` AND date <= $%d`, argIndex)
		args = append(args, options.To)
		argIndex++
	}

	if options.Category != "" {
		condition += fmt.Sprintf(` AND category = $%d`, argIndex)
		args = append(args, options.Category)
		argIndex++
	}

	return condition, args

}
