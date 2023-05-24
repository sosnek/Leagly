package guilds

import (
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

var MY_DISCORD_ID = "220732095083839488"
var DB *bbolt.DB

type DiscordGuild struct {
	ID             string
	Region         string
	Region2        string
	Prefix         string
	JoinDate       string
	AutoPatchNotes bool
	PatchNotesCh   string
	Members        int
}

func SetupDB() (*bolt.DB, error) {
	db, err := bolt.Open("leagly.db", 0600, nil)
	if err != nil {
		return nil, errors.New("could not open db, " + err.Error())
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Guilds"))
		if err != nil {
			return errors.New("could not create Guilds bucket: " + err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("could not set up buckets, " + err.Error())
	}
	return db, nil
}

func Update(db *bolt.DB, key string, guild DiscordGuild) error {
	guildBytes, err := json.Marshal(guild)
	if err != nil {
		return errors.New("could not marshal guilds json: " + err.Error())
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("Guilds")).Put([]byte(key), guildBytes)
		if err != nil {
			return errors.New("could not update guilds: " + err.Error())
		}
		return nil
	})
	return err
}

func Add(db *bolt.DB, key string, guild DiscordGuild) error {
	err := db.Update(func(tx *bolt.Tx) error {
		guildBytes, err := json.Marshal(guild)
		err = tx.Bucket([]byte("Guilds")).Put([]byte(key), []byte(guildBytes))
		if err != nil {
			return errors.New("could not add guild" + err.Error())
		}
		return nil
	})
	return err
}

func View(db *bolt.DB, key string) (DiscordGuild, error) {
	var guild string
	var guildUnmarshal DiscordGuild
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Guilds"))
		guild = string(b.Get([]byte(key)))
		return nil
	})
	if err != nil {
		return guildUnmarshal, err
	}
	err = json.Unmarshal([]byte(guild), &guildUnmarshal)
	return guildUnmarshal, err
}

func Delete(db *bolt.DB, key string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Guilds"))
		b.Delete([]byte(key))
		return nil
	})
	return err
}

func ChannelsWithAutoPatchNotes() []string {
	var patchNoteChannels []string
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Guilds"))
		b.ForEach(func(k, v []byte) error {
			var tmp DiscordGuild
			err := json.Unmarshal(v, &tmp)
			if err != nil {
				return err
			}
			if tmp.AutoPatchNotes {
				patchNoteChannels = append(patchNoteChannels, tmp.PatchNotesCh)
			}
			return nil
		})
		return nil
	})
	return patchNoteChannels
}

func GetGuildCount() int {
	var count int
	DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Guilds"))
		b.ForEach(func(k, v []byte) error {
			count++
			return nil
		})
		return nil
	})
	return count
}

func HasDebugPermissions(authID string) bool {
	return authID == MY_DISCORD_ID
}
