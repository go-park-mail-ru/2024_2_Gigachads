Relation User:
{user_id} -> email, username, password, avatar_url

Relation Email_message:
{email_id} -> email_parent_id, sender_id, recipient_id, title, description, date, isRead

Relation Csrf_token:
{csrf_id} -> user_id, hash, expire_date

Relation Session:
{session_id} -> user_id, hash, expire_date
