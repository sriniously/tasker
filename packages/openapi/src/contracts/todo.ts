import { getSecurityMetadata } from "../utils.js";
import {
  schemaWithPagination,
  ZPopulatedTodo,
  ZTodo,
  ZTodoAttachment,
  ZTodoStats,
} from "@tasker/zod";
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const todoContract = c.router(
  {
    getTodos: {
      summary: "Get all todos",
      path: "/todos",
      method: "GET",
      description: "Get all todos",
      query: z.object({
        page: z.number().min(1).optional(),
        limit: z.number().min(1).max(100).optional(),
        sort: z
          .enum([
            "created_at",
            "updated_at",
            "title",
            "priority",
            "due_date",
            "status",
          ])
          .optional(),
        order: z.enum(["asc", "desc"]).optional(),
        search: z.string().min(1).optional(),
        status: ZTodo.shape.status.optional(),
        priority: ZTodo.shape.priority.optional(),
        categoryId: z.string().uuid().optional(),
        parentTodoId: z.string().uuid().optional(),
        dueFrom: z.string().datetime().optional(),
        dueTo: z.string().datetime().optional(),
        overdue: z.boolean().optional(),
        completed: z.boolean().optional(),
      }),
      responses: {
        200: schemaWithPagination(ZPopulatedTodo),
      },
      metadata: metadata,
    },

    createTodo: {
      summary: "Create a new todo",
      path: "/todos",
      method: "POST",
      description: "Create a new todo",
      body: ZTodo.pick({
        title: true,
        description: true,
        priority: true,
        dueDate: true,
        parentTodoId: true,
        categoryId: true,
        metadata: true,
      })
        .partial()
        .required({
          title: true,
        }),
      responses: {
        201: ZTodo,
      },
      metadata: metadata,
    },

    getTodoById: {
      summary: "Get todo by ID",
      path: "/todos/:id",
      method: "GET",
      description: "Get todo by ID",
      responses: {
        200: ZPopulatedTodo,
      },
      metadata: metadata,
    },

    updateTodo: {
      summary: "Update todo",
      path: "/todos/:id",
      method: "PATCH",
      description: "Update todo",
      body: ZTodo.pick({
        title: true,
        description: true,
        status: true,
        priority: true,
        dueDate: true,
        parentTodoId: true,
        categoryId: true,
        metadata: true,
      }).partial(),
      responses: {
        200: ZTodo,
      },
      metadata: metadata,
    },

    deleteTodo: {
      summary: "Delete todo",
      path: "/todos/:id",
      method: "DELETE",
      description: "Delete todo",
      responses: {
        204: z.void(),
      },
      metadata: metadata,
    },

    getTodoStats: {
      summary: "Get todo statistics",
      path: "/todos/stats",
      method: "GET",
      description: "Get todo statistics",
      responses: {
        200: ZTodoStats,
      },
      metadata: metadata,
    },

    uploadTodoAttachment: {
      summary: "Upload attachment to todo",
      path: "/todos/:id/attachments",
      method: "POST",
      description: "Upload a file attachment to a todo",
      contentType: "multipart/form-data",
      body: z.object({
        file: z.object({
          type: z.literal("file"),
        }),
      }),
      responses: {
        201: ZTodoAttachment,
      },
      metadata: metadata,
    },

    deleteTodoAttachment: {
      summary: "Delete todo attachment",
      path: "/todos/:id/attachments/:attachmentId",
      method: "DELETE",
      description: "Delete a file attachment from a todo",
      responses: {
        204: z.void(),
      },
      metadata: metadata,
    },

    getAttachmentPresignedURL: {
      summary: "Get attachment download URL",
      path: "/todos/:id/attachments/:attachmentId/download",
      method: "GET",
      description: "Get a presigned URL to download an attachment",
      responses: {
        200: z.object({
          url: z.string().url(),
        }),
      },
      metadata: metadata,
    },
  },
  {
    pathPrefix: "/v1",
  }
);
