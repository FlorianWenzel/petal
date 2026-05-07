create extension if not exists "pgcrypto";

create table users (
  id uuid primary key default gen_random_uuid(),
  email text not null unique,
  password_hash text not null,
  created_at timestamptz not null default now()
);

create table sessions (
  id text primary key,
  user_id uuid not null references users(id) on delete cascade,
  expires_at timestamptz not null
);

create index sessions_user_id_idx on sessions(user_id);

create table entries (
  user_id uuid not null references users(id) on delete cascade,
  day date not null,
  mood text not null,
  text text,
  updated_at timestamptz not null default now(),
  primary key (user_id, day)
);

alter table entries
  add constraint entries_mood_chk
  check (mood in ('happy','calm','love','sad','energized'));
