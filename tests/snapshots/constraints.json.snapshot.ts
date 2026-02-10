/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export const enum IfThenElseSampleKindEnum {
  A = "a",
  B = "b",
}

export type Components = {
  schemas: {
    Account: {
      email: string;
      id?: string;
      password: string;
    };
    DependentRequiredSample: ({
      billingAddress?: string;
      creditCard?: string;
      name?: string;
    } & ({
      creditCard?: never;
    } | {
      billingAddress: string;
      creditCard: string;
    }));
    IfThenElseSample: ((({
      aOnly?: string;
      bOnly?: number;
      kind: IfThenElseSampleKindEnum;
    } & {
      kind: "a";
    }) & {
      aOnly: unknown;
    }) | (({
      aOnly?: string;
      bOnly?: number;
      kind: IfThenElseSampleKindEnum;
    } & {
      kind: "b";
    }) & {
      bOnly: unknown;
    }));
    NotString: string;
  };
};

export type Routes = {
  "/accounts": {
    post: {
      requestBody: {
        email: string;
        password: string;
      };
      responses: {
        201: {
          email: string;
          id?: string;
        };
      };
    };
  };
};
