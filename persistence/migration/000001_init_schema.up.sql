CREATE TABLE "users"
(
    id         uuid PRIMARY KEY,
    name       varchar        NOT NULL,
    username   varchar unique not null,
    email      varchar unique not null,
    password   varchar        not null,
    balance    bigint                  default 0,
    currency   varchar  not null,
    created_at timestamp      NOT NULL DEFAULT NOW(),
    updated_at timestamp
);

CREATE TABLE "banks"
(
    "id"                  uuid PRIMARY KEY,
    "name"                varchar       not null,
    "user_id"             uuid          not null,
    "branch_number"       varchar       not null,
    "account_number"      varchar       NOT NULL,
    "account_holder_name" varchar       NOT NULL,
    "reference"           varchar,
    "currency"            varchar not null,
    "expire_at"           timestamp     NOT NULL,
    "created_at"          timestamptz   NOT NULL DEFAULT (now()),
    updated_at            timestamp
);

/*ENUM ('FROM_BANK_TO_ACCOUNT_DEPOSIT','FROM_ACCOUNT_TO_ACCOUNT')*/
CREATE TABLE "transfers"
(
    "id"         uuid PRIMARY KEY,
    "type"    varchar       NOT NULL,
    "from_id"    uuid        not null,
    "to_id"      uuid        not null,
    "amount"     bigint      NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "banks"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
CREATE
INDEX ON "transfers" ("type","from_id");
CREATE UNIQUE  INDEX  ON "banks" ("user_id","name","account_holder_name","account_number");
CREATE
INDEX ON "transfers" ("type","to_id");
CREATE
INDEX ON "users" ("email");
CREATE
INDEX ON "users" ("username");
CREATE
INDEX ON "banks" ("user_id");