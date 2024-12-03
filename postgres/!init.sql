CREATE DATABASE payments;

CREATE DATABASE reservations;

CREATE DATABASE loyalties;


\c payments
CREATE TABLE payment
(
    id          SERIAL PRIMARY KEY,
    payment_uid uuid        NOT NULL,
    status      VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'CANCELED')),
    price       INT         NOT NULL
);

\c reservations
CREATE TABLE hotels
(
    id        SERIAL PRIMARY KEY,
    hotel_uid uuid         NOT NULL UNIQUE,
    name      VARCHAR(255) NOT NULL,
    country   VARCHAR(80)  NOT NULL,
    city      VARCHAR(80)  NOT NULL,
    address   VARCHAR(255) NOT NULL,
    stars     INT,
    price     INT          NOT NULL
);

CREATE TABLE reservation
(
    id              SERIAL PRIMARY KEY,
    reservation_uid uuid UNIQUE NOT NULL,
    username        VARCHAR(80) NOT NULL,
    payment_uid     uuid        NOT NULL,
    hotel_id        INT REFERENCES hotels (id),
    status          VARCHAR(20) NOT NULL
        CHECK (status IN ('PAID', 'CANCELED')),
    start_date      TIMESTAMP WITH TIME ZONE,
    end_date        TIMESTAMP WITH TIME ZONE
);

INSERT INTO public.hotels(hotel_uid, name, country, city, address, stars, price)
VALUES ('049161bb-badd-4fa8-9d90-87c9a82b0668'::uuid, 'Ararat Park Hyatt Moscow', 'Россия', 'Москва', 'Неглинная ул., 4', 5, 10000);

\c loyalties
CREATE TABLE loyalty
(
    id                SERIAL PRIMARY KEY,
    username          VARCHAR(80) NOT NULL UNIQUE,
    reservation_count INT         NOT NULL DEFAULT 0,
    status            VARCHAR(80) NOT NULL DEFAULT 'BRONZE'
        CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
    discount          INT         NOT NULL
);


INSERT INTO public.loyalty (username, reservation_count, status, discount)
VALUES ('Test Max', 25, 'GOLD', 10);