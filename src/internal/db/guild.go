package db

import "context"

type DbGuild struct {
	GuildId string            `firestore:"guildId"`
	Name    string            `firestore:"name"`
	Users   map[string]string `firestore:"users"`
}

func (f *FirestoreDatabase) GetGuild(guildId string) (*DbGuild, error) {
	doc, err := f.client.Collection("guilds").Doc(guildId).Get(context.Background())
	if err != nil {
		return nil, err
	}
	if !doc.Exists() {
		return nil, nil
	}
	var guild DbGuild
	err = doc.DataTo(&guild)
	if err != nil {
		return nil, err
	}
	return &guild, nil
}
