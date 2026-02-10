/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export type Components = {
  securitySchemes: {
    ApiKeyAuth: {
      in: "header";
      name: "X-API-Key";
      type: "apiKey";
    };
    BearerAuth: {
      bearerFormat: "JWT";
      scheme: "bearer";
      type: "http";
    };
    OAuth2Auth: {
      flows: {
        authorizationCode: {
          authorizationUrl: string;
          tokenUrl: string;
          scopes: {
            "read:users": string;
            "write:users": string;
          };
        };
      };
      type: "oauth2";
    };
    OpenIdConnect: {
      openIdConnectUrl: "https://auth.example.com/.well-known/openid-configuration";
      type: "openIdConnect";
    };
  };
};

export type Routes = {
  "/secure": {
    get: {
      security: ({
        ApiKeyAuth: string[];
      } | {
        BearerAuth: string[];
      })[];
      responses: {
        200: string;
      };
    };
  };
};
