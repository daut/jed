BEGIN TRANSACTION;

INSERT INTO products (name, price, description) VALUES ('Apple', 1.00, 'A delicious apple.');
INSERT INTO products (name, price, description) VALUES ('Banana', 0.50, 'A ripe banana.');
INSERT INTO products (name, price, description) VALUES ('Cherry', 0.10, 'A sweet cherry.');
INSERT INTO products (name, price, description) VALUES ('Date', 0.25, 'A sticky date.');

COMMIT;
