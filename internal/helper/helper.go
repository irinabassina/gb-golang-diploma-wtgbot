package helper

import (
	"WarehouseTgBot/internal/database"
	"bytes"
	"github.com/olekukonko/tablewriter"
	"strconv"
	"time"
)

func ConvertToHTML(cats []database.GoodCategory) string {
	buf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"ID", "NAME", "DESC", "UNIT", "COST", "CREATED_BY", "CREATED", "UPDATED"})
	for _, c := range cats {
		table.Append([]string{strconv.FormatInt(c.ID, 10), c.Name, c.Description, c.Unit,
			strconv.FormatFloat(c.Cost, 'f', 2, 64), strconv.FormatInt(c.CreatedBy, 10),
			c.CreatedAt.Format(time.Layout), c.UpdatedAt.Format(time.Layout)})
	}
	table.Render()
}

func ConvertToHTML(users []database.User) string {
	buf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"ID", "NAME", "ROLE", "CREATED", "UPDATED"})
	for _, u := range users {
		table.Append([]string{strconv.FormatInt(u.ID, 10), u.Name, u.Role, u.CreatedAt.Format(time.Layout), u.UpdatedAt.Format(time.Layout)})
	}
	table.Render()
	return "<pre>\n" + buf.String() + "\n</pre>"
}
