package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/tasker/internal/model/todo"
	"github.com/sriniously/tasker/internal/repository"
	testing_pkg "github.com/sriniously/tasker/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoRepository_CreateTodo(t *testing.T) {
	_, testServer, cleanup := testing_pkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	todoRepo := repository.NewTodoRepository(testServer)

	t.Run("create todo successfully", func(t *testing.T) {
		userID := uuid.New().String()
		dueDate := time.Now().Add(24 * time.Hour)
		payload := &todo.CreateTodoPayload{
			Title:       "Test Todo",
			Description: testing_pkg.Ptr("Test todo description"),
			Priority:    testing_pkg.Ptr(todo.PriorityHigh),
			DueDate:     &dueDate,
		}

		result, err := todoRepo.CreateTodo(ctx, userID, payload)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.NotEqual(t, uuid.Nil, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, payload.Title, result.Title)
		assert.Equal(t, payload.Description, result.Description)
		assert.Equal(t, *payload.Priority, result.Priority)
		assert.Equal(t, payload.DueDate.Unix(), result.DueDate.Unix())
		assert.Equal(t, todo.StatusDraft, result.Status)
		assert.Nil(t, result.CompletedAt)
		testing_pkg.AssertTimestampsValid(t, result)
	})

	t.Run("create todo with minimum required fields", func(t *testing.T) {
		userID := uuid.New().String()
		payload := &todo.CreateTodoPayload{
			Title: "Minimal Todo",
		}

		result, err := todoRepo.CreateTodo(ctx, userID, payload)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, payload.Title, result.Title)
		assert.Nil(t, result.Description)
		assert.Equal(t, todo.PriorityMedium, result.Priority)
		assert.Nil(t, result.DueDate)
	})

	t.Run("create todo with metadata", func(t *testing.T) {
		userID := uuid.New().String()
		metadata := &todo.Metadata{
			Tags:  []string{"work", "urgent"},
			Color: testing_pkg.Ptr("#ff0000"),
		}
		payload := &todo.CreateTodoPayload{
			Title:    "Todo with Metadata",
			Metadata: metadata,
		}

		result, err := todoRepo.CreateTodo(ctx, userID, payload)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, metadata.Tags, result.Metadata.Tags)
		assert.Equal(t, metadata.Color, result.Metadata.Color)
	})

	t.Run("with canceled context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()

		userID := uuid.New().String()
		payload := &todo.CreateTodoPayload{
			Title: "Canceled Todo",
		}

		result, err := todoRepo.CreateTodo(canceledCtx, userID, payload)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTodoRepository_GetTodoByID(t *testing.T) {
	_, testServer, cleanup := testing_pkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	todoRepo := repository.NewTodoRepository(testServer)

	// Create test todo
	userID := uuid.New().String()
	testTodo := createTestTodo(t, ctx, todoRepo, userID)

	t.Run("get todo by id successfully", func(t *testing.T) {
		result, err := todoRepo.GetTodoByID(ctx, userID, testTodo.ID)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, testTodo.ID, result.ID)
		assert.Equal(t, testTodo.Title, result.Title)
		assert.Equal(t, testTodo.UserID, result.UserID)
		assert.NotNil(t, result.Children)
		assert.NotNil(t, result.Comments)
	})

	t.Run("get non-existent todo", func(t *testing.T) {
		nonExistentID := uuid.New()

		result, err := todoRepo.GetTodoByID(ctx, userID, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("get todo with wrong user id", func(t *testing.T) {
		wrongUserID := uuid.New().String()

		result, err := todoRepo.GetTodoByID(ctx, wrongUserID, testTodo.ID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("with canceled context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()

		result, err := todoRepo.GetTodoByID(canceledCtx, userID, testTodo.ID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTodoRepository_CheckTodoExists(t *testing.T) {
	_, testServer, cleanup := testing_pkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	todoRepo := repository.NewTodoRepository(testServer)

	// Create test todo
	userID := uuid.New().String()
	testTodo := createTestTodo(t, ctx, todoRepo, userID)

	t.Run("check existing todo", func(t *testing.T) {
		result, err := todoRepo.CheckTodoExists(ctx, userID, testTodo.ID)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, testTodo.ID, result.ID)
		assert.Equal(t, testTodo.Title, result.Title)
	})

	t.Run("check non-existent todo", func(t *testing.T) {
		nonExistentID := uuid.New()

		result, err := todoRepo.CheckTodoExists(ctx, userID, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("check todo with wrong user id", func(t *testing.T) {
		wrongUserID := uuid.New().String()

		result, err := todoRepo.CheckTodoExists(ctx, wrongUserID, testTodo.ID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTodoRepository_GetTodos(t *testing.T) {
	_, testServer, cleanup := testing_pkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	todoRepo := repository.NewTodoRepository(testServer)

	// Create test todos
	userID := uuid.New().String()
	_ = createTestTodos(t, ctx, todoRepo, userID, 3)

	t.Run("get todos with default pagination", func(t *testing.T) {
		page := 1
		limit := 20
		query := &todo.GetTodosQuery{
			Page:  &page,
			Limit: &limit,
		}

		result, err := todoRepo.GetTodos(ctx, userID, query)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Data), 3)
		assert.Equal(t, page, result.Page)
		assert.Equal(t, limit, result.Limit)
		assert.GreaterOrEqual(t, result.Total, 3)
		assert.GreaterOrEqual(t, result.TotalPages, 1)
	})

	t.Run("get todos with pagination", func(t *testing.T) {
		page := 1
		limit := 2
		query := &todo.GetTodosQuery{
			Page:  &page,
			Limit: &limit,
		}

		result, err := todoRepo.GetTodos(ctx, userID, query)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, page, result.Page)
		assert.Equal(t, limit, result.Limit)
		assert.GreaterOrEqual(t, result.Total, 3)
		assert.GreaterOrEqual(t, result.TotalPages, 2)
	})

	t.Run("filter by status", func(t *testing.T) {
		page := 1
		limit := 20
		status := todo.StatusDraft
		query := &todo.GetTodosQuery{
			Page:   &page,
			Limit:  &limit,
			Status: &status,
		}

		result, err := todoRepo.GetTodos(ctx, userID, query)
		require.NoError(t, err)
		require.NotNil(t, result)

		for _, todoItem := range result.Data {
			assert.Equal(t, todo.StatusDraft, todoItem.Status)
		}
	})

	t.Run("filter by priority", func(t *testing.T) {
		page := 1
		limit := 20
		priority := todo.PriorityHigh
		query := &todo.GetTodosQuery{
			Page:     &page,
			Limit:    &limit,
			Priority: &priority,
		}

		result, err := todoRepo.GetTodos(ctx, userID, query)
		require.NoError(t, err)
		require.NotNil(t, result)

		for _, todoItem := range result.Data {
			assert.Equal(t, todo.PriorityHigh, todoItem.Priority)
		}
	})

	t.Run("search by title", func(t *testing.T) {
		page := 1
		limit := 20
		search := "Test"
		query := &todo.GetTodosQuery{
			Page:   &page,
			Limit:  &limit,
			Search: &search,
		}

		result, err := todoRepo.GetTodos(ctx, userID, query)
		require.NoError(t, err)
		require.NotNil(t, result)

		for _, todoItem := range result.Data {
			assert.Contains(t, todoItem.Title, "Test")
		}
	})

	t.Run("with canceled context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()

		page := 1
		limit := 20
		query := &todo.GetTodosQuery{
			Page:  &page,
			Limit: &limit,
		}

		result, err := todoRepo.GetTodos(canceledCtx, userID, query)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTodoRepository_UpdateTodo(t *testing.T) {
	_, testServer, cleanup := testing_pkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	todoRepo := repository.NewTodoRepository(testServer)

	// Create test todo
	userID := uuid.New().String()
	testTodo := createTestTodo(t, ctx, todoRepo, userID)

	t.Run("update todo title successfully", func(t *testing.T) {
		newTitle := "Updated Todo Title"
		payload := &todo.UpdateTodoPayload{
			ID:    testTodo.ID,
			Title: &newTitle,
		}

		result, err := todoRepo.UpdateTodo(ctx, userID, payload)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, testTodo.ID, result.ID)
		assert.Equal(t, newTitle, result.Title)
		assert.True(t, result.UpdatedAt.After(testTodo.UpdatedAt))
	})

	t.Run("update todo status to completed", func(t *testing.T) {
		status := todo.StatusCompleted
		payload := &todo.UpdateTodoPayload{
			ID:     testTodo.ID,
			Status: &status,
		}

		result, err := todoRepo.UpdateTodo(ctx, userID, payload)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, todo.StatusCompleted, result.Status)
		assert.NotNil(t, result.CompletedAt)
	})

	t.Run("update multiple fields successfully", func(t *testing.T) {
		newTitle := "Multi Update Todo"
		newPriority := todo.PriorityLow
		payload := &todo.UpdateTodoPayload{
			ID:       testTodo.ID,
			Title:    &newTitle,
			Priority: &newPriority,
		}

		result, err := todoRepo.UpdateTodo(ctx, userID, payload)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, newTitle, result.Title)
		assert.Equal(t, newPriority, result.Priority)
	})

	t.Run("update with no fields should fail", func(t *testing.T) {
		payload := &todo.UpdateTodoPayload{
			ID: testTodo.ID,
		}

		result, err := todoRepo.UpdateTodo(ctx, userID, payload)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no fields to update")
	})

	t.Run("update non-existent todo", func(t *testing.T) {
		nonExistentID := uuid.New()
		newTitle := "Non Existent Todo"
		payload := &todo.UpdateTodoPayload{
			ID:    nonExistentID,
			Title: &newTitle,
		}

		result, err := todoRepo.UpdateTodo(ctx, userID, payload)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("with canceled context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()

		newTitle := "Canceled Update"
		payload := &todo.UpdateTodoPayload{
			ID:    testTodo.ID,
			Title: &newTitle,
		}

		result, err := todoRepo.UpdateTodo(canceledCtx, userID, payload)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTodoRepository_DeleteTodo(t *testing.T) {
	_, testServer, cleanup := testing_pkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	todoRepo := repository.NewTodoRepository(testServer)

	// Create test todo
	userID := uuid.New().String()
	testTodo := createTestTodo(t, ctx, todoRepo, userID)

	t.Run("delete todo successfully", func(t *testing.T) {
		err := todoRepo.DeleteTodo(ctx, userID, testTodo.ID)
		require.NoError(t, err)

		// Verify todo is deleted
		result, err := todoRepo.GetTodoByID(ctx, userID, testTodo.ID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("delete non-existent todo", func(t *testing.T) {
		nonExistentID := uuid.New()

		err := todoRepo.DeleteTodo(ctx, userID, nonExistentID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "todo not found")
	})

	t.Run("with canceled context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()

		testTodo := createTestTodo(t, ctx, todoRepo, userID)

		err := todoRepo.DeleteTodo(canceledCtx, userID, testTodo.ID)
		assert.Error(t, err)
	})
}

func TestTodoRepository_GetTodoStats(t *testing.T) {
	_, testServer, cleanup := testing_pkg.SetupTest(t)
	defer cleanup()

	ctx := context.Background()
	todoRepo := repository.NewTodoRepository(testServer)

	// Create test todos with different statuses
	userID := uuid.New().String()
	createTestTodos(t, ctx, todoRepo, userID, 5)

	t.Run("get todo stats successfully", func(t *testing.T) {
		result, err := todoRepo.GetTodoStats(ctx, userID)
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.GreaterOrEqual(t, result.Total, 5)
		assert.GreaterOrEqual(t, result.Draft, 0)
		assert.GreaterOrEqual(t, result.Active, 0)
		assert.GreaterOrEqual(t, result.Completed, 0)
		assert.GreaterOrEqual(t, result.Archived, 0)
		assert.GreaterOrEqual(t, result.Overdue, 0)
	})

	t.Run("with canceled context", func(t *testing.T) {
		canceledCtx, cancel := context.WithCancel(ctx)
		cancel()

		result, err := todoRepo.GetTodoStats(canceledCtx, userID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func createTestTodo(t *testing.T, ctx context.Context, repo *repository.TodoRepository, userID string) *todo.Todo {
	t.Helper()

	dueDate := time.Now().Add(24 * time.Hour)
	payload := &todo.CreateTodoPayload{
		Title:       "Test Todo",
		Description: testing_pkg.Ptr("Test todo description"),
		Priority:    testing_pkg.Ptr(todo.PriorityHigh),
		DueDate:     &dueDate,
	}

	result, err := repo.CreateTodo(ctx, userID, payload)
	require.NoError(t, err)

	return result
}

func createTestTodos(t *testing.T, ctx context.Context, repo *repository.TodoRepository, userID string, count int) []*todo.Todo {
	t.Helper()

	todos := make([]*todo.Todo, 0, count)

	for i := 0; i < count; i++ {
		dueDate := time.Now().Add(time.Duration(i+1) * 24 * time.Hour)
		payload := &todo.CreateTodoPayload{
			Title:       fmt.Sprintf("Test Todo %d", i+1),
			Description: testing_pkg.Ptr(fmt.Sprintf("Test todo description %d", i+1)),
			Priority:    testing_pkg.Ptr(todo.PriorityHigh),
			DueDate:     &dueDate,
		}

		result, err := repo.CreateTodo(ctx, userID, payload)
		require.NoError(t, err)
		todos = append(todos, result)

		// Add a small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	return todos
}
