export const QUERY_KEYS = {
  TODOS: {
    ALL_TODOS: "allTodos",
    GET_TODO_BY_ID: "getTodoById",
    TODO_STATS: "todoStats",
    TODO_ATTACHMENTS: "todoAttachments",
  },
  CATEGORIES: {
    ALL_CATEGORIES: "allCategories", 
    GET_CATEGORY_BY_ID: "getCategoryById",
  },
  COMMENTS: {
    GET_COMMENTS_BY_TODO_ID: "getCommentsByTodoId",
  },
} as const satisfies Record<Uppercase<string>, object>;