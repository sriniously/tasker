import { z } from "zod";

export type PaginatedResponse<T> = {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
};

export const schemaWithPagination = <T>(
  schema: z.ZodSchema<T>
): z.ZodSchema<PaginatedResponse<T>> =>
  z.object({
    data: z.array(schema),
    total: z.number(),
    page: z.number(),
    limit: z.number(),
    totalPages: z.number(),
  });
