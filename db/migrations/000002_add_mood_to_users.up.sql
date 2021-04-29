BEGIN;

ALTER TABLE users ADD COLUMN mood STRING;
ALTER TABLE users ADD CONSTRAINT check_mood CHECK (mood IN ('happy', 'sad', 'neutral'));

COMMIT;
