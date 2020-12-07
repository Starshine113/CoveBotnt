package cbdb

import (
	"context"
	"errors"
	"time"
)

// Trigger ...
type Trigger struct {
	ID       int
	GuildID  string
	Creator  string
	Modified time.Time
	Trigger  string
	Response string
}

// AddTrigger ...
func (db *Db) AddTrigger(t *Trigger) (*Trigger, error) {
	if t == nil {
		return nil, errors.New("trigger was nil")
	}
	if t.GuildID == "" || t.Creator == "" || t.Trigger == "" || t.Response == "" {
		return nil, errors.New("one or more required fields was nil")
	}
	if len(t.Trigger) > 99 {
		t.Trigger = t.Trigger[:99]
	}
	if len(t.Response) > 1999 {
		t.Response = t.Response[:1999]
	}
	var timestamp time.Time
	var id int

	err := db.Pool.QueryRow(context.Background(), "insert into public.triggers (guild_id, created_by, trigger, response) values ($1, $2, $3, $4) returning id, modified", t.GuildID, t.Creator, t.Trigger, t.Response).Scan(&id, &timestamp)
	t.ID = id
	t.Modified = timestamp
	return t, err
}

// Triggers gets all triggers for a guild
func (db *Db) Triggers(id string) (out []*Trigger, err error) {
	rows, err := db.Pool.Query(context.Background(), "select * from public.triggers where guild_id = $1", id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id                                  int
			guildID, creator, trigger, response string
			modified                            time.Time
		)

		rows.Scan(&id, &guildID, &creator, &modified, &trigger, &response)
		out = append(out, &Trigger{
			ID:       id,
			GuildID:  guildID,
			Creator:  creator,
			Modified: modified,
			Trigger:  trigger,
			Response: response,
		})
	}

	return
}

// RemoveTrigger ...
func (db *Db) RemoveTrigger(guildID string, triggerID int) (err error) {
	_, err = db.Pool.Exec(context.Background(), "delete from public.triggers where id = $1 and guild_id = $2", triggerID, guildID)
	return
}
