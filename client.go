package mmrta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const abbrevPrefix = "Mega Man "

type Client struct {
}

type Series struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Game struct {
	Id           int        `json:"id"`
	SeriesId     int        `json:"series_id"`
	Name         string     `json:"name"`
	ShortName    string     `json:"short_name"`
	HasGametime  int        `json:"has_gametime"`
	UsesGametime int        `json:"uses_gametime"`
	Forum        string     `json:"forum"`
	Categories   []Category `json:"categories"`
}

func (g *Game) AbbrevName() string {
	return strings.TrimPrefix(g.Name, abbrevPrefix)
}

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	UserLevel int    `json:"user_level"`
	Country   string `json:"country"`
	Twitter   string `json:"twitter"`
	Twitch    string `json:"twitch"`
	Hitbox    string `json:"hitbot"`
	YouTube   string `json:"youtube"`
}

type Run struct {
	Id            int     `json:"id"`
	Runner        *string `json:"runner"`
	UserId        *int    `json:"user_id"`
	GameId        *int    `json:"game_id"`
	Game          *Game   `json:"game"`
	User          *User   `json:"user"`
	Category      string  `json:"category"`
	Version       string  `json:"version"`
	VersionDetail *string `json:"version_detail"`
	Time          int     `json:"time"`
	GameTime      *int    `json:"game_time"`
	ConvertedTime int     `json:"converted_time"`
	Video         string  `json:"video"`
	Verified      int     `json:"verified"`
	VerifiedBy    string  `json:"verified_by"`
	Notes         string  `json:"notes"`
	SubmittedBy   string  `json:"submitted_by"`
	Rank          int     `json:"rank"`
}

func (r *Run) AbbrevCat() string {
	return strings.TrimPrefix(r.Category, abbrevPrefix)
}

func (r *Run) Duration() time.Duration {
	return time.Duration(r.Time*10) * time.Millisecond
}

func (r *Run) PrettyTime() string {
	d := r.Duration()
	hours := d.Truncate(time.Hour)
	d -= hours
	mins := d.Truncate(time.Minute)
	d -= mins

	if hours == 0.0 {
		return fmt.Sprintf("%02d:%02d", int(mins.Minutes()), int(d.Seconds()))
	} else {
		return fmt.Sprintf("%d:%02d:%02d", int(hours.Hours()),
			int(mins.Minutes()), int(d.Seconds()))

	}
}

type response struct {
	Series []*Series `json:"series"`
	Games  []*Game   `json:"games"`
	Runs   []*Run    `json:"runs"`
}

const urlBase = "https://megamanleaderboards.net/api/"

func NewClient() (*Client, error) {
	return &Client{}, nil
}

func (c *Client) getRequest(endpoint string, args map[string]string) ([]byte, error) {
	u, err := url.Parse(urlBase)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, endpoint+".php")
	q := u.Query()
	for param, value := range args {
		q.Set(param, value)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func decodeResponse(data []byte) (*response, error) {
	filteredString := regexp.MustCompile(`\\'`).ReplaceAllString(string(data), "'")
	var val response
	err := json.Unmarshal([]byte(filteredString), &val)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

func (c *Client) getJsonRequest(endpoint string, args map[string]string) (*response, error) {
	data, err := c.getRequest(endpoint, args)
	if err != nil {
		return nil, err
	}

	return decodeResponse(data)
}

func (c *Client) GetSeries() ([]*Series, error) {
	resp, err := c.getJsonRequest("series", nil)
	if err != nil {
		return nil, err
	}

	return resp.Series, err
}

func (c *Client) GetGames() ([]*Game, error) {
	resp, err := c.getJsonRequest("games", nil)
	if err != nil {
		return nil, err
	}

	return resp.Games, err
}

func (c *Client) GetGameById(id int) ([]*Game, error) {
	resp, err := c.getJsonRequest("games",
		map[string]string{"game": strconv.FormatInt(int64(id), 10)})
	if err != nil {
		return nil, err
	}

	if len(resp.Games) != 1 {
		return nil, fmt.Errorf("Got %d games for id %d.  Expected 0.", len(resp.Games), id)
	}
	return resp.Games, err
}

func (c *Client) GetGamesBySeries(series int) ([]*Game, error) {
	resp, err := c.getJsonRequest("games",
		map[string]string{"series": strconv.FormatInt(int64(series), 10)})
	if err != nil {
		return nil, err
	}

	return resp.Games, err
}

func (c *Client) GetUnverifiedRuns(expanded bool) ([]*Run, error) {
	var args map[string]string
	if expanded {
		args = map[string]string{"ex": "1"}
	}
	resp, err := c.getJsonRequest("runs", args)
	if err != nil {
		return nil, err
	}

	return resp.Runs, err
}
