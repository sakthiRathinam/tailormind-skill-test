const asyncHandler = require("express-async-handler");
const { getAllStudents, addNewStudent, getStudentDetail, setStudentStatus, updateStudent } = require("./students-service");
const handleGetAllStudents = asyncHandler(async (req, res) => {
    console.log("Getting all students");
    const dummyStudents = [
        {
            id: 1,
            name: "John Doe",
            email: "john.doe@example.com",
            grade: "10th",
            section: "A",
            status: "active"
        },
        {
            id: 2,
            name: "Jane Smith",
            email: "jane.smith@example.com", 
            grade: "11th",
            section: "B",
            status: "active"
        },
        {
            id: 3,
            name: "Mike Johnson",
            email: "mike.johnson@example.com",
            grade: "9th", 
            section: "C",
            status: "inactive"
        }
    ];
    
    res.json({ students: dummyStudents });
});

const handleAddStudent = asyncHandler(async (req, res) => {
    //write your code
        
});

const handleUpdateStudent = asyncHandler(async (req, res) => {
    //write your code

});

const handleGetStudentDetail = asyncHandler(async (req, res) => {
    //write your code

});

const handleStudentStatus = asyncHandler(async (req, res) => {
    //write your code

});

module.exports = {
    handleGetAllStudents,
    handleGetStudentDetail,
    handleAddStudent,
    handleStudentStatus,
    handleUpdateStudent,
};
