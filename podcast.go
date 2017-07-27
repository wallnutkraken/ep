package ep

import (
	"encoding/json"
	"errors"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"os"
	"strings"
)

type Podcast struct {
	Name          string    `json:"name"`
	URL           string    `json:"feed_url"`
	Tag           string    `json:"tag"`
	EpisodicItems []Episode `json:"episodes"`
}

func (p *Podcast) Write() error {
	data, _ := json.Marshal(p)
	file, err := os.Create(feedsDir + dir_seperator + p.Tag)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return file.Close()
}

func AddPodcast(tag string, url string) error {
	if len(tag) == 0 {
		errors.New("Tag is required")
	}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	podcast := Podcast{Name: feed.Title, URL: url, Tag: tag}
	err = podcast.UpdateEpisodes()
	if err != nil {
		return errors.New("failed updating episodes: " + err.Error())
	}
	if exists(feedsDir + dir_seperator + tag) {
		return errors.New("A podcast with the tag \"" + tag + "\" already exists")
	}
	return podcast.Write()
}

func ListAll() ([]Podcast, error) {
	dir, err := ioutil.ReadDir(feedsDir)
	if err != nil {
		return nil, err
	}
	podcasts := make([]Podcast, 0)

	for _, dirEntry := range dir {
		podcast := Podcast{}
		data, err := ioutil.ReadFile(feedsDir + dir_seperator + dirEntry.Name())
		if err != nil {
			/* Ignore */
			continue
		}
		err = json.Unmarshal(data, &podcast)
		if err != nil {
			/* Ignore */
			continue
		}
		podcasts = append(podcasts, podcast)
	}

	return podcasts, nil
}

func GetPodcast(tag string) (Podcast, error) {
	data, err := ioutil.ReadFile(feedsDir + dir_seperator + tag)
	if err != nil {
		return Podcast{}, err
	}
	p := Podcast{}
	err = json.Unmarshal(data, &p)
	return p, err
}

func (p *Podcast) UpdateEpisodes() error {
	feed, err := gofeed.NewParser().ParseURL(p.URL)
	if err != nil {
		return err
	}
	episodes := make([]Episode, 0)

	/* Reverse loop, to start from earliest items */
	for index := len(feed.Items) - 1; index > 0; index-- {
		item := feed.Items[index]
		ep := Episode{}
		ep.Title = item.Title

		/* Find the actual URL to the audio */
		for _, encl := range item.Enclosures {
			if strings.HasPrefix(encl.Type, "audio") {
				ep.URL = encl.URL
				/* We got the one we want; break */
				break
			}
		}
		episodes = append(episodes, ep)
	}

	p.EpisodicItems = episodes
	return nil
}

// Remove physically removes the podcast from the drive
func (p *Podcast) Remove() error {
	return os.Remove(feedsDir + dir_seperator + p.Tag)
}