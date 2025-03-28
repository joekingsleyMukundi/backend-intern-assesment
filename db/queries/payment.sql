-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY id
LIMIT $1
OFFSET $1;


-- name: CreatePayment :one
INSERT INTO payments (
    owner, amount, status
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdatePayment :one
UPDATE payments
    set status = $2
WHERE id = $1
RETURNING *;
