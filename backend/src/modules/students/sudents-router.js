const express = require("express");
const router = express.Router();
const studentController = require("./students-controller");

// GET /api/v1/students/:id - Get single student details
router.get("/:id", studentController.handleGetStudentDetail);

module.exports = { studentsRoutes: router };
