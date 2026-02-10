/*
 * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT
 *
 * Generator: openapi-tsgen@dev
 * OpenAPI version: 3.1.1
 * Generated at: 2026-02-10T17:22:30Z
 */

export const enum BikeKindBikeEnum {
  BIKE = "bike",
}

export const enum CarDoorsEnum {
  VALUE_2 = 2,
  VALUE_4 = 4,
}

export const enum CarKindCarEnum {
  CAR = "car",
}

export const enum CatHuntingSkillEnum {
  CLUELESS = "clueless",
  LAZY = "lazy",
}

export const enum CatPetTypeCatEnum {
  CAT = "cat",
}

export const enum DogPetTypeDogEnum {
  DOG = "dog",
}

export const enum MixedAnyAllOneAEnum {
  A = "a",
}

export const enum MixedAnyAllOneBEnum {
  B = "b",
}

export type Components = {
  schemas: {
    Bike: {
      hasBell: boolean;
      kind: BikeKindBikeEnum;
    };
    Car: {
      doors: CarDoorsEnum;
      kind: CarKindCarEnum;
    };
    Cat: (Components["schemas"]["PetBase"] & {
      huntingSkill: CatHuntingSkillEnum;
      petType?: CatPetTypeCatEnum;
    });
    Dog: (Components["schemas"]["PetBase"] & {
      packSize: number;
      petType?: DogPetTypeDogEnum;
    });
    MaybeString: (string | null);
    MixedAnyAllOne: ((string | number) & (MixedAnyAllOneAEnum | MixedAnyAllOneBEnum));
    OneOfWithNull: (string | null);
    Pet: (Components["schemas"]["Cat"] | Components["schemas"]["Dog"]);
    PetBase: {
      name: string;
      petType: string;
    };
    PolyRequest: {
      maybe?: Components["schemas"]["MaybeString"];
      mixed?: Components["schemas"]["MixedAnyAllOne"];
      oneOrNull?: Components["schemas"]["OneOfWithNull"];
      pet: Components["schemas"]["Pet"];
      stringOrNumber?: Components["schemas"]["StringOrNumber"];
      vehicle: Components["schemas"]["Vehicle"];
    };
    PolyResponse: {
      pet: Components["schemas"]["Pet"];
    };
    StringOrNumber: (string | number);
    Vehicle: (Components["schemas"]["Car"] | Components["schemas"]["Bike"]);
  };
};

export type Routes = {
  "/polymorph": {
    post: {
      requestBody: {
        maybe?: (string | null);
        mixed?: ((string | number) & (MixedAnyAllOneAEnum | MixedAnyAllOneBEnum));
        oneOrNull?: (string | null);
        pet: (({
        name: string;
        petType: string;
      } & {
        huntingSkill: CatHuntingSkillEnum;
        petType?: CatPetTypeCatEnum;
      }) | ({
        name: string;
        petType: string;
      } & {
        packSize: number;
        petType?: DogPetTypeDogEnum;
      }));
        stringOrNumber?: (string | number);
        vehicle: ({
        doors: CarDoorsEnum;
        kind: CarKindCarEnum;
      } | {
        hasBell: boolean;
        kind: BikeKindBikeEnum;
      });
      };
      responses: {
        200: {
          pet: (({
          name: string;
          petType: string;
        } & {
          huntingSkill: CatHuntingSkillEnum;
          petType?: CatPetTypeCatEnum;
        }) | ({
          name: string;
          petType: string;
        } & {
          packSize: number;
          petType?: DogPetTypeDogEnum;
        }));
        };
      };
    };
  };
};
