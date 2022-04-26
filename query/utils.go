package query

import (
	"Leagly/config"
	"Leagly/guilds"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func UpdateVersionAsync(s *discordgo.Session) {
	uptimeTicker := time.NewTicker(15 * time.Minute)
	for {
		select {
		case <-uptimeTicker.C:
			newVersion := GetLeagueVersion()
			if newVersion != Version && newVersion != "12.7.1" { //Error happened if newVersion is 12.7.1
				Version = newVersion
				guildsWithAutoUpdates := guilds.GuildsWithAutoPatchNotes()
				send, err := PatchNotes()
				if err != nil {
					log.Println("Error: Unable to get patchnotes " + err.Error())
					return
				}
				for i := range guildsWithAutoUpdates {
					PatchNotesCh, err := Decrypt(guildsWithAutoUpdates[i], []byte(config.EncryptionKey))
					if err != nil {
						log.Println("Error: Unable to decrypt discord channel " + err.Error())
						return
					}
					s.ChannelMessageSendComplex(string(PatchNotesCh), send)
					log.Println("Sent patchnotes to: " + string(PatchNotesCh))
				}
			}
		}
	}
}

// A champion object is created on program start up containing all the names of all the champions and their ID's. Use this method to retrieve a name by ID
func GetChampion(champID string) string {
	for k, v := range champ3 {
		if champID == v.Key {
			return k // K is the champion name
		}
	}
	return champID
}

// An emoji object is created on program start up containing all the names and ID's of the emojis Leagly has access to.
func GetEmoji(emojiName string) string {
	for i := range emojis {
		for x := range emojis[i] {
			if emojis[i][x].Name == emojiName {
				return emojis[i][x].ID
			}
		}
	}
	return ""
}

func ParseVersion(version string) []string {
	return strings.Split(version, ".")
}

func patchNotesImgRegex(html []byte, version string) error {
	imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)

	subMatchSlice := imageRegExp.FindAllStringSubmatch(string(html), -1)
	for _, imagesUrl := range subMatchSlice {
		if strings.Contains(imagesUrl[1], "1920x1080") || strings.Contains(imagesUrl[1], "Infographic") || strings.Contains(imagesUrl[1], "LOL") {
			resp2, err := http.Get(imagesUrl[1])
			if err != nil {
				log.Println("Unable to get patchnotes image URL" + err.Error())
				return err
			}
			defer resp2.Body.Close()
			if resp2.StatusCode != 200 {
				log.Println("Unable to get URL with status code error:" + strconv.Itoa(resp2.StatusCode) + resp2.Status)
				return errors.New("Unable to get patchnotes image URL with status code error: " + resp2.Status)
			}
			file, err := os.Create("./patchNotes/" + version + ".png")
			if err != nil {
				log.Println("Error creating patchnotes file. Error: " + err.Error())
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, resp2.Body)
			if err != nil {
				log.Println("Error copying patchnotes file. Error: " + err.Error())
				return err
			}

		}
	}
	return nil
}

///
func InitializeEmojis(s *discordgo.Session) {
	var emoji [][]*discordgo.Emoji
	emoji1, _ := s.GuildEmojis("937465588446539920")
	emoji2, _ := s.GuildEmojis("937453232517693502")
	emoji3, _ := s.GuildEmojis("937481122198200320")
	emoji4, _ := s.GuildEmojis("937537071902503005")
	emoji5, _ := s.GuildEmojis("937482778499485756")
	emoji6, _ := s.GuildEmojis("938569984748163112")
	emoji7, _ := s.GuildEmojis("938569677326671913")
	emoji8, _ := s.GuildEmojis("938569400724910110")
	emoji9, _ := s.GuildEmojis("946539173597302804")
	emoji = append(emoji, emoji1, emoji2, emoji3, emoji4, emoji5, emoji6, emoji7, emoji8, emoji9)
	emojis = emoji
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
