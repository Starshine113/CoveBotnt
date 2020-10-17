package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var initDBSql = `create type modaction as enum ('warn', 'mute', 'unmute', 'pause', 'hardmute', 'kick', 'tempban', 'ban');

create table if not exists guild_settings
(
    guild_id			bigint primary key,
    starboard_channel	bigint not null default 0,
    react_limit			int not null default 100,
	emoji				text not null default '⭐',
    sender_can_react	boolean default false,
	react_to_starboard	boolean default true,
	
	mod_roles			bigint[] not null default array[0],
	helper_roles		bigint[] not null default array[0],
	mod_log				bigint not null default 0,
	mute_role			bigint not null default 0,
	pause_role			bigint not null default 0,

	gatekeeper_roles	bigint[] not null default array[0],
	member_roles		bigint[] not null default array[0],
	gatekeeper_channel	bigint not null default 0,
	gatekeeper_message	text not null default 'Please wait to be approved, {mention}.',
	welcome_channel		bigint not null default 0,
	welcome_message		text not null default 'Welcome to {guild}, {mention}!'
);

create table if not exists starboard_messages
(
    message_id				bigint primary key,
    channel_id				bigint,
    server_id				bigint not null references guild_settings (guild_id) on delete cascade,
    starboard_message_id	bigint
);

create table if not exists starboard_blacklisted_channels
(
    channel_id	bigint primary key,
    server_id	bigint not null references guild_settings (guild_id) on delete cascade
);

create table if not exists info
(
    id						int primary key not null default 1, -- enforced only equal to 1
    schema_version			int,
    constraint singleton	check (id = 1) -- enforce singleton table/row
);

insert into info (schema_version) values (1);`

func initDB() (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), config.Auth.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v", err)
	}
	if err := initDBIfNotInitialised(db); err != nil {
		fmt.Fprintf(os.Stderr, "[%v] There was an error while initialising the database: %v\n", time.Now().Format(time.RFC3339), err)
		os.Exit(1)
	}
	sugar.Infof("Connected to database.")
	return db, nil
}

func initDBIfNotInitialised(db *pgxpool.Pool) error {
	var exists bool
	err := db.QueryRow(context.Background(), "select exists (select from information_schema.tables where table_schema = 'public' and table_name = 'info')").Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil // the database has been initialised so we're done
	}

	// ...it's not initialised and we have to do that
	_, err = db.Exec(context.Background(), initDBSql)
	if err != nil {
		return err
	}
	sugar.Infof("Successfully initialised the database.")
	return nil
}
