// Messages
export const NOT_FOUND_SUFFIX = "does not exist";
export const INTERNAL_SERVER_ERROR = "Internal Server Error";

// HTTP status codes
export enum HTTP_STATUS {
    OK = 200,
    CREATED = 201,
    NO_CONTENT = 204,
    NOT_MODIFIED = 304,
    BAD_REQUEST = 400,
    UNAUTHORIZED = 401,
    NOT_FOUND = 404,
    INTERNAL_SERVER_ERROR = 500,
}
