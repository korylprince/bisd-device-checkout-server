package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

//CheckoutDevice checks out the device with the given bagTag to the student with the given otherID
func CheckoutDevice(ctx context.Context, otherID, bagTag string) error {
	student, err := GetStudent(ctx, otherID)
	if err != nil {
		return err
	}

	status, err := student.Status(ctx)
	if err != nil {
		return err
	}

	if status.Type == "none" {
		return &Error{Description: fmt.Sprintf("Student unable to check out Chromebook: %s", status.Reason), Err: nil, RequestError: true}
	}

	tx := ctx.Value(InventoryTransactionKey).(*sql.Tx)
	commitUser := ctx.Value(UserKey).(*User)

	note := fmt.Sprintf("\n%s %s: Checked out Bag Tag %s (%s) to %s\n",
		time.Now().Format("01/02/06"),
		commitUser.DisplayName,
		bagTag,
		strings.Replace(status.Type, "_", " ", -1),
		student.Name(),
	)

	res, err := tx.Exec(`
	UPDATE devices SET User = ?, Status = "Checked Out", Notes = CONCAT(Notes, ?)
	WHERE bag_tag = ? AND model = "C740-C4PE" AND Status = "Storage";
	`, student.Name(), note, bagTag)

	if err != nil {
		return &Error{Description: fmt.Sprintf("Could not update Device(%s)", bagTag), Err: err}
	}

	if n, _ := res.RowsAffected(); n != 1 {
		return &Error{Description: fmt.Sprintf("Device with Bag Tag %s is missing or not in storage", bagTag), Err: nil, RequestError: true}
	}

	return nil
}
