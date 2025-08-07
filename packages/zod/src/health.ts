import { z } from "zod";

const ZHealthCheck = z.object({
  status: z.string(),
  response_time: z.string(),
  error: z.string().optional(),
});

export const ZHealthResponse = z.object({
  status: z.enum(["healthy", "unhealthy"]),
  timestamp: z.string().datetime(),
  environment: z.string(),
  checks: z.object({
    database: ZHealthCheck,
    redis: ZHealthCheck.optional(),
  }),
});
