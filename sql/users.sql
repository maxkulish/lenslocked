CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name text,
    email text NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount INT,
    description text
);

SELECT users.id, users.name, users.email, orders.id as order_id, orders.amount, orders.description FROM users
		INNER JOIN orders ON users.id=orders.user_id;