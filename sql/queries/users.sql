-- name: InsertUser :one
INSERT INTO users (
    name, 
    institution,
    address,
    approx_address,
    location
) VALUES (
    $1,
    $2,
    $3,
    $4,
    ST_SetSRID(ST_MakePoint(@longitude::double precision, @latitude::double precision), 4326)
)
RETURNING *;

-- name: GetUsers :many
SELECT 
    id,
    name,
    institution,
    address,
    ST_X(location)::double precision AS longitude,
    ST_Y(location)::double precision AS latitude
From users LIMIT $1;


-- name: InsertNode :one
INSERT INTO nodes (
    name, 
    route_id,
  	sequence,
    location
) VALUES (
    $1,
    $2,
    $3,
    ST_SetSRID(ST_MakePoint(@longitude::double precision, @latitude::double precision), 4326)
)
RETURNING *;