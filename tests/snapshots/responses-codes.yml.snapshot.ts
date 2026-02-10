/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export type Routes = {
  "/status": {
    get: {
      responses: {
        200: {
          ok?: boolean;
        };
        "4XX": string;
        default: {
          message: string;
        };
      };
    };
  };
};
