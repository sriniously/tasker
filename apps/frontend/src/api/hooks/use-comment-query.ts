import { type TApiClient, useApiClient } from "@/api";
import { QUERY_KEYS } from "@/api/query-utils";
import type { TRequests } from "@/api/types";
import { showApiErrorToast } from "@/api/utils";
import type { apiContract } from "@tasker/openapi/contracts";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ServerInferResponseBody } from "@ts-rest/core";

// Type definitions for Comment operations
export type TGetCommentsByTodoIdResponse = ServerInferResponseBody<
  typeof apiContract.Comment.getCommentsByTodoId,
  200
>;

export type TAddCommentPayload = TRequests["Comment"]["addComment"]["body"];
export type TAddCommentResponse = ServerInferResponseBody<
  typeof apiContract.Comment.addComment,
  201
>;

export type TUpdateCommentPayload =
  TRequests["Comment"]["updateComment"]["body"];
export type TUpdateCommentResponse = ServerInferResponseBody<
  typeof apiContract.Comment.updateComment,
  200
>;

// API functions
const fetchCommentsByTodoId = async ({
  api,
  todoId,
}: {
  api: TApiClient;
  todoId: string;
}): Promise<TGetCommentsByTodoIdResponse> => {
  const res = await api.Comment.getCommentsByTodoId({ params: { id: todoId } });

  if (res.status === 200) {
    return res.body;
  } else {
    throw res.body;
  }
};

const addComment = async ({
  api,
  todoId,
  data,
}: {
  api: TApiClient;
  todoId: string;
  data: TAddCommentPayload;
}): Promise<TAddCommentResponse> => {
  const res = await api.Comment.addComment({
    params: { id: todoId },
    body: data,
  });

  if (res.status === 201) {
    return res.body;
  } else {
    throw res.body;
  }
};

const updateComment = async ({
  api,
  id,
  data,
}: {
  api: TApiClient;
  id: string;
  data: TUpdateCommentPayload;
}): Promise<TUpdateCommentResponse> => {
  const res = await api.Comment.updateComment({ params: { id }, body: data });

  if (res.status === 200) {
    return res.body;
  } else {
    throw res.body;
  }
};

const deleteComment = async ({
  api,
  id,
}: {
  api: TApiClient;
  id: string;
}): Promise<void> => {
  const res = await api.Comment.deleteComment({ params: { id } });

  if (res.status !== 204) {
    throw res.body;
  }
};

// React Query hooks
export const useGetCommentsByTodoId = ({
  todoId,
  enabled = true,
}: {
  todoId: string;
  enabled?: boolean;
}) => {
  const api = useApiClient();

  return useQuery({
    queryKey: [QUERY_KEYS.COMMENTS.GET_COMMENTS_BY_TODO_ID, todoId],
    queryFn: () => fetchCommentsByTodoId({ api, todoId }),
    enabled: enabled && !!todoId,
    placeholderData: [],
  });
};

export const useAddComment = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ body }: { body: TAddCommentPayload & { todoId: string } }) =>
      addComment({ api, todoId: body.todoId, data: body }),
    onSuccess: (_, { body }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.COMMENTS.GET_COMMENTS_BY_TODO_ID, body.todoId],
      });
      // Also invalidate the todo to update comment count
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.GET_TODO_BY_ID, body.todoId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.ALL_TODOS],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to add comment");
    },
  });
};

export const useUpdateComment = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      commentId,
      body,
    }: {
      commentId: string;
      body: TUpdateCommentPayload;
    }) => updateComment({ api, id: commentId, data: body }),
    onSuccess: (updatedComment) => {
      // Invalidate comments for the todo this comment belongs to
      queryClient.invalidateQueries({
        queryKey: [
          QUERY_KEYS.COMMENTS.GET_COMMENTS_BY_TODO_ID,
          updatedComment.todoId,
        ],
      });
      // Also invalidate the todo
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.GET_TODO_BY_ID, updatedComment.todoId],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.ALL_TODOS],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to update comment");
    },
  });
};

export const useDeleteComment = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ commentId }: { commentId: string }) =>
      deleteComment({ api, id: commentId }),
    onMutate: async ({ commentId }) => {
      // We need to find which todo this comment belongs to
      // This is a bit tricky without the todoId, but we can get it from the cache
      const queryCache = queryClient.getQueryCache();
      const commentQueries = queryCache.findAll({
        queryKey: [QUERY_KEYS.COMMENTS.GET_COMMENTS_BY_TODO_ID],
      });

      let todoId: string | null = null;
      for (const query of commentQueries) {
        const comments = query.state.data as
          | TGetCommentsByTodoIdResponse
          | undefined;
        if (comments && Array.isArray(comments)) {
          const comment = comments.find((c) => c.id === commentId);
          if (comment) {
            todoId = comment.todoId;
            break;
          }
        }
      }

      return { todoId };
    },
    onSuccess: (_, __, context) => {
      if (context?.todoId) {
        queryClient.invalidateQueries({
          queryKey: [
            QUERY_KEYS.COMMENTS.GET_COMMENTS_BY_TODO_ID,
            context.todoId,
          ],
        });
        queryClient.invalidateQueries({
          queryKey: [QUERY_KEYS.TODOS.GET_TODO_BY_ID, context.todoId],
        });
      }
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.ALL_TODOS],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to delete comment");
    },
  });
};
