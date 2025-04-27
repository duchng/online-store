CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(255),
    role          VARCHAR(50) DEFAULT 'user',
    created_at    TIMESTAMP DEFAULT current_timestamp,
    updated_at    TIMESTAMP DEFAULT current_timestamp
);

select create_updated_at_trigger('users');

CREATE TABLE products
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    description    TEXT,
    price          DECIMAL(10, 2),
    stock_quantity INT,
    status         VARCHAR(50),
    created_at     TIMESTAMP DEFAULT current_timestamp,
    updated_at     TIMESTAMP DEFAULT current_timestamp
);

select create_updated_at_trigger('products');

CREATE TABLE categories
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_at     TIMESTAMP DEFAULT current_timestamp,
    updated_at     TIMESTAMP DEFAULT current_timestamp
);

select create_updated_at_trigger('categories');

CREATE TABLE product_categories
(
    product_id  INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (product_id, category_id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

create index product_categories_product_id_idx on product_categories (product_id);
create index product_categories_category_id_idx on product_categories (category_id);

CREATE TABLE reviews
(
    id         SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    user_id    INT NOT NULL,
    rating     INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    COMMENT    TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

select create_updated_at_trigger('reviews');

create index reviews_product_id_idx on reviews (product_id);
create index reviews_user_id_idx on reviews (user_id);

CREATE TABLE wishlist
(
    user_id    INT NOT NULL,
    product_id INT NOT NULL,
    added_at   TIMESTAMP DEFAULT current_timestamp,
    PRIMARY KEY (user_id, product_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

create index wishlist_user_id_idx on wishlist (user_id);
