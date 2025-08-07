import { initContract } from "@ts-rest/core";
import { healthContract } from "./health.js";

const c = initContract();

export const apiContract = c.router({
  Health: healthContract,
});
