const asyncHandler = require("express-async-handler");
const { getAllStudents, addNewStudent, getStudentDetail, setStudentStatus, updateStudent } = require("./students-service");

const handleGetAllStudents = asyncHandler(async (req, res) => {
    const studentData = req.body;
    const students = await getAllStudents(studentData);
    res.json({
        success: true,
        students: students
    });
});


const handleAddStudent = asyncHandler(async (req, res) => {
    const studentData = req.body;
    const result = await addNewStudent(studentData);
    
    res.status(201).json({
        success: true,
        message: "Student added successfully",
        students: result
    });
});

const handleUpdateStudent = asyncHandler(async (req, res) => {
    const studentId = req.params.id;
    const studentData = req.body;
    const result = await updateStudent(studentId, studentData);
    
    res.json({
        success: true,
        message: "Student updated successfully",
        students: result
    });
});

const handleGetStudentDetail = asyncHandler(async (req, res) => {
    console.log("Fetching student details for ID:", req.params.id);
    const studentId = req.params.id;
    const student = await getStudentDetail(studentId);
    
    if (!student) {
        return res.status(404).json({
            success: false,
            message: "Student not found"
        });
    }
    
    res.json(student);
});

const handleStudentStatus = asyncHandler(async (req, res) => {
    const studentId = req.params.id;
    const { status } = req.body;
    const result = await setStudentStatus(studentId, status);
    
    res.json({
        success: true,
        message: "Student status updated successfully",
        students: result
    });
});

module.exports = {
    handleGetAllStudents,
    handleGetStudentDetail,
    handleAddStudent,
    handleStudentStatus,
    handleUpdateStudent,
};
