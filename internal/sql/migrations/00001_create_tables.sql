-- +goose Up
-- +goose StatementBegin
-- Enable foreign key constraints (required for SQLite)
PRAGMA foreign_keys = ON;

-- Networks: top-level networking resource
CREATE TABLE networks (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

-- Subnets: must belong to a network
-- Note: Application should enforce that each network has at least one subnet
CREATE TABLE subnets (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	network_id INTEGER NOT NULL REFERENCES networks(id) ON DELETE CASCADE
);

-- Static IPs: can be assigned to servers or load balancers, can have domain names
CREATE TABLE static_ips (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	ip_address TEXT NOT NULL UNIQUE
);

-- Domain names: can be associated with static IPs
CREATE TABLE domain_names (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	static_ip_id INTEGER REFERENCES static_ips(id) ON DELETE SET NULL
);

-- Servers: must belong to a subnet, optionally can have a static IP
CREATE TABLE servers (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	subnet_id INTEGER NOT NULL REFERENCES subnets(id) ON DELETE CASCADE,
	static_ip_id INTEGER REFERENCES static_ips(id) ON DELETE SET NULL
);

-- Load balancers: must belong to a subnet, optionally can have a static IP
CREATE TABLE load_balancers (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	subnet_id INTEGER NOT NULL REFERENCES subnets(id) ON DELETE CASCADE,
	static_ip_id INTEGER REFERENCES static_ips(id) ON DELETE SET NULL
);

-- Databases: must belong to a subnet
CREATE TABLE databases (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	subnet_id INTEGER NOT NULL REFERENCES subnets(id) ON DELETE CASCADE
);

-- Buckets: standalone storage resource with no relationships
CREATE TABLE buckets (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS buckets;
DROP TABLE IF EXISTS databases;
DROP TABLE IF EXISTS load_balancers;
DROP TABLE IF EXISTS servers;
DROP TABLE IF EXISTS domain_names;
DROP TABLE IF EXISTS static_ips;
DROP TABLE IF EXISTS subnets;
DROP TABLE IF EXISTS networks;
-- +goose StatementEnd
