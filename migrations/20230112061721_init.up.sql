SET statement_timeout = 0;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--bun:split

CREATE TABLE user_account (
  id UUID NOT NULL PRIMARY KEY UNIQUE DEFAULT uuid_generate_v4(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  email_verified BOOLEAN NOT NULL DEFAULT false,
  image TEXT
);

CREATE UNIQUE INDEX unique_user_on_email
  ON user_account (email);


CREATE TABLE category (
  id UUID NOT NULL PRIMARY KEY UNIQUE DEFAULT uuid_generate_v4(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_id  UUID NOT NULL REFERENCES user_account (id),
  name TEXT NOT NULL,
  description TEXT
);

CREATE UNIQUE INDEX unique_category_on_user_id_and_name
  ON category (user_id, name);


CREATE TABLE bookmark (
  id UUID NOT NULL PRIMARY KEY UNIQUE DEFAULT uuid_generate_v4(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_id UUID NOT NULL REFERENCES user_account (id),
  category_id UUID NOT NULL REFERENCES category (id),
  url TEXT NOT NULL,
  image TEXT,
  description TEXT
);

CREATE UNIQUE INDEX unique_bookmark_on_user_id_and_url
  ON bookmark (user_id, url);


CREATE TABLE tag (
  id UUID NOT NULL PRIMARY KEY UNIQUE DEFAULT uuid_generate_v4(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  user_id UUID NOT NULL REFERENCES user_account (id),
  name TEXT NOT NULL,
  description TEXT,
  emoji TEXT
);

CREATE UNIQUE INDEX unique_tag_on_user_id_and_name
  ON tag (user_id, name);


CREATE TABLE tag_on_bookmark (
  bookmark_id  UUID NOT NULL REFERENCES bookmark (id),
  tag_id UUID NOT NULL REFERENCES tag (id),

  constraint pk_tag_on_bookmark
    PRIMARY KEY (bookmark_id, tag_id)
);


-- NOTE: We may move this table completely to redis!
CREATE TABLE session (
  token TEXT NOT NULL,
  user_id UUID NOT NULL REFERENCES user_account (id),
  expires TIMESTAMPTZ NOT NULL,
  raw_data TEXT NOT NULL,

  constraint pk_session
    PRIMARY KEY (token, user_id)
);
