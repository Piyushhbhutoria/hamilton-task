CREATE TYPE "txntype" AS ENUM (
  'DEBIT',
  'CREDIT'
);

CREATE SEQUENCE "users_seq" START 1;

CREATE TABLE "users" (
  "id" bigint PRIMARY KEY DEFAULT (nextval('users_seq'::regclass)),
  "username" text,
  "created_at" timestamp DEFAULT (now()),
  "user_id" uuid UNIQUE DEFAULT (uuid_generate_v4())
);

CREATE SEQUENCE "wallet_seq" START 1;

CREATE TABLE "wallets" (
  "id" bigint PRIMARY KEY DEFAULT (nextval('wallet_seq'::regclass)),
  "user_id" uuid,
  "currency" text,
  "balance" numeric DEFAULT 0,
  "wallet_id" uuid UNIQUE DEFAULT (uuid_generate_v4())
);

CREATE SEQUENCE "ledger_seq" START 1;

CREATE TABLE "ledger" (
  "id" bigint PRIMARY KEY DEFAULT (nextval('ledger_seq'::regclass)),
  "created_at" timestamp DEFAULT (now()),
  "transaction_type" txntype,
  "amount" numeric DEFAULT 0,
  "description" text,
  "wallet_id" uuid
);

CREATE UNIQUE INDEX ON "users" ("user_id");

CREATE UNIQUE INDEX ON "wallets" ("wallet_id");

CREATE UNIQUE INDEX ON "wallets" ("user_id", "currency");

ALTER TABLE "wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "ledger" ADD FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("wallet_id");

CREATE VIEW wallet_view AS  SELECT w.id,
    w.user_id,
    w.currency,
    w.wallet_id,
    w.balance,
    to_json(ARRAY( SELECT json_build_object('created_at', l.created_at, 'transaction_type', l.transaction_type, 'amount', l.amount, 'description', l.description) AS json_build_object
           FROM ledger l
          WHERE w.wallet_id = l.wallet_id)) AS ledger
   FROM wallets w;

CREATE VIEW ledger_view AS SELECT l.id,
    l.created_at,
    l.transaction_type,
    l.amount,
    l.description,
    l.wallet_id,
    u.user_id,
    w.currency
   FROM ledger l
     JOIN wallets w ON w.wallet_id = l.wallet_id
     JOIN users u ON u.user_id = w.user_id;
