CREATE SEQUENCE success_seq START 1;
CREATE SEQUENCE error_seq START 1;

CREATE TABLE relay_count (
    app_public_key char(64) NOT NULL,
    day date NOT NULL,
    success INTEGER DEFAULT nextval('success_seq') NOT NULL,
    error INTEGER DEFAULT nextval('error_seq') NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    PRIMARY KEY (app_public_key, day)
);
