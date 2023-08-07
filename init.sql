CREATE TABLE IF NOT EXISTS "roles"(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR UNIQUE NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "users"(
    "id" BIGSERIAL PRIMARY KEY,
    "email" VARCHAR UNIQUE NOT NULL,
    "password" VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    "phone" VARCHAR NOT NULL,
    "role" BIGINT NOT NULL,
    "photo" BYTEA,
    "photo_format" VARCHAR,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_role FOREIGN KEY("role") REFERENCES "roles"("id")
);

CREATE TABLE IF NOT EXISTS "referral_statuses"(
    "id" BIGSERIAL PRIMARY KEY,
    "message" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS "user_details"(
    "user_id" BIGINT PRIMARY KEY,
    "referral_code" VARCHAR NOT NULL,
    "claimed_referral" VARCHAR,
    "total_transactions" NUMERIC NOT NULL DEFAULT 0,
    "gacha_quota" INTEGER NOT NULL DEFAULT 0,
    "balance" NUMERIC NOT NULL DEFAULT 0,
    "referral_status" BIGINT,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_user FOREIGN KEY("user_id") REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "addresses"(
    "id" BIGSERIAL PRIMARY KEY,
    "recipient_name" VARCHAR NOT NULL,
    "full_address" VARCHAR NOT NULL,
    "recipient_phone" VARCHAR NOT NULL,
    "user_id" BIGINT NOT NULL,
    "label" VARCHAR NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_role FOREIGN KEY("user_id") REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "shipping_statuses"(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "categories"(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "description" VARCHAR,
    "price" NUMERIC NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "sizes"(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "description" VARCHAR,
    "price" NUMERIC NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "promos"(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "min_fee" NUMERIC NOT NULL,
    "discount" NUMERIC NOT NULL,
    "max_discount" NUMERIC NOT NULL,
    "quota" INTEGER NOT NULL,
    "limited" BOOLEAN NOT NULL,
    "exp_date" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "user_vouchers"(
    "id" BIGSERIAL PRIMARY KEY,
    "promo_id" BIGINT NOT NULL,
    "user_id" BIGINT NOT NULL,
    "exp_date" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_promo FOREIGN KEY("promo_id") REFERENCES "promos"("id"),
    CONSTRAINT fk_user FOREIGN KEY("user_id") REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "payments"(
    "id" BIGSERIAL PRIMARY KEY,
    "status" BOOLEAN NOT NULL,
    "total_cost" NUMERIC NOT NULL,
    "total_discount" NUMERIC NOT NULL,
    "voucher_id" BIGINT,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_voucher FOREIGN KEY("voucher_id") REFERENCES "user_vouchers"("id")
);



CREATE TABLE IF NOT EXISTS "transactions"(
    "id" BIGSERIAL PRIMARY KEY,
    "description" VARCHAR,
    "amount" NUMERIC NOT NULL,
    "user_id" BIGINT NOT NULL,
    "payment_id" BIGINT,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_user FOREIGN KEY("user_id") REFERENCES "users"("id"),
    CONSTRAINT fk_payment FOREIGN KEY("payment_id") REFERENCES "payments"("id")
);

CREATE TABLE IF NOT EXISTS "shippings"(
    "id" BIGSERIAL PRIMARY KEY,
    "size_id" BIGINT NOT NULL,
    "category_id" BIGINT NOT NULL,
    "address_id" BIGINT NOT NULL,
    "payment_id" BIGINT NOT NULL,
    "status_id" BIGINT NOT NULL,
    "review_comment" VARCHAR,
    "review_rating" INTEGER,
    "user_id" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_size FOREIGN KEY("size_id") REFERENCES "sizes"("id"),
    CONSTRAINT fk_category FOREIGN KEY("category_id") REFERENCES "categories"("id"),
    CONSTRAINT fk_address FOREIGN KEY("address_id") REFERENCES "addresses"("id"),
    CONSTRAINT fk_payment FOREIGN KEY("payment_id") REFERENCES "payments"("id"),
    CONSTRAINT fk_status FOREIGN KEY("status_id") REFERENCES "shipping_statuses"("id"),
    CONSTRAINT fk_user FOREIGN KEY("user_id") REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "add_ons"(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "description" VARCHAR,
    "price" NUMERIC NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "deleted_at" TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS "add_on_shippings"(
    "shipping_id" BIGINT NOT NULL,
    "add_on_id" BIGINT NOT NULL,
    PRIMARY KEY("shipping_id", "add_on_id"),
    CONSTRAINT fk_shipping FOREIGN KEY("shipping_id") REFERENCES "shippings"("id"),
    CONSTRAINT fk_add_ons FOREIGN KEY("add_on_id") REFERENCES "add_ons"("id")
);

INSERT INTO "roles"("name", "created_at", "updated_at") VALUES
                   ('user', now(), now()),
                   ('admin', now(), now());

INSERT INTO "users" ("email", "password", "name", "phone", "role","created_at", "updated_at") VALUES 
                    ('budi@gmail.com','$2a$04$VJAf6MggcmvZAfxY1RDK8e.IiRamUhzYeGoDocpN5IUhBV/dvH9QS', 'Budi Sutanto', '081318235566', 1, now(), now()),
                    ('tata@copi.com','$2a$04$G7qUIOnZH2wgR3cLDOQSpe0nwUxoLmnHkqpogoGinJgWSyuL8RwuW', 'Tata Tatata', '082112340987', 2, now(), now()),
                    ('admin@gmail.com','$2a$04$bOI9QPQg9DqB9iBfbOkuoeDs9JklaqfTAisQvq8zY.Ox2b831D9qW', 'Admin Kurir', '082222340987', 2, now(), now()),
                    ('user@gmail.com','$2a$10$vK.JR/dYaUvvST4xhnS7A.rNCIP2sFc2tUDx9FialwiPuOkm0buiG', 'User', '082222340987', 1, now(), now());

INSERT INTO "user_details"("user_id", "referral_code", "claimed_referral","referral_status" "created_at", "updated_at") VALUES
                    (1, 'QiGPTHrjZJBmZnmy', null, 1, NOW(), NOW()),
                    (4, 'nxjwANKpGnDdBcCC', null, 1, NOW(), NOW());
                    
INSERT INTO "add_ons" ("name", "description", "price", "created_at", "updated_at") VALUES
                    ('Safe Package', 'Adds more protection to your package', 50000, NOW(), NOW()),
                    ('Cooler', 'Keeps your items cool during delivery', 20000, NOW(), NOW()),
                    ('Heatkeeper', 'Keeps your items warm during delivery', 15000, NOW(), NOW()),
                    ('Insurance', 'Insure your items for a flat fee for claims up to Rp 100.000.000', 20000, NOW(), NOW());

INSERT INTO "sizes" ("name", "description", "price", "created_at", "updated_at") VALUES
                    ('Large', 'VOLUME > 50x50x50cm OR WEIGHT > 3kg ', 150000, NOW(), NOW()),
                    ('Medium', '25x25x25cm < VOLUME <= 50x50x50cm OR 2kg < WEIGHT <= 3kg', 100000, NOW(), NOW()),
                    ('Small', 'VOLUME <= 25x25x25cm OR WEIGHT <= 2kg', 75000, NOW(), NOW());

INSERT INTO "shipping_statuses" ("name", "created_at", "updated_at") VALUES ('Waiting for payment', NOW(), NOW()), ('Processing', NOW(), NOW()), ('In Transit', NOW(), NOW()), ('Delivering to Destination', NOW(), NOW()), ('Arrived', NOW(), NOW());
INSERT INTO "categories" ("name", "description", "price", "created_at", "updated_at") VALUES
                    ('Food and Beverages', 'Send your food or beverages safely with us', 30000, NOW(), NOW()),
                    ('Fragile', 'Send your fragile package safely with us', 25000, NOW(), NOW()),
                    ('Documents', 'Send your documents safely with us', 20000, NOW(), NOW()),
                    ('Electronics', 'Send your electronic items safely with us', 30000, NOW(), NOW()),
                    ('Medical', 'Send your medical goods safely and hygienically with us', 30000, NOW(), NOW());

INSERT INTO "addresses" ("recipient_name", "full_address", "recipient_phone", "user_id", "label", "created_at") VALUES
                        ('Raven Treves', '896 Starling Junction', '421-883-1042', 4, 'Tú', NOW()),
                        ('Joelie Pescod', '5991 Petterle Plaza', '588-210-4143', 4, 'Maëly', NOW()),
                        ('Clair Coarser', '91662 Algoma Parkway', '718-724-7424', 4, 'Ráo', NOW()),
                        ('Fredric Fery', '32 Miller Crossing', '950-139-2966', 4, 'Géraldine', NOW()),
                        ('Bekki Whaites', '324 Doe Crossing Pass', '872-153-4373', 4, 'Clélia', NOW()),
                        ('Edin Womack', '1 Scofield Hill', '190-663-1826', 4, 'Gösta', NOW()),
                        ('Giselle Stickford', '11 Menomonie Crossing', '708-890-0100', 4, 'Mahélie', NOW()),
                        ('Christel Strover', '6 Dawn Circle', '668-974-9934', 4, 'Aloïs', NOW()),
                        ('Mitchell Prozescky', '4 Hoepker Park', '760-485-4543', 4, 'Börje', NOW()),
                        ('Justino McKee', '47 Garrison Alley', '197-361-0500', 4, 'Joséphine', NOW()),
                        ('Jerrilyn Hansom', '8198 Almo Alley', '256-894-7287', 4, 'Sòng', NOW()),
                        ('Merle Marshland', '14168 Rusk Place', '692-502-5152', 4, 'Mahélie', NOW()),
                        ('Lorna Geldeard', '12 Bartillon Circle', '440-259-4261', 4, 'Océanne', NOW()),
                        ('Ikey Cullen', '90552 Moland Crossing', '725-274-2677', 4, 'Lóng', NOW()),
                        ('Chancey Addlestone', '483 Mallory Street', '508-853-7175', 4, 'Dorothée', NOW()),
                        ('Flossy Shafto', '11692 Morningstar Circle', '873-245-9655', 4, 'Bérangère', NOW()),
                        ('Miof mela Muslim', '2 Farwell Terrace', '140-936-5032', 4, 'Cunégonde', NOW()),
                        ('Esta Middlebrook', '507 Scofield Pass', '530-185-4395', 4, 'Mårten', NOW()),
                        ('Hort Neylon', '200 Arizona Parkway', '249-295-8972', 4, 'Miléna', NOW()),
                        ('Kameko Gentle', '125 Grim Terrace', '776-350-3856', 4, 'Anaïs', NOW());
INSERT INTO promos ("name", "min_fee", "discount", "max_discount", "quota", "limited","exp_date",  "created_at", "updated_at") VALUES
                   ('Always Cheaper 1', 20000, 0.4, 20000, 1000, false, NOW()+INTERVAL '1 WEEK',  NOW(), NOW()), 
                   ('Always Cheaper 2', 20000, 0.6, 20000, 1000, false, NOW()+INTERVAL '1 WEEK', NOW(), NOW()),
                   ('Always Cheaper 2', 20000, 0.8, 20000, 1000, false, NOW()+INTERVAL '1 WEEK', NOW(), NOW()), 
                   ('Always Discount', 20000, 0.05, 20000, 1000, false, NOW()+INTERVAL '365 YEAR', NOW(), NOW());

INSERT INTO "referral_statuses" ("message") VALUES 
                        ('NONE'), ('UNCLAIMED'), ('CLAIMED_USER'), ('CLAIMED_FULL');

