import http from "http";
import url from "url";

const wrapResponse = (
    res: http.ServerResponse,
    code: number,
    body?: object
) => {
    const headers = {
        "Content-Type": "application/json",
    };
    res.writeHead(code, headers);
    res.end(JSON.stringify(body));
};

const getRequestBody = (req: http.IncomingMessage) => {
    let body = "";
    // Process the chunks of data and append it to the body
    req.on("data", (chunk) => {
        body += chunk.toString();
    });
    req.on("end", () => {});
    return body;
};

const getRequestQueryParams = (req: http.IncomingMessage) => {
    if (req.url === undefined) {
        throw new Error(
            "Undefined request URL, unable to get query parameters"
        );
    }
    return url.parse(req.url, true).query;
};

const getRequestPathname = (req: http.IncomingMessage) => {
    if (req.url === undefined) {
        throw new Error("Undefined request URL, unable to get path parameters");
    }
    return url.parse(req.url, true).pathname;
};

const wrapText = (text: string) => {
    return {
        message: `${text}`,
    };
};

export {
    getRequestBody,
    getRequestPathname,
    getRequestQueryParams,
    wrapResponse,
    wrapText,
};
