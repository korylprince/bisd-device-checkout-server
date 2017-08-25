package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

//Student represents a Skyward Student
type Student struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	OtherID   string `json:"other_id"`
	Grade     int    `json:"grade"`
}

//formalizeName tries to return the pretty (capitalized) version of the given name
func formalizeName(name string) string {
	name = strings.Title(strings.ToLower(strings.TrimSpace(name)))

	//Fix Mc...
	if strings.HasPrefix(name, "Mc") && len(name) >= 3 {
		name = name[0:2] + strings.ToUpper(string(name[2])) + name[3:]
	}
	//Fix hyphenated names
	if idx := strings.Index(name, "-"); idx > -1 {
		name = name[0:idx] + "-" + formalizeName(name[idx+1:])
	}

	return name
}

//getStudentName returns the formalized name of the student with the given OtherID or an error if one occurred
func getStudentName(ctx context.Context, otherID string) (string, error) {
	tx := ctx.Value(SkywardTransactionKey).(*sql.Tx)

	var (
		firstName string
		lastName  string
	)

	row := tx.QueryRow(`
	SELECT name."FIRST-NAME", name."LAST-NAME" FROM PUB.NAME AS name
	INNER JOIN PUB.STUDENT AS student ON name."NAME-ID" = student."NAME-ID"
	WHERE student."OTHER-ID" = ?
	WITH (READPAST NOWAIT)
	`, otherID)
	err := row.Scan(
		&firstName,
		&lastName,
	)

	switch {
	case err == sql.ErrNoRows:
		return "", &Error{Description: fmt.Sprintf("Student could not be found with OtherID: %s", otherID), Err: err, RequestError: true}
	case err != nil:
		return "", &Error{Description: fmt.Sprintf("Could not query Student(%s)", otherID), Err: err}
	}

	return formalizeName(firstName) + " " + formalizeName(lastName), nil
}

//calculateCharge parses the charge syntax and returns a proper total charge
func calculateCharge(charges string) float32 {
	var total float32
	for _, charge := range strings.Split(charges, "|") {
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

//validateStudent validates that a student doesn't have a device checked out already and has no open charges
func validateStudent(ctx context.Context, name string) error {
	tx := ctx.Value(InventoryTransactionKey).(*sql.Tx)

	rows, err := tx.Query(`SELECT 1 FROM devices WHERE user = ? AND status = 'Checked Out' AND model = "C740-C4PE";`, name)
	if err != nil {
		return &Error{Description: fmt.Sprintf("Could not query devices for Student: %s", name), Err: err}
	}
	if rows.Next() {
		rows.Close()
		return &Error{Description: fmt.Sprintf("Student %s has existing device checked out", name), Err: err, RequestError: true}
	}
	rows.Close()

	rows, err = tx.Query(`SELECT amount_paid, charges FROM charges WHERE user=?;`, name)
	if err != nil {
		return &Error{Description: fmt.Sprintf("Could not query charges for Student: %s", name), Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		var paid float32
		var charges string
		if err = rows.Scan(&paid, &charges); err != nil {
			return &Error{Description: "Could not scan charge row", Err: err}
		}
		if paid < calculateCharge(charges) {
			return &Error{Description: fmt.Sprintf("Student %s has unpaid charges", name), Err: err, RequestError: true}
		}
	}

	if err := rows.Err(); err != nil {
		return &Error{Description: fmt.Sprintf("Could not scan charges for Student: %s", name), Err: err}
	}

	return nil
}

//GetStudentList returns a list of all Students
func GetStudentList(ctx context.Context) ([]*Student, error) {
	tx := ctx.Value(SkywardTransactionKey).(*sql.Tx)

	rows, err := tx.Query(`
	SELECT
		name."FIRST-NAME" AS First_Name,
		name."LAST-NAME" AS Last_Name,
		student."OTHER-ID" AS Other_I_D,
		(12 - (student."GRAD-YR" - entity."SCHOOL-YEAR")) AS Grade
	FROM PUB.NAME AS name
	INNER JOIN PUB."STUDENT" AS student ON 
		name."NAME-ID" = student."NAME-ID"

	INNER JOIN PUB."STUDENT-ENTITY" as sentity ON
		sentity."STUDENT-ID" = student."STUDENT-ID"

	INNER JOIN PUB."ENTITY" as entity ON
		entity."ENTITY-ID" = sentity."ENTITY-ID"

	WHERE sentity."STUDENT-STATUS" = 'A' AND
	student."GRAD-YR" >= entity."SCHOOL-YEAR" AND
	(student."GRAD-YR" - entity."SCHOOL-YEAR") < 6

	WITH (READPAST NOWAIT)
	`)
	if err != nil {
		return nil, &Error{Description: "Could not query Student list", Err: err}
	}
	defer rows.Close()

	var students []*Student

	for rows.Next() {
		s := new(Student)
		if err := rows.Scan(&(s.FirstName), &(s.LastName), &(s.OtherID), &(s.Grade)); err != nil {
			return nil, &Error{Description: "Could not scan Student row", Err: err}
		}

		s.FirstName = formalizeName(s.FirstName)
		s.LastName = formalizeName(s.LastName)

		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, &Error{Description: "Could not scan Student rows", Err: err}
	}

	return students, nil
}
