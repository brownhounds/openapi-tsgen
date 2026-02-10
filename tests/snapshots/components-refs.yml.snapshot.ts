/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:29Z
 */

export type Components = {
  schemas: {
    User: {
      email: string;
      id: string;
    };
  };
  responses: {
    NotFound: string;
    UserList: {
      headers: {
        "X-Rate-Limit"?: Components["headers"]["RateLimit"];
      };
      body: {
      email: string;
      id: string;
    }[];
    };
  };
  requestBodies: {
    UserUpdate: {
      email: string;
      id: string;
    };
  };
  parameters: {
    Limit: number;
    TraceId: string;
  };
  headers: {
    RateLimit: number;
  };
};

export type Routes = {
  "/users": {
    get: {
      query: {
        limit?: Components["parameters"]["Limit"];
      };
      headers: {
        "x-trace-id"?: Components["parameters"]["TraceId"];
      };
      responses: {
        200: Components["responses"]["UserList"];
      };
    };
  };
  "/users/{id}": {
    post: {
      params: {
        id: string;
      };
      requestBody: Components["requestBodies"]["UserUpdate"];
      responses: {
        204: never;
        404: Components["responses"]["NotFound"];
      };
    };
  };
};
