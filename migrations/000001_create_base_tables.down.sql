-- Drop tables in reverse order of creation (due to foreign key constraints)
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS action_types;
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS shares;
DROP TABLE IF EXISTS shopping_list_items;
DROP TABLE IF EXISTS announcements;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS storage;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS status;
DROP TABLE IF EXISTS roles;