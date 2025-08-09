import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { z } from "zod";

extendZodWithOpenApi(z);

export * from "./utils.js";
export * from "./health.js";
export * from "./todo/index.js";
export * from "./category/index.js";
export * from "./comment/index.js";
