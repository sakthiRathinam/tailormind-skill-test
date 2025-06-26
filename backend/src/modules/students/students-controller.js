const asyncHandler = require("express-async-handler");
const { getStudentDetail } = require("./students-service");

const handleGetStudentDetail = asyncHandler(async (req, res) => {
    const { id } = req.params;
    
    if (!id || isNaN(parseInt(id))) {
        return res.status(400).json({
            success: false,
            message: "Invalid student ID provided"
        });
    }
    
    const student = await getStudentDetail(parseInt(id));
    
    if (!student) {
        return res.status(404).json({
            success: false,
            message: "Student not found"
        });
    }
    
    res.json({
        success: true,
        message: "Student details retrieved successfully",
        data: student
    });
});

module.exports = {
    handleGetStudentDetail
};
