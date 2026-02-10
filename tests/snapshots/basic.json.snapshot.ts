/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export const enum StatusEnum {
  ACTIVE = "active",
  INACTIVE = "inactive",
}

export type Components = {
  schemas: {
    Status: StatusEnum;
    User: {
      id: string;
      status?: Components["schemas"]["Status"];
    };
  };
};

export type Routes = {
  "/ping": {
    get: {
      query: {
        limit?: number;
      };
      headers: {
        trace_id?: string;
      };
      responses: {
        200: string;
      };
    };
  };
};
