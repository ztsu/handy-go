CREATE TABLE decks (
    "id"                   VARCHAR(36)  PRIMARY KEY,
    "name"                 TEXT         NOT NULL,
    "userId"               VARCHAR(36)  NOT NULL,

    "createdAt"            TIMESTAMP WITH TIME ZONE DEFAULT now()    NOT NULL,
    "modifiedAt"           TIMESTAMP WITH TIME ZONE DEFAULT now()    NOT NULL,

    CONSTRAINT "decks_userId_fkey" FOREIGN KEY ("userId") REFERENCES users ON DELETE CASCADE
);

CREATE TRIGGER "decks_updateModifiedAt"
    BEFORE UPDATE
    ON decks
    FOR EACH ROW EXECUTE PROCEDURE "updateModifiedAt"();