CREATE OR REPLACE FUNCTION "updateModifiedAt"()
    RETURNS TRIGGER AS $$
BEGIN
    NEW."modifiedAt" = now();
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TABLE users (
    "id"                   VARCHAR(36) PRIMARY KEY,
    "email"                TEXT                                      NOT NULL,

    "createdAt"            TIMESTAMP WITH TIME ZONE DEFAULT now()    NOT NULL,
    "modifiedAt"           TIMESTAMP WITH TIME ZONE DEFAULT now()    NOT NULL,

    CONSTRAINT "users_email_ukey" UNIQUE ("email")
);

CREATE TRIGGER "users_updateModifiedAt"
    BEFORE UPDATE
    ON users
    FOR EACH ROW EXECUTE PROCEDURE "updateModifiedAt"();