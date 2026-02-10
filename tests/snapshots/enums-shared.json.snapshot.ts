/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export const enum OrderStateEnum {
  OPEN = "open",
  CLOSED = "closed",
}

export const enum PriorityEnum {
  LOW = "low",
  MEDIUM = "medium",
  HIGH = "high",
}

export const enum StatusEnum {
  ACTIVE = "active",
  INACTIVE = "inactive",
  PENDING = "pending",
}

export type Components = {
  schemas: {
    Order: {
      id: string;
      priority: Components["schemas"]["Priority"];
      state: Components["schemas"]["OrderState"];
      status: Components["schemas"]["Status"];
    };
    OrderState: OrderStateEnum;
    Priority: PriorityEnum;
    Status: StatusEnum;
    User: {
      id: string;
      priority?: Components["schemas"]["Priority"];
      status: Components["schemas"]["Status"];
    };
  };
};

export type Routes = {
  "/orders": {
    post: {
      requestBody: {
        id: string;
        priority: PriorityEnum;
        state: OrderStateEnum;
        status: StatusEnum;
      };
      responses: {
        201: {
          id: string;
          priority: PriorityEnum;
          state: OrderStateEnum;
          status: StatusEnum;
        };
      };
    };
  };
  "/users": {
    get: {
      query: {
        status?: StatusEnum;
      };
      responses: {
        200: {
          id: string;
          priority?: PriorityEnum;
          status: StatusEnum;
        }[];
      };
    };
  };
};
