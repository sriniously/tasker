import { type TApiClient, useApiClient } from "@/api";
import { QUERY_KEYS } from "@/api/query-utils";
import type { TRequests } from "@/api/types";
import { showApiErrorToast } from "@/api/utils";
import type { apiContract } from "@tasker/openapi/contracts";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { ServerInferResponseBody } from "@ts-rest/core";

// Type definitions for Todo Attachment operations
export type TUploadTodoAttachmentParams =
  TRequests["Todo"]["uploadTodoAttachment"]["params"];
export type TUploadTodoAttachmentBody =
  TRequests["Todo"]["uploadTodoAttachment"]["body"];
export type TUploadTodoAttachmentResponse = ServerInferResponseBody<
  typeof apiContract.Todo.uploadTodoAttachment,
  201
>;

export type TDeleteTodoAttachmentParams =
  TRequests["Todo"]["deleteTodoAttachment"]["params"];
export type TDeleteTodoAttachmentResponse = ServerInferResponseBody<
  typeof apiContract.Todo.deleteTodoAttachment,
  204
>;

export type TGetTodoAttachmentPresignedUrlParams =
  TRequests["Todo"]["getAttachmentPresignedURL"]["params"];
export type TGetTodoAttachmentPresignedUrlResponse = ServerInferResponseBody<
  typeof apiContract.Todo.getAttachmentPresignedURL,
  200
>;

// API functions
const uploadTodoAttachment = async ({
  api,
  todoId,
  file,
}: {
  api: TApiClient;
  todoId: string;
  file: File;
}): Promise<TUploadTodoAttachmentResponse> => {
  const formData = new FormData();
  formData.append("file", file);

  const res = await api.Todo.uploadTodoAttachment({
    params: { id: todoId },
    body: formData,
    extraHeaders: {
      "Content-Type": "multipart/form-data",
    },
  });

  if (res.status === 201) {
    return res.body;
  } else {
    throw res.body;
  }
};

const deleteTodoAttachment = async ({
  api,
  todoId,
  attachmentId,
}: {
  api: TApiClient;
  todoId: string;
  attachmentId: string;
}): Promise<void> => {
  const res = await api.Todo.deleteTodoAttachment({
    params: { id: todoId, attachmentId },
  });

  if (res.status !== 204) {
    throw res.body;
  }
};

const getTodoAttachmentPresignedUrl = async ({
  api,
  todoId,
  attachmentId,
}: {
  api: TApiClient;
  todoId: string;
  attachmentId: string;
}): Promise<TGetTodoAttachmentPresignedUrlResponse> => {
  const res = await api.Todo.getAttachmentPresignedURL({
    params: { id: todoId, attachmentId },
  });

  if (res.status === 200) {
    return res.body;
  } else {
    throw res.body;
  }
};

// React Query hooks
export const useUploadTodoAttachment = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ todoId, file }: { todoId: string; file: File }) =>
      uploadTodoAttachment({ api, todoId, file }),
    onSuccess: (_, { todoId }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.GET_TODO_BY_ID, todoId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.TODO_ATTACHMENTS, todoId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.ALL_TODOS],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to upload attachment");
    },
  });
};

export const useDeleteTodoAttachment = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      todoId,
      attachmentId,
    }: {
      todoId: string;
      attachmentId: string;
    }) => deleteTodoAttachment({ api, todoId, attachmentId }),
    onSuccess: (_, { todoId }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.GET_TODO_BY_ID, todoId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.TODO_ATTACHMENTS, todoId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.ALL_TODOS],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to delete attachment");
    },
  });
};

export const useGetTodoAttachmentPresignedUrl = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: ({
      todoId,
      attachmentId,
    }: {
      todoId: string;
      attachmentId: string;
    }) => getTodoAttachmentPresignedUrl({ api, todoId, attachmentId }),
    onError: (err) => {
      showApiErrorToast(err, "Failed to get download link");
    },
  });
};
