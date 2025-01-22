package service

import (
	"task-viewer/db"
	"task-viewer/model"
	"time"
)

// CreateTask 创建任务
func CreateTask(name, description string) (*model.Task, error) {
	result, err := db.DB.Exec(`
        INSERT INTO tasks (name, description)
        VALUES (?, ?)
    `, name, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetTask(int(id))
}

// GetTask 获取任务
func GetTask(id int) (*model.Task, error) {
	var task model.Task
	err := db.DB.QueryRow(`
        SELECT id, name, description, created_at, updated_at
        FROM tasks
        WHERE id = ?
    `, id).Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllTasks 获取所有任务
func GetAllTasks() ([]*model.Task, error) {
	rows, err := db.DB.Query(`
        SELECT id, name, description, created_at, updated_at
        FROM tasks
        ORDER BY id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

// UpdateTask 更新任务
func UpdateTask(id int, name, description string) (*model.Task, error) {
	_, err := db.DB.Exec(`
        UPDATE tasks
        SET name = ?, description = ?, updated_at = ?
        WHERE id = ?
    `, name, description, time.Now(), id)
	if err != nil {
		return nil, err
	}
	return GetTask(id)
}

// DeleteTask 删除任务
func DeleteTask(id int) error {
	// 首先删除相关的依赖关系
	_, err := db.DB.Exec(`
        DELETE FROM task_dependencies
        WHERE source_id = ? OR target_id = ?
    `, id, id)
	if err != nil {
		return err
	}

	// 然后删除任务
	_, err = db.DB.Exec(`
        DELETE FROM tasks
        WHERE id = ?
    `, id)
	return err
}

// AddDependency 添加依赖关系
func AddDependency(sourceID, targetID int) error {
	_, err := db.DB.Exec(`
        INSERT INTO task_dependencies (source_id, target_id)
        VALUES (?, ?)
    `, sourceID, targetID)
	return err
}

// RemoveDependency 移除依赖关系
func RemoveDependency(sourceID, targetID int) error {
	_, err := db.DB.Exec(`
        DELETE FROM task_dependencies
        WHERE source_id = ? AND target_id = ?
    `, sourceID, targetID)
	return err
}

// GetTaskDependencies 获取任务的依赖关系
func GetTaskDependencies(taskID int) ([]*model.Task, []*model.Task, error) {
	// 获取前置依赖
	upstreamRows, err := db.DB.Query(`
        SELECT t.id, t.name, t.description, t.created_at, t.updated_at
        FROM tasks t
        JOIN task_dependencies td ON t.id = td.source_id
        WHERE td.target_id = ?
    `, taskID)
	if err != nil {
		return nil, nil, err
	}
	defer upstreamRows.Close()

	var upstream []*model.Task
	for upstreamRows.Next() {
		var task model.Task
		err := upstreamRows.Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, nil, err
		}
		upstream = append(upstream, &task)
	}

	// 获取后置依赖
	downstreamRows, err := db.DB.Query(`
        SELECT t.id, t.name, t.description, t.created_at, t.updated_at
        FROM tasks t
        JOIN task_dependencies td ON t.id = td.target_id
        WHERE td.source_id = ?
    `, taskID)
	if err != nil {
		return nil, nil, err
	}
	defer downstreamRows.Close()

	var downstream []*model.Task
	for downstreamRows.Next() {
		var task model.Task
		err := downstreamRows.Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, nil, err
		}
		downstream = append(downstream, &task)
	}

	return upstream, downstream, nil
}

// GetAllDependencies 获取所有依赖关系
func GetAllDependencies() ([]*model.TaskDependency, error) {
	rows, err := db.DB.Query(`
        SELECT id, source_id, target_id, created_at
        FROM task_dependencies
        ORDER BY id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dependencies []*model.TaskDependency
	for rows.Next() {
		var dep model.TaskDependency
		err := rows.Scan(&dep.ID, &dep.SourceID, &dep.TargetID, &dep.CreatedAt)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, &dep)
	}
	return dependencies, nil
}
