const { ApiError } = require("../../utils");
const { findStudentById } = require("./students-repository");

const getStudentDetail = async (id) => {
    try {
        const student = await findStudentById(id);
        return student;
    } catch (error) {
        console.error("Error fetching student detail:", error);
        throw new ApiError(500, "Failed to retrieve student details");
    }
};

module.exports = {
    getStudentDetail
};
