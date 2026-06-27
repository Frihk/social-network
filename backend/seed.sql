-- seed.sql
-- This script inserts dummy users, posts, and groups for testing purposes.
-- The password for all dummy users is "password".

-- 1. Insert Dummy Users
-- Using a standard bcrypt hash for "password"
INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, nickname, about_me, is_private, created_at, updated_at)
VALUES 
('dummy_user_1', 'alice@example.com', '$2a$10$wO./60H5.M4.E9C3R8Z6..n48.cM/F83Xh.y/XJ0m9H1C5L7yE7eO', 'Alice', 'Smith', '1990-01-01', 'alicesmith', 'Hello World! I am Alice.', false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('dummy_user_2', 'bob@example.com', '$2a$10$wO./60H5.M4.E9C3R8Z6..n48.cM/F83Xh.y/XJ0m9H1C5L7yE7eO', 'Bob', 'Jones', '1992-02-02', 'bobjones', 'Just a guy named Bob.', false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('dummy_user_3', 'charlie@example.com', '$2a$10$wO./60H5.M4.E9C3R8Z6..n48.cM/F83Xh.y/XJ0m9H1C5L7yE7eO', 'Charlie', 'Brown', '1995-03-03', 'cbrown', 'Private profile testing.', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 2. Insert Dummy Posts
INSERT INTO posts (user_id, content, privacy, created_at)
VALUES 
('dummy_user_1', 'Hello everyone! This is my first public post.', 'public', CURRENT_TIMESTAMP),
('dummy_user_2', 'Just chilling today, what is everyone up to?', 'public', CURRENT_TIMESTAMP),
('dummy_user_3', 'This is a private post, only my followers can see it.', 'almost_private', CURRENT_TIMESTAMP);

-- 3. Insert Dummy Groups
INSERT INTO groups (id, creator_id, title, description, created_at)
VALUES 
('dummy_group_1', 'dummy_user_1', 'Tech Enthusiasts', 'A group for discussing the latest in tech, gadgets, and software.', CURRENT_TIMESTAMP),
('dummy_group_2', 'dummy_user_2', 'Movie Buffs', 'For people who love watching and reviewing movies.', CURRENT_TIMESTAMP);
