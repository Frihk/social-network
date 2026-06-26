CREATE TABLE IF NOT EXISTS posts (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER NOT NULL,
    group_id     INTEGER,
    content      TEXT NOT NULL,
    image_path   TEXT,
    privacy      TEXT NOT NULL CHECK(privacy IN ('public', 'almost_private', 'private')),
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)  REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS comments (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id      INTEGER NOT NULL,
    user_id      INTEGER NOT NULL,
    content      TEXT NOT NULL,
    image_path   TEXT,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS post_allowed_viewers (
    post_id  INTEGER NOT NULL,
    user_id  INTEGER NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
