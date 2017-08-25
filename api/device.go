package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

//CheckoutDevice checks out the device with the given Bag Tag to the student with the given Other ID
func CheckoutDevice(ctx context.Context, bagTag, otherID string, redBag bool) error {
	user, err := getStudentName(ctx, otherID)
	if err != nil {
		return err
	}

	err = validateStudent(ctx, user, redBag)
	if err != nil {
		return err
	}

	tx := ctx.Value(InventoryTransactionKey).(*sql.Tx)
	commitUser := ctx.Value(UserKey).(*User)

	note := fmt.Sprintf("\n%s %s: Checked out Bag Tag %s to %s\n",
		time.Now().Format("01/02/06"),
		commitUser.DisplayName,
		bagTag,
		user,
	)

	res, err := tx.Exec(`
	UPDATE devices SET User = ?, Status = "Checked Out", Notes = CONCAT(Notes, ?)
	WHERE bag_tag = ? AND model = "C740-C4PE" AND Status = "Storage";
	`, user, note, bagTag)

	if err != nil {
		return &Error{Description: fmt.Sprintf("Could not update Device(%s)", bagTag), Err: err}
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return &Error{Description: fmt.Sprintf("Device with Bag Tag %s is missing or not in storage", bagTag), Err: nil, RequestError: true}
	}

	return nil
}
