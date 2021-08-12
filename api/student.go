package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

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

	//Fix name with multiple names
	if len(strings.Split(name, " ")) > 1 {
		var names []string
		for _, n := range strings.Split(name, " ") {
			names = append(names, formalizeName(n))
		}
		return strings.Join(names, " ")
	}

	return name
}

//Student represents a Skyward Student
type Student struct {
	FirstName  string
	LastName   string
	OtherID    string
	Grade      int
	T2E2Status *string
}

//Name returns to formalized name of the Student
func (s Student) Name() string {
	return fmt.Sprintf("%s %s", s.FirstName, s.LastName)
}

//GetStudent returns the Student with the given otherID
func GetStudent(ctx context.Context, otherID string) (*Student, error) {
	tx := ctx.Value(SkywardTransactionKey).(*sql.Tx)

	s := &Student{}

	var (
		firstName *string
		lastName  *string
	)

	err := tx.QueryRow(`
	SELECT
		name."FIRST-NAME" AS First_Name,
		name."LAST-NAME" AS Last_Name,
		student."OTHER-ID" AS Other_I_D,
		(12 - (student."GRAD-YR" - entity."SCHOOL-YEAR")) AS Grade,
		data."STATUS" AS Status
	FROM PUB.NAME AS name
	INNER JOIN PUB."STUDENT" AS student ON
		name."NAME-ID" = student."NAME-ID"

	INNER JOIN PUB."STUDENT-ENTITY" as sentity ON
		sentity."STUDENT-ID" = student."STUDENT-ID"

	INNER JOIN PUB."ENTITY" as entity ON
		entity."ENTITY-ID" = sentity."ENTITY-ID"

    LEFT JOIN (
		SELECT
			data."QUDDAT-SRC-ID" AS "STUDENT-ID",
			data."QUDDAT-CHAR" AS "STATUS"

		FROM PUB."QUDDAT-DATA" AS data

		INNER JOIN PUB."QUDTBL-TABLES" AS tables ON
			data."QUDDAT-STORAGE-TYPE" = 'Custom Student' AND
			data."QUDTBL-TABLE-ID" = tables."QUDTBL-TABLE-ID" AND
			tables."QUDTBL-DESC" = '21-22 T2E2'

		INNER JOIN PUB."QUDFLD-FIELDS" AS fields ON
			data."QUDFLD-FIELD-ID" = fields."QUDFLD-FIELD-ID" AND
			fields."QUDFLD-FIELD-LABEL" = 'Can_Take_Chromebook_Home'
    ) AS data ON
		student."STUDENT-ID" = data."STUDENT-ID"

	WHERE sentity."STUDENT-STATUS" = 'A' AND
	student."GRAD-YR" >= entity."SCHOOL-YEAR" AND
	(student."GRAD-YR" - entity."SCHOOL-YEAR") < 6 AND
	student."OTHER-ID" = ?

	WITH (NOLOCK)
	`, otherID).Scan(
		&(firstName),
		&(lastName),
		&(s.OtherID),
		&(s.Grade),
		&(s.T2E2Status),
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, &Error{Description: fmt.Sprintf("Student could not be found with OtherID: %s", otherID), Err: err, RequestError: true}
	case err != nil:
		return nil, &Error{Description: fmt.Sprintf("Could not query Student(%s)", otherID), Err: err}
	}

	s.FirstName = formalizeName(*firstName)
	s.LastName = formalizeName(*lastName)

	return s, nil
}

//GetStudentList returns a list of all Students
func GetStudentList(ctx context.Context) ([]*Student, error) {
	tx := ctx.Value(SkywardTransactionKey).(*sql.Tx)

	rows, err := tx.Query(`
	SELECT
		name."FIRST-NAME" AS First_Name,
		name."LAST-NAME" AS Last_Name,
		student."OTHER-ID" AS Other_I_D,
		(12 - (student."GRAD-YR" - entity."SCHOOL-YEAR")) AS Grade,
		data."STATUS" AS Status
	FROM PUB.NAME AS name
	INNER JOIN PUB."STUDENT" AS student ON
		name."NAME-ID" = student."NAME-ID"

	INNER JOIN PUB."STUDENT-ENTITY" as sentity ON
		sentity."STUDENT-ID" = student."STUDENT-ID"

	INNER JOIN PUB."ENTITY" as entity ON
		entity."ENTITY-ID" = sentity."ENTITY-ID"

    LEFT JOIN (
		SELECT
			data."QUDDAT-SRC-ID" AS "STUDENT-ID",
			data."QUDDAT-CHAR" AS "STATUS"

		FROM PUB."QUDDAT-DATA" AS data

		INNER JOIN PUB."QUDTBL-TABLES" AS tables ON
			data."QUDDAT-STORAGE-TYPE" = 'Custom Student' AND
			data."QUDTBL-TABLE-ID" = tables."QUDTBL-TABLE-ID" AND
			tables."QUDTBL-DESC" = '21-22 T2E2'

		INNER JOIN PUB."QUDFLD-FIELDS" AS fields ON
			data."QUDFLD-FIELD-ID" = fields."QUDFLD-FIELD-ID" AND
			fields."QUDFLD-FIELD-LABEL" = 'Can_Take_Chromebook_Home'
    ) AS data ON
		student."STUDENT-ID" = data."STUDENT-ID"

	WHERE sentity."STUDENT-STATUS" = 'A' AND
	student."GRAD-YR" >= entity."SCHOOL-YEAR" AND
	(student."GRAD-YR" - entity."SCHOOL-YEAR") < 6

	WITH (NOLOCK)
	`)
	if err != nil {
		return nil, &Error{Description: "Could not query Student list", Err: err}
	}
	defer rows.Close()

	var students []*Student

	var (
		firstName *string
		lastName  *string
	)

	for rows.Next() {
		s := new(Student)
		if err := rows.Scan(&(firstName), &(lastName), &(s.OtherID), &(s.Grade), &(s.T2E2Status)); err != nil {
			return nil, &Error{Description: "Could not scan Student row", Err: err}
		}

		s.FirstName = formalizeName(*firstName)
		s.LastName = formalizeName(*lastName)

		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, &Error{Description: "Could not scan Student rows", Err: err}
	}

	return students, nil
}
