-- name: InsertRoute :one
INSERT INTO routes (name, geom, distance, duration) VALUES (
    $1,
    ST_SetSRID(ST_LineFromEncodedPolyline(@polyline::text, 5), 4326),
    @distance,
    @duration
)
RETURNING id;
