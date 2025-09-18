const {
  generateToken,
  generateCsrfHmacHash,
} = require("../utils");
const { v4: uuidV4 } = require("uuid");
const dotenv = require("dotenv");
dotenv.config();

const csrfToken = uuidV4();
const csrfHmacHash = generateCsrfHmacHash(csrfToken);
const accessToken = generateToken(
      { id: "golang-service", csrf_hmac: csrfHmacHash },
      env.JWT_ACCESS_TOKEN_SECRET,
      env.JWT_ACCESS_TOKEN_TIME_IN_MS
    );
console.log("Access Token:", accessToken);
console.log("CSRF Token:", csrfToken);
console.log("Use the above tokens in the Go service for authentication.");