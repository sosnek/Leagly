package query

import (
	"Leagly/guilds"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
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
			s.UpdateListeningStatus("/help")
			if newVersion != Version {
				Version = newVersion
				channelsWithAutoUpdates := guilds.ChannelsWithAutoPatchNotes()
				send, err := PatchNotes()
				if err != nil {
					log.Println("Error: Unable to get patchnotes " + err.Error())
					return
				}
				time.Sleep(5000) //I think it's trying to send the patchnotes before the image has been saved. Lets try waiting 5 seconds.
				for i := range channelsWithAutoUpdates {
					_, err := s.ChannelMessageSendComplex(channelsWithAutoUpdates[i], send)
					if err != nil {
						log.Println("Error while sending out patchnotes to " + channelsWithAutoUpdates[i] + ".  Error: " + err.Error())
					}
					log.Println("Sent patchnotes to channel ID: " + channelsWithAutoUpdates[i])
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
		if strings.Contains(imagesUrl[1], "1920x1080") || strings.Contains(imagesUrl[1], "Infographic") || strings.Contains(imagesUrl[1], "PatchNotes") {
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
			return err
		}
	}
	return nil
}

// /
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

// Because discord embeds only support 2x2 images at a maximum, I decided to use a method
// that combines 3 images into one to be use in a 1x3 format. Unfortunately discord also
// has limitations on image size. Embed will be constrained if the image is greater than 300px
// As a result, i limit 3 images combined to be at a maximum of 299 px in length :D
func mergeImages(imageName []string) string {

	var imgFile []*os.File
	var img []image.Image
	for n := 0; n < len(imageName); n++ {
		imgFile1, err := os.Open(imageName[n])
		if err != nil {
			fmt.Println(err)
		}
		imgFile = append(imgFile, imgFile1)
		img1, _, err := image.Decode(imgFile1)
		if err != nil {
			fmt.Println(err)
		}
		img = append(img, img1)
	}

	var sp3 image.Point
	var sp image.Point
	if len(img) == 3 {
		sp = image.Point{(img[0].Bounds().Dx() - 20), 0}
		sp3 = image.Point{sp.X + sp.X, 0}
	} else if len(img) == 2 {
		sp = image.Point{(img[0].Bounds().Dx()), 0}
		sp3 = image.Point{sp.X, 0}
	} else {
		sp = image.Point{(img[0].Bounds().Dx()), 0}
		sp3 = image.Point{0, 0}
	}
	r2 := image.Rectangle{sp, sp.Add(img[0].Bounds().Size())}

	r3 := image.Rectangle{sp3, sp3.Add(img[0].Bounds().Size())} //all images are same size anyways
	if len(img) == 3 {
		r3.Max.X = 299 //Discord embed width will be constrained if the image is 300px in width or greater
	} else if len(img) == 2 {
		r3.Max.X = sp.X + sp.X
	} else {
		r3.Max.X = sp.X
	}
	r := image.Rectangle{image.Point{0, 0}, r3.Max}

	rgba := image.NewRGBA(r)
	draw.Draw(rgba, img[0].Bounds(), img[0], image.Point{0, 0}, draw.Src)
	if len(img) == 2 {
		draw.Draw(rgba, r2, img[1], image.Point{0, 0}, draw.Src)
	} else if len(img) == 3 {
		draw.Draw(rgba, r2, img[1], image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r3, img[2], image.Point{0, 0}, draw.Src)
	}

	out, err := os.Create("./output.png")
	if err != nil {
		fmt.Println(err)
	}
	var opt jpeg.Options
	opt.Quality = 80

	jpeg.Encode(out, rgba, &opt)
	return "./output.png"
}
