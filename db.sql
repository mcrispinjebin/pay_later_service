CREATE DATABASE simple_service;
USE simple_service;

CREATE TABLE user
(
    user_id                int PRIMARY KEY AUTO_INCREMENT,
    user_name              varchar(30),
    user_email             varchar(30),
    credit_limit_offered   float,
    available_credit_limit float,
    created_at             BIGINT
);

CREATE TABLE merchant
(
    merchant_id      int PRIMARY KEY AUTO_INCREMENT,
    merchant_name    varchar(30),
    discount_percent float,
    created_at       BIGINT,
    updated_at       BIGINT
);

CREATE TABLE orders
(
    order_id     int PRIMARY KEY AUTO_INCREMENT,
    user_id      int REFERENCES user (user_id),
    merchant_id  int REFERENCES merchant (merchant_id),
    order_amount float,
    order_status enum('success', 'failed'),
    created_at   BIGINT
);

CREATE TABLE ledger
(
    ledger_id  int PRIMARY KEY AUTO_INCREMENT,
    user_id    int REFERENCES user (user_id),
    amount     float,
    status     enum('success', 'failed'),
    created_at BIGINT
);

CREATE TABLE payout
(
    payout_id     int PRIMARY KEY AUTO_INCREMENT,
    order_id      int REFERENCES orders (order_id),
    payout_amount float,
    payout_status enum('success', 'failed'),
    created_at    BIGINT
);


