const jwt = require("jsonwebtoken");
const { ApiError } = require("../utils");
const { env } = require("../config");

const authenticateToken = (req, res, next) => {
  const accessToken = req.cookies.accessToken;
  const refreshToken = req.cookies.refreshToken;
  if (req.headers["internal-service"]) {
    return basicAuth(req, res, next);
  }
  if (!accessToken || !refreshToken) {
    throw new ApiError(401, "Unauthorized. Please provide valid tokens.");
  }

  jwt.verify(accessToken, env.JWT_ACCESS_TOKEN_SECRET, (err, user) => {
    if (err) {
      throw new ApiError(
        401,
        "Unauthorized. Please provide valid access token."
      );
    }

    jwt.verify(
      refreshToken,
      env.JWT_REFRESH_TOKEN_SECRET,
      (err, refreshToken) => {
        if (err) {
          throw new ApiError(
            401,
            "Unauthorized. Please provide valid refresh token."
          );
        }

        req.user = user;
        req.refreshToken = refreshToken;
        next();
      }
    );
  });
};


const basicAuth = (req, res, next) => {
  const authHeader = req.headers["x-auth-token"];
  if (!authHeader) {
    throw new ApiError(401, "Unauthorized. Please provide valid tokens.");
  }
  if (authHeader !== env.NODEJS_AUTH_TOKEN) {
    throw new ApiError(401, "Unauthorized. auth token is not valid.");
  }
  next();
};


module.exports = { authenticateToken, basicAuth };
