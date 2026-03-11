package postgres

import (
	"Hades/internal/models"
	"fmt"
)

func buildQuery(options models.Options) (string, []any) {

	condition, args := buildCondition(options)

	query := `

	SELECT id, type, amount, date, category, description, created_at
	FROM items` + condition

	return query, args

}

func buildCondition(options models.Options) (string, []any) {

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

	if options.Sort != "" {
		condition += fmt.Sprintf(` ORDER BY date %s`, options.Sort)
	}

	return condition, args

}
