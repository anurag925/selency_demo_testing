const jwt = require("jsonwebtoken");
const { ApiError } = require("../utils");
const { env } = require("../config");

const authenticateServiceToken = (req, res, next) => {
  console.log("Authenticating token...");
  const serviceToken = req.headers["x-service-token"];

  console.log("Service Token:", serviceToken);

  if (!serviceToken) {
    throw new ApiError(401, "Unauthorized. Please provide a valid service token.");
  }

  jwt.verify(serviceToken, env.SERVICE_TOKEN_SECRET, (err, service) => {
    if (err) {
      throw new ApiError(
        401,
        "Unauthorized. Please provide a valid service token."
      );
    }

    req.service = service;
    next();
  });
}

module.exports = { authenticateServiceToken };