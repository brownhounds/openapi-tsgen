/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export type Routes = {
  "/echo": {
    post: {
      requestBody: ({
        message?: string;
      } | string | {
        name?: string;
      } | Record<string, unknown> | {
        description?: string;
        file: string;
      });
      responses: {
        200: (Record<string, unknown> | string);
      };
    };
  };
};
