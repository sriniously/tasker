import { z } from "zod";

const envVarsSchema = z.object({
  VITE_CLERK_PUBLISHABLE_KEY: z
    .string()
    .min(1, "VITE_CLERK_PUBLISHABLE_KEY is required"),
  VITE_API_URL: z.string().url().default("http://localhost:3000"),
  VITE_ENV: z.enum(["production", "development", "local"]).default("local"),
});

const parseResult = envVarsSchema.safeParse(process.env);

if (!parseResult.success) {
  console.error(
    "❌ Invalid environment variables:",
    parseResult.error.flatten(),
  );
  throw new Error("Invalid environment variables");
}

export const env = parseResult.data;
