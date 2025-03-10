// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: recipe.sql

package store

import (
	"context"
)

const createRecipe = `-- name: CreateRecipe :one
INSERT INTO recipe (id, name, description) VALUES (?, ?, ?) RETURNING id, name, description
`

type CreateRecipeParams struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) CreateRecipe(ctx context.Context, arg CreateRecipeParams) (Recipe, error) {
	row := q.db.QueryRowContext(ctx, createRecipe, arg.ID, arg.Name, arg.Description)
	var i Recipe
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const getRecipe = `-- name: GetRecipe :one
SELECT id, name, description FROM recipe WHERE id = ?
`

func (q *Queries) GetRecipe(ctx context.Context, id string) (Recipe, error) {
	row := q.db.QueryRowContext(ctx, getRecipe, id)
	var i Recipe
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const getRecipeWithIngredients = `-- name: GetRecipeWithIngredients :one
SELECT recipe.id, recipe.name, recipe.description, recipe_id, ingredient_id, quantity, unit, ingredient.id, ingredient.name, ingredient.description FROM recipe
JOIN dosing ON recipe.id = dosing.recipe_id
JOIN ingredient ON dosing.ingredient_id = ingredient.id
WHERE recipe.id = ?
`

type GetRecipeWithIngredientsRow struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	RecipeID      string `json:"recipe_id"`
	IngredientID  string `json:"ingredient_id"`
	Quantity      int64  `json:"quantity" validate:"required,gt=0"`
	Unit          string `json:"unit"`
	ID_2          string `json:"id_2"`
	Name_2        string `json:"name_2"`
	Description_2 string `json:"description_2"`
}

func (q *Queries) GetRecipeWithIngredients(ctx context.Context, id string) (GetRecipeWithIngredientsRow, error) {
	row := q.db.QueryRowContext(ctx, getRecipeWithIngredients, id)
	var i GetRecipeWithIngredientsRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.RecipeID,
		&i.IngredientID,
		&i.Quantity,
		&i.Unit,
		&i.ID_2,
		&i.Name_2,
		&i.Description_2,
	)
	return i, err
}

const getRecipes = `-- name: GetRecipes :many
SELECT id, name, description FROM recipe
`

func (q *Queries) GetRecipes(ctx context.Context) ([]Recipe, error) {
	rows, err := q.db.QueryContext(ctx, getRecipes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
