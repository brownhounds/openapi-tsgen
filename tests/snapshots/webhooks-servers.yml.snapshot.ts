/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export type Servers = {
  description: "Primary API";
  url: "https://api.example.com";
}[];

export type Routes = {
  "/events": {
    post: {
      servers: {
        description: "Primary API";
        url: "https://api.example.com";
      }[];
      requestBody: {
        name?: string;
      };
      responses: {
        201: never;
      };
    };
  };
};

export type Webhooks = {
  "user.created": {
    post: {
      servers: {
        description: "Primary API";
        url: "https://api.example.com";
      }[];
      requestBody: {
        id: string;
      };
      responses: {
        200: never;
      };
    };
  };
};
