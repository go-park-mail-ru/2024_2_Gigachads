CREATE USER app_user WITH PASSWORD 'app_user_password';

GRANT SELECT, INSERT, UPDATE ON profile TO app_user;
GRANT SELECT, INSERT, UPDATE ON message TO app_user;
GRANT SELECT, INSERT, UPDATE ON folder TO app_user;
GRANT SELECT, INSERT, UPDATE ON email_transaction TO app_user;
GRANT SELECT, INSERT, UPDATE ON attachment TO app_user;
GRANT SELECT, INSERT, UPDATE ON session TO app_user;
GRANT SELECT, INSERT, UPDATE ON csrf_token TO app_user;
