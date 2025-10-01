CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    institution VARCHAR(100) NOT NULL,
    address VARCHAR(200) NOT NULL,
  	approx_address BOOLEAN NOT NULL,
    location GEOMETRY(Point, 4326),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE routes (
    id serial NOT NULL,
    name character varying(255) NOT NULL,
    geom geometry (LineString, 4326) NOT NULL,
    distance real NOT NULL,
    duration real NOT NULL
);

CREATE TABLE nodes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
  	route_id integer NOT NULL,
		sequence integer NOT NULL,
    location GEOMETRY(Point, 4326) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
