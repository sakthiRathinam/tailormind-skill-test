const { processDBRequest } = require("../../utils");

const findStudentById = async (id) => {
    const query = `
        SELECT
            id,
            name,
            email,
            phone,
            gender,
            dob,
            class,
            section,
            roll,
            father_name AS "fatherName",
            father_phone AS "fatherPhone",
            mother_name AS "motherName",
            mother_phone AS "motherPhone",
            guardian_name AS "guardianName",
            guardian_phone AS "guardianPhone",
            relation_of_guardian AS "relationOfGuardian",
            current_address AS "currentAddress",
            permanent_address AS "permanentAddress",
            admission_date AS "admissionDate",
            system_access AS "systemAccess",
            reporter_name AS "reporterName"
        FROM students
        WHERE id = $1`;
    
    const queryParams = [id];
    const { rows } = await processDBRequest({ query, queryParams });
    return rows[0];
};

module.exports = {
    findStudentById
};
