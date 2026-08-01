-- +goose Up
ALTER TABLE users
ADD COLUMN is_chirpy_red boolean;
