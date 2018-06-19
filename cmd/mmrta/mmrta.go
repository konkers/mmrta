package main

import (
	"log"
	"os"
	"strconv"

	"github.com/konkers/mmrta"
	"github.com/olekukonko/tablewriter"
)

func main() {
	c, err := mmrta.NewClient()
	if err != nil {
		log.Fatalf("Can't create new client: %v", err)
	}

	runs, err := c.GetUnverifiedRuns(true)
	if err != nil {
		log.Fatalf("Can't get runs: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Run Id", "Game", "Category", "User"})
	for _, r := range runs {
		table.Append([]string{
			strconv.FormatInt(int64(r.Id), 10),
			r.Game.Name,
			r.Category,
			r.User.Name,
		})
	}
	table.Render()

}
