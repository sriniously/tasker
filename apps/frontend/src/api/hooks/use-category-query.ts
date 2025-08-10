import { type TApiClient, useApiClient } from "@/api";
import { QUERY_KEYS } from "@/api/query-utils";
import type { TRequests } from "@/api/types";
import { showApiErrorToast } from "@/api/utils";
import type { apiContract } from "@tasker/openapi/contracts";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ServerInferResponseBody } from "@ts-rest/core";

// Type definitions for Category operations
export type TGetCategoriesQuery = TRequests["Category"]["getCategories"]["query"];
export type TGetCategoriesResponse = ServerInferResponseBody<
  typeof apiContract.Category.getCategories,
  200
>;

export type TGetCategoryByIdResponse = ServerInferResponseBody<
  typeof apiContract.Category.getCategoryById,
  200
>;

export type TCreateCategoryPayload = TRequests["Category"]["createCategory"]["body"];
export type TCreateCategoryResponse = ServerInferResponseBody<
  typeof apiContract.Category.createCategory,
  201
>;

export type TUpdateCategoryPayload = TRequests["Category"]["updateCategory"]["body"];
export type TUpdateCategoryResponse = ServerInferResponseBody<
  typeof apiContract.Category.updateCategory,
  200
>;

// API functions
const fetchAllCategories = async ({
  api,
  query,
}: {
  api: TApiClient;
  query?: TGetCategoriesQuery;
}): Promise<TGetCategoriesResponse> => {
  const res = await api.Category.getCategories({ query });
  
  if (res.status === 200) {
    return res.body;
  } else {
    throw res.body;
  }
};

const fetchCategoryById = async ({
  api,
  id,
}: {
  api: TApiClient;
  id: string;
}): Promise<TGetCategoryByIdResponse> => {
  const res = await api.Category.getCategoryById({ params: { id } });
  
  if (res.status === 200) {
    return res.body;
  } else {
    throw res.body;
  }
};

const createCategory = async ({
  api,
  data,
}: {
  api: TApiClient;
  data: TCreateCategoryPayload;
}): Promise<TCreateCategoryResponse> => {
  const res = await api.Category.createCategory({ body: data });
  
  if (res.status === 201) {
    return res.body;
  } else {
    throw res.body;
  }
};

const updateCategory = async ({
  api,
  id,
  data,
}: {
  api: TApiClient;
  id: string;
  data: TUpdateCategoryPayload;
}): Promise<TUpdateCategoryResponse> => {
  const res = await api.Category.updateCategory({ params: { id }, body: data });
  
  if (res.status === 200) {
    return res.body;
  } else {
    throw res.body;
  }
};

const deleteCategory = async ({
  api,
  id,
}: {
  api: TApiClient;
  id: string;
}): Promise<void> => {
  const res = await api.Category.deleteCategory({ params: { id } });
  
  if (res.status !== 204) {
    throw res.body;
  }
};

// React Query hooks
export const useGetAllCategories = ({
  query,
}: {
  query?: TGetCategoriesQuery;
} = {}) => {
  const api = useApiClient();

  return useQuery({
    queryKey: [QUERY_KEYS.CATEGORIES.ALL_CATEGORIES, query],
    queryFn: () => fetchAllCategories({ api, query }),
    placeholderData: {
      data: [],
      total: 0,
      page: 1,
      limit: 50,
      totalPages: 0,
    },
  });
};

export const useGetCategoryById = ({
  id,
  enabled = true,
}: {
  id: string;
  enabled?: boolean;
}) => {
  const api = useApiClient();

  return useQuery({
    queryKey: [QUERY_KEYS.CATEGORIES.GET_CATEGORY_BY_ID, id],
    queryFn: () => fetchCategoryById({ api, id }),
    enabled: enabled && !!id,
  });
};

export const useCreateCategory = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ body }: { body: TCreateCategoryPayload }) => createCategory({ api, data: body }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.CATEGORIES.ALL_CATEGORIES],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to create category");
    },
  });
};

export const useUpdateCategory = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ categoryId, body }: { categoryId: string; body: TUpdateCategoryPayload }) =>
      updateCategory({ api, id: categoryId, data: body }),
    onSuccess: (_, { categoryId }) => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.CATEGORIES.ALL_CATEGORIES],
      });
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.CATEGORIES.GET_CATEGORY_BY_ID, categoryId],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to update category");
    },
  });
};

export const useDeleteCategory = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ categoryId }: { categoryId: string }) => deleteCategory({ api, id: categoryId }),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.CATEGORIES.ALL_CATEGORIES],
      });
      // Also invalidate todos since they might reference deleted categories
      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.TODOS.ALL_TODOS],
      });
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to delete category");
    },
  });
};