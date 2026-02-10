/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export type Routes = {
  "/items/{id}": {
    get: {
      params: {
        id: string;
      };
      query: {
        q?: string;
      };
      headers: {
        "X-Trace-Id"?: string;
      };
      cookies: {
        session?: string;
      };
      responses: {
        200: {
          headers: {
            "X-Request-Id"?: string;
          };
          body: {
          id?: string;
          q?: string;
        };
        };
      };
    };
  };
};
