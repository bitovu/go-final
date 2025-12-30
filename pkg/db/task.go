package db

import (
	"database/sql"
	"fmt"
)

const (
	addTaskQuery = `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`

	getTaskQuery = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?
	`

	tasksQuery = `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`

	updateTaskQuery = `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	updateDateQuery = `
		UPDATE scheduler
		SET date = ?
		WHERE id = ?
	`

	deleteTaskQuery = `
		DELETE FROM scheduler
		WHERE id = ?
	`
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	// определите запрос
	res, err := db.Exec(addTaskQuery,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := db.Query(tasksQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Title,
			&t.Comment,
			&t.Repeat,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	err := db.QueryRow(
		getTaskQuery,
		sql.Named("id", id),
	).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	res, err := db.Exec(updateDateQuery,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := db.Exec(deleteTaskQuery, id, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := db.Exec(
		updateDateQuery,
		id,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
