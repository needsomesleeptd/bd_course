CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    page_count INT NOT NULL,
    document_name VARCHAR(255) NOT NULL,
    checks_count INT NOT NULL,
    creator_id BIGINT NOT NULL,
    creation_time TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS markup_types (
    id serial  PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    creator_id INT NOT NULL,
    class_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
   id serial PRIMARY KEY,
   login VARCHAR(255) UNIQUE NOT NULL,
   password VARCHAR(255) NOT NULL,
   name VARCHAR(255),
   surname VARCHAR(255) NOT NULL,
   role INT NOT NULL,
   user_group VARCHAR(255) NOT NULL 
);

CREATE TABLE IF NOT EXISTS markups (
    id serial PRIMARY KEY NOT NULL,
    page_data BYTEA NOT NULL,
    error_bb JSONB DEFAULT '[]' NOT NULL,
    class_label BIGINT NOT NULL,
    creator_id BIGINT NOT NULL
);


CREATE TABLE IF NOT EXISTS document_queues (
    id serial PRIMARY KEY NOT NULL,
    doc_id UUID NOT NULL,
    status SMALLINT NOT NULL,
    CONSTRAINT fk_doc_queue_docs FOREIGN KEY (doc_id) REFERENCES documents(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS comments (
    id serial PRIMARY KEY NOT NULL,
    doc_id UUID NOT NULL,
    description text NOT NULL,
    creator_id INT NOT NULL,
    CONSTRAINT fk_comment_docs FOREIGN KEY (doc_id) REFERENCES documents(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_users FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS achievments (
    id serial PRIMARY KEY NOT NULL,
    description text NOT NULL,
    creator_id INT NOT NULL,
    granted_to_id INT NOT NULL,
    CONSTRAINT fk_comment_creator FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_granted FOREIGN KEY (granted_to_id) REFERENCES users(id) ON DELETE CASCADE
);


ALTER TABLE markups ADD CONSTRAINT fk_markup_markup_type FOREIGN KEY ( class_label ) REFERENCES markup_types( id );

ALTER TABLE documents ADD CONSTRAINT fk_document_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE markups ADD CONSTRAINT fk_markup_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE markup_types ADD CONSTRAINT fk_markup_types_user FOREIGN KEY ( creator_id ) REFERENCES users( id );

ALTER TABLE document_queues ADD CONSTRAINT fk_doc_queue_docs FOREIGN KEY ( doc_id ) REFERENCES documents( id );


