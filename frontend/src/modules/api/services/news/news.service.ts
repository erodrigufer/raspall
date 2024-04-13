import { Article } from "../../../articles/types";
import { generateGetMethod } from "../utils";

export const newsURL = "/v1/news";

// export const generateCarParkService = (): CarParkService => {
//   return {
//     find: generateFindMethod(carParkUrl),
//     get: generateGetMethod(carParkUrl),
//   };
// };
//

export const newsFunction = generateGetMethod<Article[]>(newsURL);
