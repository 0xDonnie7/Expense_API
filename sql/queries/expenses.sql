-- name: CreateExpense :one
INSERT INTO expenses (
    id,
    user_id,
    amount,
    category,
    description,
    date,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetExpenseByID :one
SELECT * FROM expenses WHERE id = $1 AND user_id = $2;

-- name: ListExpensesByUser :many
SELECT * FROM expenses WHERE user_id = $1 ORDER BY date DESC;

-- name: UpdateExpense :one
UPDATE expenses 
SET amount = $1, category = $2, description = $3, date = $4, updated_at = $5
WHERE id = $6 AND user_id = $7
RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = $1 AND user_id = $2;

-- name: ListExpensesByUserAndDateRange :many
SELECT * FROM expenses
WHERE user_id = $1
AND date >= $2
AND date <= $3
ORDER BY date DESC;

