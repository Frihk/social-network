"use client";

import { useEffect, useMemo, useState } from "react";
import {
  fetchNotifications,
  markAllNotificationsRead,
  markNotificationRead,
} from "@/lib/notifications";

const notificationCopy = {
  follow_request: "sent you a follow request",
  group_invite: "invited you to join a group",
  group_request: "requested to join your group",
  group_join_request: "requested to join your group",
  event_created: "created a new group event",
  group_event: "created a new group event",
  private_message: "sent you a message",
  group_message: "sent a group message",
};

function normalizeNotification(raw) {
  const data = raw?.data || raw?.notification || raw;
  return {
    id: data.id || `${data.type || "notification"}-${Date.now()}`,
    type: data.type || "notification",
    actor_id: data.actor_id || data.actorId || data.user_id || "",
    entity_id: data.entity_id || data.related_id || data.entityId || "",
    is_read: Number(data.is_read || data.read || 0),
    created_at: data.created_at || new Date().toISOString(),
  };
}

function notificationText(notification) {
  const action = notificationCopy[notification.type] || "sent you a notification";
  const actor = notification.actor_id ? `User ${notification.actor_id}` : "Someone";
  return `${actor} ${action}`;
}

export default function NotificationBell() {
  const [notifications, setNotifications] = useState([]);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const unreadCount = useMemo(
    () => notifications.filter((notification) => Number(notification.is_read) === 0).length,
    [notifications]
  );

  const loadNotifications = async () => {
    setLoading(true);
    setError("");
    try {
      const data = await fetchNotifications();
      setNotifications((data || []).map(normalizeNotification));
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadNotifications();

    const handleRealtimeNotification = (event) => {
      const notification = normalizeNotification(event.detail);
      setNotifications((prev) => {
        if (prev.some((item) => item.id === notification.id)) return prev;
        return [notification, ...prev];
      });
    };

    window.addEventListener("realtime-notification", handleRealtimeNotification);
    return () => window.removeEventListener("realtime-notification", handleRealtimeNotification);
  }, []);

  const handleMarkRead = async (notification) => {
    setNotifications((prev) =>
      prev.map((item) => item.id === notification.id ? { ...item, is_read: 1 } : item)
    );

    if (!notification.id.includes("-")) {
      try {
        await markNotificationRead(notification.id);
      } catch (err) {
        setError(err.message);
      }
    }
  };

  const handleMarkAllRead = async () => {
    setNotifications((prev) => prev.map((notification) => ({ ...notification, is_read: 1 })));
    try {
      await markAllNotificationsRead();
    } catch (err) {
      setError(err.message);
      loadNotifications();
    }
  };

  return (
    <div className="relative">
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        className="relative rounded-full border border-gray-200 bg-white px-3 py-2 text-sm font-semibold text-gray-700 hover:bg-gray-50"
        aria-label="Notifications"
      >
        Notifications
        {unreadCount > 0 && (
          <span className="absolute -right-2 -top-2 min-w-5 rounded-full bg-red-500 px-1.5 py-0.5 text-xs text-white">
            {unreadCount}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-20 mt-2 w-80 overflow-hidden rounded-lg border border-gray-200 bg-white shadow-lg">
          <div className="flex items-center justify-between border-b border-gray-100 px-4 py-3">
            <h2 className="text-sm font-bold text-gray-900">Notifications</h2>
            <button
              type="button"
              onClick={handleMarkAllRead}
              className="text-xs font-semibold text-blue-600 hover:text-blue-700"
            >
              Mark all read
            </button>
          </div>

          {error && <p className="px-4 py-3 text-sm text-red-600">{error}</p>}
          {loading && <p className="px-4 py-3 text-sm text-gray-500">Loading...</p>}

          {!loading && notifications.length === 0 && (
            <p className="px-4 py-6 text-center text-sm text-gray-500">No notifications yet.</p>
          )}

          <div className="max-h-96 overflow-y-auto">
            {notifications.map((notification) => (
              <button
                key={notification.id}
                type="button"
                onClick={() => handleMarkRead(notification)}
                className={`block w-full border-b border-gray-100 px-4 py-3 text-left hover:bg-gray-50 ${
                  Number(notification.is_read) === 0 ? "bg-blue-50" : "bg-white"
                }`}
              >
                <p className="text-sm font-medium text-gray-900">{notificationText(notification)}</p>
                {notification.entity_id && (
                  <p className="mt-1 text-xs text-gray-500">Reference: {notification.entity_id}</p>
                )}
                <p className="mt-1 text-xs text-gray-400">
                  {new Date(notification.created_at).toLocaleString()}
                </p>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
