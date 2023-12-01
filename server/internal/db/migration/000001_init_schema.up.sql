create table
  public.games (
    id uuid not null default uuid_generate_v4 (),
    fen text null,
    history text not null default ''::text,
    completed boolean not null default false,
    date_started timestamp with time zone not null default now(),
    date_finished timestamp with time zone null,
    current_state text not null,
    ruleset text not null default 'default'::text,
    type text not null default ''::text,
    user_1 uuid not null,
    user_2 uuid not null,
    result text null,
    constraint games_pkey primary key (id),
    constraint games_id_key unique (id),
    constraint games_user_1_fkey foreign key (user_1) references auth.users (id),
    constraint games_user_2_fkey foreign key (user_2) references auth.users (id)
  ) tablespace pg_default;

create table
  public.undo_request (
    id bigint generated by default as identity,
    game_id uuid not null,
    sender_id uuid not null,
    receiver_id uuid not null,
    status text not null default 'pending'::text,
    constraint undo_request_pkey primary key (id),
    constraint undo_request_game_id_sender_id_unique unique (game_id, sender_id),
    constraint undo_request_game_id_fkey foreign key (game_id) references games (id) on update cascade on delete cascade,
    constraint undo_request_receiver_id_fkey foreign key (receiver_id) references auth.users (id) on update cascade on delete cascade,
    constraint undo_request_sender_id_fkey foreign key (sender_id) references auth.users (id) on update cascade on delete cascade
  ) tablespace pg_default;

  create table
  public.room_list (
    id uuid not null default gen_random_uuid (),
    host_id uuid not null,
    description text not null,
    rules text not null,
    type text not null,
    color text not null,
    created_at timestamp with time zone not null default now(),
    constraint room_list_pkey primary key (id),
    constraint room_list_host_fkey foreign key (host) references auth.users (id) on update cascade on delete cascade
  ) tablespace pg_default;

  create table
  public.profiles (
    id uuid not null,
    username text not null,
    is_username_onboard_complete boolean not null default false,
    constraint profiles_pkey primary key (id),
    constraint profiles_username_key unique (username),
    constraint profiles_id_fkey foreign key (id) references auth.users (id) on delete cascade
  ) tablespace pg_default;

CREATE SCHEMA auth;

  create table
  auth.users (
    instance_id uuid null,
    id uuid not null,
    aud character varying(255) null,
    role character varying(255) null,
    email character varying(255) null,
    encrypted_password character varying(255) null,
    email_confirmed_at timestamp with time zone null,
    invited_at timestamp with time zone null,
    confirmation_token character varying(255) null,
    confirmation_sent_at timestamp with time zone null,
    recovery_token character varying(255) null,
    recovery_sent_at timestamp with time zone null,
    email_change_token_new character varying(255) null,
    email_change character varying(255) null,
    email_change_sent_at timestamp with time zone null,
    last_sign_in_at timestamp with time zone null,
    raw_app_meta_data jsonb null,
    raw_user_meta_data jsonb null,
    is_super_admin boolean null,
    created_at timestamp with time zone null,
    updated_at timestamp with time zone null,
    phone text null default null::character varying,
    phone_confirmed_at timestamp with time zone null,
    phone_change text null default ''::character varying,
    phone_change_token character varying(255) null default ''::character varying,
    phone_change_sent_at timestamp with time zone null,
    confirmed_at timestamp with time zone null,
    email_change_token_current character varying(255) null default ''::character varying,
    email_change_confirm_status smallint null default 0,
    banned_until timestamp with time zone null,
    reauthentication_token character varying(255) null default ''::character varying,
    reauthentication_sent_at timestamp with time zone null,
    is_sso_user boolean not null default false,
    deleted_at timestamp with time zone null,
    constraint users_pkey primary key (id),
    constraint users_phone_key unique (phone),
    constraint users_email_change_confirm_status_check check (
      (
        (email_change_confirm_status >= 0)
        and (email_change_confirm_status <= 2)
      )
    )
  );

create index if not exists users_instance_id_idx on auth.users using btree (instance_id);

create index if not exists users_instance_id_email_idx on auth.users using btree (instance_id, lower((email)::text));

create unique index confirmation_token_idx on auth.users using btree (confirmation_token)
where
  ((confirmation_token)::text !~ '^[0-9 ]*$'::text);

create unique index recovery_token_idx on auth.users using btree (recovery_token)
where
  ((recovery_token)::text !~ '^[0-9 ]*$'::text);

create unique index email_change_token_current_idx on auth.users using btree (email_change_token_current)
where
  (
    (email_change_token_current)::text !~ '^[0-9 ]*$'::text
  );

create unique index email_change_token_new_idx on auth.users using btree (email_change_token_new)
where
  (
    (email_change_token_new)::text !~ '^[0-9 ]*$'::text
  );

create unique index reauthentication_token_idx on auth.users using btree (reauthentication_token)
where
  (
    (reauthentication_token)::text !~ '^[0-9 ]*$'::text
  );

create unique index users_email_partial_key on auth.users using btree (email)
where
  (is_sso_user = false);
