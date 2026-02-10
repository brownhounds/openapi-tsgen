/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export type Components = {
  schemas: {
    FreeForm: Record<string, unknown>;
    MapPayload: {
      freeForm: Components["schemas"]["FreeForm"];
      mixedMap?: Components["schemas"]["MixedMap"];
      patternAndAdditional?: Components["schemas"]["PatternAndAdditional"];
      patterned: Components["schemas"]["Patterned"];
      stringMap?: Components["schemas"]["StringMap"];
    };
    MapResponse: {
      freeForm?: Components["schemas"]["FreeForm"];
      patterned?: Components["schemas"]["Patterned"];
    };
    MixedMap: ({
      id: string;
    } & Record<string, (string | number)>);
    PatternAndAdditional: ({ [K in `s-${string}`]?: string } & Record<string, (boolean | string)>);
    Patterned: { [K in (`${number}` | `x-${string}`)]?: (number | string) };
    StringMap: Record<string, string>;
  };
};

export type Routes = {
  "/maps": {
    post: {
      requestBody: {
        freeForm: Record<string, unknown>;
        mixedMap?: ({
        id: string;
      } & Record<string, (string | number)>);
        patternAndAdditional?: ({ [K in `s-${string}`]?: string } & Record<string, (boolean | string)>);
        patterned: { [K in (`${number}` | `x-${string}`)]?: (number | string) };
        stringMap?: Record<string, string>;
      };
      responses: {
        200: {
          freeForm?: Record<string, unknown>;
          patterned?: { [K in (`${number}` | `x-${string}`)]?: (number | string) };
        };
      };
    };
  };
};
