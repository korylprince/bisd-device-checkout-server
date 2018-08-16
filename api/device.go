package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func getDevice(ctx context.Context, bagTag string) (string, error) {
	tx := ctx.Value(InventoryTransactionKey).(*sql.Tx)

	var (
		user   *string
		status *string
		model  *string
	)
	err := tx.QueryRow("SELECT User, Status, Model FROM devices WHERE bag_tag = ?;", bagTag).Scan(
		&(user),
		&(status),
		&(model),
	)

	switch {
	case err == sql.ErrNoRows:
		return fmt.Sprintf("Bag Tag %s doesn't exist", bagTag), nil
	case err != nil:
		return "", &Error{Description: fmt.Sprintf("Could not query Device(%s)", bagTag), Err: err}
	}

	if *model != "C740-C4PE" {
		return fmt.Sprintf(`Bag Tag %s isn't a Chromebook (Model is "%s")`, bagTag, *model), nil
	}

	if *status != "Storage" {
		return fmt.Sprintf(`Bag Tag %s is not in Storage (Status is "%s")`, bagTag, *status), nil
	}

	if *user != "" {
		return fmt.Sprintf(`Bag Tag %s already has a User assigned (User is "%s")`, bagTag, *user), nil
	}

	return "", nil
}

func getDeviceList(ctx context.Context, name string) ([]int, error) {
	tx := ctx.Value(InventoryTransactionKey).(*sql.Tx)

	rows, err := tx.Query(`SELECT id FROM devices WHERE user=?;`, name)
	if err != nil {
		return nil, &Error{Description: "Could not query Device list", Err: err}
	}
	defer rows.Close()

	var devices []int

	for rows.Next() {
		var d *int
		if err := rows.Scan(&(d)); err != nil {
			return nil, &Error{Description: "Could not scan Device row", Err: err}
		}

		devices = append(devices, *d)
	}

	if err := rows.Err(); err != nil {
		return nil, &Error{Description: "Could not scan Device rows", Err: err}
	}

	return devices, nil
}

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

	deviceStatus, err := getDevice(ctx, bagTag)
	if err != nil {
		return err
	}

	if deviceStatus != "" {
		return &Error{Description: deviceStatus, Err: nil, RequestError: true}
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
