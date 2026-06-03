# Social Network — Issues

---

## Issue #1 — Authentication & Profiles
**Assigned to:** Member 1

You are responsible for the foundation of the entire project. Every other member depends on your work — nothing can be tested or built until the database is running, users can register and log in, and the auth middleware is in place. Build the backend first, then move to the frontend.

### What you are building

**On the backend**, you will set up the SQLite database and write the migration runner so the database tables are created automatically every time the server starts. You will write the `users` table migration with all required columns: `id` (UUID), `email`, `password_hash`, `first_name`, `last_name`, `date_of_birth`, `avatar_path`, `nickname`, `about_me`, `is_private`, and `created_at`. You will then write the Go model functions that query this table (`CreateUser`, `GetUserByEmail`, `GetUserByID`, `UpdateUser`, `SetProfilePrivacy`).

Once the model is ready, you will write the auth handlers. The register handler (`POST /api/auth/register`) takes a user's details, hashes the password with bcrypt, saves the user, and returns the new user's ID. The login handler (`POST /api/auth/login`) checks the credentials, creates a UUID session, stores it, and sets it as an `HttpOnly` cookie named `session_token` on the response. The logout handler (`POST /api/auth/logout`) deletes the session and clears the cookie.

You will also write two middleware functions. The auth middleware reads the `session_token` cookie, validates it, and attaches the `user_id` to the request context — every protected route uses this. The CORS middleware allows the Next.js frontend (running on port 3000) to make requests to the backend, including sending cookies cross-origin. Both must be registered in `server.go` before any other member can test their work.

Finally, write the user profile handler (`GET /api/users/:id` and `PUT /api/users/:id/privacy`) so the frontend can fetch and update profile data.

**On the frontend**, you will build the login page (`/login`) and register page (`/register`). Both are forms — login takes email and password, register takes all user fields including an optional avatar upload. On success, both redirect appropriately. You will set up `AuthContext`, which calls `GET /api/auth/me` on app load to check whether a session cookie is valid and stores the current user globally. You will also write `middleware.js` which redirects unauthenticated users away from protected pages and redirects logged-in users away from the login/register pages. Finally, build the profile page (`/profile/[id]`) which shows the user's info, their posts, follower/following counts, and a public/private toggle that only appears on the logged-in user's own profile.

### Checklist

#### Backend
- [ ] `backend/pkg/db/sqlite/sqlite.go` — connect to SQLite and run all migrations on startup
- [ ] `000001_create_users_table.up.sql` — create the `users` table with all required columns
- [ ] `000001_create_users_table.down.sql` — drop the `users` table
- [ ] `backend/pkg/models/user.go` — `CreateUser`, `GetUserByEmail`, `GetUserByID`, `UpdateUser`, `SetProfilePrivacy`
- [ ] `POST /api/auth/register` — validate fields, hash password with bcrypt, save user, return `201`
- [ ] `POST /api/auth/login` — verify password, create UUID session, set `HttpOnly` cookie, return `200`
- [ ] `POST /api/auth/logout` — delete session, clear cookie, return `200`
- [ ] `GET /api/auth/me` — return the logged-in user's data using the session cookie
- [ ] `GET /api/users/:id` — return a user's profile (respects public/private logic)
- [ ] `PUT /api/users/:id/privacy` — toggle `is_private` for the logged-in user
- [ ] `backend/pkg/middleware/auth.go` — validate session cookie, attach `user_id` to context, return `401` if missing
- [ ] `backend/pkg/middleware/cors.go` — allow origin `http://localhost:3000`, allow credentials
- [ ] Register all routes and both middleware in `server.go`

#### Frontend
- [ ] `src/context/AuthContext.jsx` — call `GET /api/auth/me` on load, store user, expose `login()` and `logout()`
- [ ] `src/middleware.js` — redirect unauthenticated users to `/login`, redirect logged-in users away from `/login` and `/register`
- [ ] `src/app/login/page.jsx` — email + password form, calls `POST /api/auth/login`, redirects to `/` on success
- [ ] `src/app/register/page.jsx` — full registration form with optional avatar upload, redirects to `/login` on success
- [ ] `src/app/profile/[id]/page.jsx` — show avatar, name, bio, posts, follower/following counts, public/private toggle (own profile only)
- [ ] `src/lib/auth.js` — fetch wrapper functions for all auth endpoints

---

## Issue #2 — Posts & Followers
**Assigned to:** Member 2

You are building the core social features — the ability to create posts, comment on them, and follow other users. You can start writing migrations and models immediately, but you will not be able to test your handlers until Member 1 has finished the auth middleware.

### What you are building

**On the backend**, you will write two migrations. The posts migration creates a `posts` table (with `id`, `user_id`, `content`, `image_path`, `privacy`, `created_at`), a `comments` table (with `id`, `post_id`, `user_id`, `content`, `image_path`, `created_at`), and a `post_allowed_viewers` table (with `post_id` and `user_id`) for posts with `privacy = 'private'`. The followers migration creates a `followers` table with `id`, `follower_id`, `following_id`, `status` (either `pending` or `accepted`), and `created_at`, with a unique constraint on the pair `(follower_id, following_id)`.

You will also write the image upload utility, a shared function used by posts, comments, and the profile avatar. It validates that uploaded files are JPEG, PNG, or GIF, saves them to `backend/uploads/` with a UUID filename, and returns the saved path.

The posts handlers cover the full feed logic: `GET /api/posts` returns posts the logged-in user is allowed to see (public posts, followers-only posts if they follow the author, and private posts if they are in `post_allowed_viewers`). `POST /api/posts` creates a post with the specified privacy level. `GET /api/users/:id/posts` returns all posts by a specific user. Comment endpoints allow creating and listing comments on any post.

The followers handlers manage the follow relationship. `POST /api/users/:id/follow` immediately accepts on a public profile, or creates a `pending` request and a notification on a private profile. `PUT /api/followers/:id/accept` and `decline` let the recipient respond. `POST /api/users/:id/unfollow` removes the relationship.

**On the frontend**, you will build the home feed page (`/`), which fetches and displays posts using the `PostCard` component. `PostCard` shows the author, content, image, privacy badge, timestamp, and a comment thread that expands on click. `PostForm` is the create-post form with a content textarea, a privacy dropdown (Public, Followers only, Specific followers), a multi-select follower list that appears when "Specific followers" is chosen, and an image upload input. The followers panel lives inside the profile page — it shows a Follow/Unfollow button and, on your own profile, a list of pending follow requests with Accept and Decline buttons.

### Checklist

#### Backend
- [ ] `000002_create_posts_table.up.sql` — create `posts`, `comments`, and `post_allowed_viewers` tables
- [ ] `000002_create_posts_table.down.sql` — drop all three tables
- [ ] `000004_create_followers_table.up.sql` — create `followers` table with unique constraint on `(follower_id, following_id)`
- [ ] `000004_create_followers_table.down.sql` — drop `followers` table
- [ ] `backend/pkg/utils/images.go` — validate JPEG/PNG/GIF, save with UUID filename, return path
- [ ] `backend/pkg/models/post.go` — model functions for posts and comments
- [ ] `backend/pkg/models/follower.go` — model functions for follower relationships
- [ ] `GET /api/posts` — return feed filtered by privacy rules, ordered newest first
- [ ] `POST /api/posts` — create post, handle image upload, save `post_allowed_viewers` if privacy is `private`
- [ ] `GET /api/users/:id/posts` — return all posts by a specific user
- [ ] `POST /api/posts/:id/comments` — add a comment (with optional image)
- [ ] `GET /api/posts/:id/comments` — return all comments on a post
- [ ] `POST /api/users/:id/follow` — accept immediately if public, create `pending` + notification if private
- [ ] `PUT /api/followers/:id/accept` — update status to `accepted`
- [ ] `PUT /api/followers/:id/decline` — delete the follower row
- [ ] `POST /api/users/:id/unfollow` — delete the follower row
- [ ] `GET /api/users/:id/followers` — list accepted followers
- [ ] `GET /api/users/:id/following` — list users being followed

#### Frontend
- [ ] `src/app/page.jsx` — fetch `GET /api/posts`, render a `PostCard` for each, show empty state if none
- [ ] `src/components/PostCard.jsx` — author info, content, image, privacy badge, timestamp, expandable comment thread
- [ ] `src/components/PostForm.jsx` — content textarea, privacy dropdown, conditional follower multi-select, image upload, submit
- [ ] Followers panel inside `src/app/profile/[id]/page.jsx` — Follow/Unfollow button, pending requests with Accept/Decline (own profile only)
- [ ] `src/lib/posts.js` — fetch wrappers for all post and comment endpoints
- [ ] `src/lib/followers.js` — fetch wrappers for all follower endpoints

---

## Issue #3 — Groups & Events
**Assigned to:** Member 3

You are building the groups feature — users can create groups, invite or request members, post inside groups, and create events with RSVP. You can write migrations and models immediately, but you need Member 1's auth middleware before testing handlers, and Member 2's post components before building the group detail page.

### What you are building

**On the backend**, you will write two migrations. The groups migration creates a `groups` table (`id`, `creator_id`, `title`, `description`, `created_at`) and a `group_members` table (`id`, `group_id`, `user_id`, `status` of `invited`/`requested`/`accepted`, `created_at`) with a unique constraint on `(group_id, user_id)`. The events migration creates an `events` table (`id`, `group_id`, `creator_id`, `title`, `description`, `event_time`, `created_at`) and an `event_responses` table (`id`, `event_id`, `user_id`, `response` of `going`/`not_going`) with a unique constraint on `(event_id, user_id)` so each user can only respond once per event.

The groups handlers cover: `GET /api/groups` (browse all groups), `POST /api/groups` (create a group), `GET /api/groups/:id` (group details and members, only full data for accepted members), `POST /api/groups/:id/invite` (invite a user, triggers a notification), `POST /api/groups/:id/request` (request to join, triggers a notification to the creator), `PUT /api/groups/:id/members/:userId/accept` and `decline` (creator only — accepting or declining a join request), and `PUT /api/group-invites/:id/accept` and `decline` (for the invited user to respond).

The events handlers cover: `GET /api/groups/:id/events` (list events), `POST /api/groups/:id/events` (create an event, triggers a notification to all group members), and `POST /api/events/:id/respond` (submit or update a going/not-going response using an upsert).

**On the frontend**, you will build the groups browse page (`/groups`) which lists all groups using `GroupCard` and includes a "Create Group" button that opens a modal. `GroupCard` shows the group title, description, member count, and a "Request to Join" button for groups the user has not joined. The group detail page (`/groups/[id]`) shows the full group — title, description, member list, a post feed (reusing `PostCard` and `PostForm` from Member 2), an events section, an "Invite Member" button, and a pending requests section visible only to the creator. The events widget inside the detail page lists events with their RSVP counts and Going/Not going buttons, and includes a "Create Event" form for members.

### Checklist

#### Backend
- [ ] `000003_create_groups_table.up.sql` — create `groups` and `group_members` tables with unique constraint
- [ ] `000003_create_groups_table.down.sql` — drop both tables
- [ ] `000006_create_events_table.up.sql` — create `events` and `event_responses` tables with unique constraint
- [ ] `000006_create_events_table.down.sql` — drop both tables
- [ ] `backend/pkg/models/group.go` — model functions for groups and membership
- [ ] `GET /api/groups` — return all groups
- [ ] `POST /api/groups` — create a group, set creator as first accepted member
- [ ] `GET /api/groups/:id` — return group details; return full data only to accepted members
- [ ] `POST /api/groups/:id/invite` — create `invited` membership row, send notification to invited user
- [ ] `POST /api/groups/:id/request` — create `requested` membership row, send notification to group creator
- [ ] `PUT /api/groups/:id/members/:userId/accept` — set status to `accepted` (creator only)
- [ ] `PUT /api/groups/:id/members/:userId/decline` — delete membership row (creator only)
- [ ] `PUT /api/group-invites/:id/accept` — set status to `accepted` (invited user only)
- [ ] `PUT /api/group-invites/:id/decline` — delete membership row (invited user only)
- [ ] `GET /api/groups/:id/events` — return all events for a group (members only)
- [ ] `POST /api/groups/:id/events` — create event, send notification to all accepted group members
- [ ] `POST /api/events/:id/respond` — upsert `going` or `not_going` response (members only)

#### Frontend
- [ ] `src/app/groups/page.jsx` — fetch and list all groups with `GroupCard`, "Create Group" modal with title + description form
- [ ] `src/components/GroupCard.jsx` — title, truncated description, member count, "Request to Join" / "Member" button
- [ ] `src/app/groups/[id]/page.jsx` — title, description, member list, post feed, events section, invite button, pending requests (creator only)
- [ ] Events widget inside group detail page — list events with RSVP counts, Going/Not going buttons, "Create Event" form
- [ ] `src/lib/groups.js` — fetch wrappers for all group and event endpoints

---

## Issue #4 — Chat, Notifications 
**Assigned to:** Member 4

You are building the real-time layer of the app — private messaging, group chat, and notifications. The WebSocket hub must be done before the chat and notification features can be built.

### What you are building

**On the backend**, you will write the chats migration which creates a `messages` table for private DMs (`id`, `sender_id`, `receiver_id`, `content`, `created_at`) and a `group_messages` table for group chat (`id`, `group_id`, `sender_id`, `content`, `created_at`). Both tables must use `TEXT` for the `content` column with UTF-8 support so emoji characters store and retrieve correctly.

You will then build the WebSocket hub in `backend/pkg/websocket/hub.go`. The hub maintains a map of `user_id → connection` for every currently connected client. When a client connects, they send their `session_token` — the hub validates it and registers their connection under the correct `user_id`. The hub exposes two methods used by other parts of the backend: `SendToUser(userID, message)` for private messages and notifications, and `BroadcastToGroup(groupID, message)` for group chat. If the target user is not connected, the message is saved to the database so it can be fetched later. The hub is started once as a goroutine in `server.go`.

The chat handlers provide REST endpoints for fetching history: `GET /api/chat/users` lists everyone the logged-in user can DM (at least one must follow the other), `GET /api/chat/:userId` returns the full message history between two users, and `GET /api/groups/:id/messages` returns a group's chat history (members only).

The notifications handler stores and retrieves notifications. Each notification has a `type` (`follow_request`, `group_invite`, `group_join_request`, `group_event`), an `actor_id` (who triggered it), an `entity_id` (the relevant group/event/request ID), and an `is_read` flag. `GET /api/notifications` returns all notifications for the logged-in user. `PUT /api/notifications/:id/read` and `PUT /api/notifications/read-all` mark them as read. When a notification is created anywhere in the codebase, it is also pushed in real-time via `Hub.SendToUser` if the recipient is connected.

**On the frontend**, you will write the `useWebSocket` hook which manages a single WebSocket connection for the whole app, exposes a `sendMessage(payload)` function, accepts an `onMessage(handler)` callback, and automatically reconnects if the connection drops. The chat page (`/chat`) has a left sidebar listing DM conversations and group chats, and a `ChatWindow` on the right. `ChatWindow` loads history via REST on open, then receives new messages in real-time via the hook. Sending a message goes through the WebSocket directly, not REST. The `EmojiPicker` component opens a small panel above the text input and inserts the chosen emoji at the cursor. `NotificationBell` lives in the root layout and shows a badge with the unread count. Clicking it opens a dropdown where follow requests and group invites show Accept/Decline buttons inline.


### Checklist

#### Backend
- [ ] `000005_create_chats_table.up.sql` — create `messages` and `group_messages` tables (UTF-8 content)
- [ ] `000005_create_chats_table.down.sql` — drop both tables
- [ ] `backend/pkg/websocket/hub.go` — client registry, validate session on connect, `SendToUser`, `BroadcastToGroup`, save to DB if user offline
- [ ] Register and start the hub as a goroutine in `server.go`
- [ ] `GET /api/chat/users` — list DM-eligible users (mutual or one-way follow)
- [ ] `GET /api/chat/:userId` — return message history ordered by `created_at ASC` (403 if no follow relationship)
- [ ] `GET /api/groups/:id/messages` — return group message history (members only)
- [ ] Notifications table migration (`000007` or appended to existing) — `id`, `user_id`, `type`, `actor_id`, `entity_id`, `is_read`, `created_at`
- [ ] `backend/pkg/models/notification.go` — `CreateNotification`, `GetNotificationsForUser`, `MarkAsRead`, `MarkAllAsRead`
- [ ] `GET /api/notifications` — return notifications for logged-in user, unread first
- [ ] `PUT /api/notifications/:id/read` — mark one notification as read
- [ ] `PUT /api/notifications/read-all` — mark all as read

#### Frontend
- [ ] `src/hooks/useWebSocket.js` — open connection, `sendMessage()`, `onMessage()` callback, auto-reconnect with retry
- [ ] `src/app/chat/page.jsx` — DM list sidebar + group chats tab, renders `ChatWindow` for selected conversation
- [ ] `src/components/ChatWindow.jsx` — load history on open, append incoming messages, send via WebSocket, scroll to bottom on new message
- [ ] `src/components/EmojiPicker.jsx` — emoji panel, inserts at cursor position in text input
- [ ] `src/components/NotificationBell.jsx` — unread badge, dropdown list, inline Accept/Decline for actionable notifications, marks as read on open
- [ ] Add `NotificationBell` to `src/app/layout.jsx`

