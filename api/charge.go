package api

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

//Charge represents an inventory charge
type Charge struct {
	ID         int
	AmountPaid float32
	charges    string
}

//AmountCharged is the total amount charged
func (c *Charge) AmountCharged() float32 {
	var total float32
	for _, charge := range strings.Split(c.charges, "|") {
		if split := strings.Split(strings.TrimSpace(charge), ":"); len(split) == 2 {
			c, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 32)
			if err != nil {
				continue
			}
			total += float32(c)
		}
	}
	return total
}

func getChargeList(ctx context.Context, name string) ([]*Charge, error) {
	tx := ctx.Value(InventoryTransactionKey).(*sql.Tx)

	rows, err := tx.Query(`SELECT id, amount_paid, charges FROM charges WHERE user=?;`, name)
	if err != nil {
		return nil, &Error{Description: "Could not query Charge list", Err: err}
	}
	defer rows.Close()

	var charges []*Charge

	for rows.Next() {
		c := new(Charge)
		if err := rows.Scan(&(c.ID), &(c.AmountPaid), &(c.charges)); err != nil {
			return nil, &Error{Description: "Could not scan Charge row", Err: err}
		}

		charges = append(charges, c)
	}

	if err := rows.Err(); err != nil {
		return nil, &Error{Description: "Could not scan Charge rows", Err: err}
	}

	return charges, nil
}
