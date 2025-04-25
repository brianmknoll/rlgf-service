package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d *DiscordClient) CreateDiscordEvent(name string) error {
	start, end := makeFakeTimes()
	params := &discordgo.GuildScheduledEventParams{
		Name:               name,
		Description:        "Test event",
		ScheduledStartTime: &start,
		ScheduledEndTime:   &end,
		// I think the only supported value is 2...
		// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-privacy-level
		// GUILD_ONLY | 2 | the scheduled event is only accessible to guild members
		PrivacyLevel: 2,
		Status:       1, // Scheduled
		// Type of the entity where event would be hosted
		// See field requirements
		// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-field-requirements-by-entity-type
		EntityType: 3, // External
	}
	_, err := d.sess.GuildScheduledEventCreate("test-guild-id", params)
	if err != nil {
		return err
	}
	return nil
}

// Creates a fake start and end time that is relative to now.
func makeFakeTimes() (start, end time.Time) {
	now := time.Now()
	return now.Add(2 * time.Hour).UTC(), now.Add(4 * time.Hour).UTC()
}
