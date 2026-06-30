-- MCP InSales - Database Migrations
-- Этот файл выполняется при каждом обновлении deb-пакета.
-- Используйте ADD COLUMN IF NOT EXISTS для безопасного повторного выполнения.

-- Миграция v0.1.0 (начальная версия)
-- (schema.sql уже содержит базовую схему)

-- Миграция v0.2.0 - Per-shop read-only
ALTER TABLE shops ADD COLUMN IF NOT EXISTS is_read_only BOOLEAN DEFAULT FALSE;

-- Миграция v0.3.0 (будущее)
-- ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
