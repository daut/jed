BEGIN TRANSACTION;

INSERT INTO admins (username, password) VALUES ('admin', crypt('password', gen_salt('bf')));

COMMIT;
