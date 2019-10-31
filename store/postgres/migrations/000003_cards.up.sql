CREATE TABLE cards (
    "id"           VARCHAR(36)  PRIMARY KEY,
    "userId"       VARCHAR(36)        NOT NULL,
    "from"         TEXT         NOT NULL,
    "to"           TEXT         NOT NULL,
    "word"         TEXT         NOT NULL,
    "translation"  TEXT         NOT NULL,
    "ipa"                       TEXT,

    "createdAt"            TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "modifiedAt"           TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,

    CONSTRAINT "cards_userId_fkey" FOREIGN KEY ("userId") REFERENCES users ON DELETE CASCADE
);

CREATE TRIGGER "cards_updateModifiedAt"
    BEFORE UPDATE
    ON cards
    FOR EACH ROW EXECUTE PROCEDURE "updateModifiedAt"();


CREATE TABLE deck_cards (
    "id"                   VARCHAR(36)        PRIMARY KEY,
    "deckId"               VARCHAR(36)        NOT NULL,
    "cardId"               VARCHAR(36)        NOT NULL,
    "views"                INTEGER DEFAULT 0  NOT NULL,
    "turns"                INTEGER DEFAULT 0  NOT NULL,

    "createdAt"            TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "modifiedAt"           TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,

    CONSTRAINT "deck_cards_deckId_fkey" FOREIGN KEY ("deckId") REFERENCES decks ON DELETE CASCADE,
    CONSTRAINT "deck_cards_cardId_fkey" FOREIGN KEY ("cardId") REFERENCES cards ON DELETE CASCADE
);

CREATE TRIGGER "deck_cards_updateModifiedAt"
    BEFORE UPDATE
    ON deck_cards
    FOR EACH ROW EXECUTE PROCEDURE "updateModifiedAt"();