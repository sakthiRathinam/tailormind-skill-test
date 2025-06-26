-- Create students table with essential fields
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20),
    gender VARCHAR(10),
    dob DATE,
    class VARCHAR(50),
    section VARCHAR(10),
    roll INTEGER,
    father_name VARCHAR(100),
    father_phone VARCHAR(20),
    mother_name VARCHAR(100),
    mother_phone VARCHAR(20),
    guardian_name VARCHAR(100),
    guardian_phone VARCHAR(20),
    relation_of_guardian VARCHAR(50),
    current_address TEXT,
    permanent_address TEXT,
    admission_date DATE,
    system_access BOOLEAN DEFAULT true,
    reporter_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert 10 dummy student records
INSERT INTO students (name, email, phone, gender, dob, class, section, roll, father_name, father_phone, mother_name, mother_phone, guardian_name, guardian_phone, relation_of_guardian, current_address, permanent_address, admission_date, system_access, reporter_name) VALUES

('John Doe', 'john.doe@school.com', '+1234567890', 'Male', '2005-01-15', '10th Grade', 'A', 101, 'Robert Doe', '+1234567891', 'Mary Doe', '+1234567892', 'Robert Doe', '+1234567891', 'Father', '123 Main St, City, State 12345', '123 Main St, City, State 12345', '2023-01-15', true, 'Admin User'),

('Jane Smith', 'jane.smith@school.com', '+1234567893', 'Female', '2006-03-20', '9th Grade', 'B', 102, 'Michael Smith', '+1234567894', 'Sarah Smith', '+1234567895', 'Michael Smith', '+1234567894', 'Father', '456 Oak Ave, Town, State 67890', '456 Oak Ave, Town, State 67890', '2023-02-10', true, 'Teacher Johnson'),

('Alex Wilson', 'alex.wilson@school.com', '+1234567896', 'Male', '2004-07-08', '11th Grade', 'A', 103, 'David Wilson', '+1234567897', 'Lisa Wilson', '+1234567898', 'David Wilson', '+1234567897', 'Father', '789 Pine St, Village, State 11111', '789 Pine St, Village, State 11111', '2022-08-20', true, 'Principal Davis'),

('Emily Johnson', 'emily.johnson@school.com', '+1234567899', 'Female', '2005-11-25', '10th Grade', 'C', 104, 'James Johnson', '+1234567800', 'Patricia Johnson', '+1234567801', 'James Johnson', '+1234567800', 'Father', '321 Elm Street, Downtown, State 22222', '321 Elm Street, Downtown, State 22222', '2023-01-05', false, 'Teacher Brown'),

('Michael Chen', 'michael.chen@school.com', '+1234567802', 'Male', '2006-09-12', '9th Grade', 'A', 105, 'Wei Chen', '+1234567803', 'Li Chen', '+1234567804', 'Wei Chen', '+1234567803', 'Father', '654 Maple Drive, Suburb, State 33333', '654 Maple Drive, Suburb, State 33333', '2023-03-15', true, 'Vice Principal Lee'),

('Sarah Davis', 'sarah.davis@school.com', '+1234567805', 'Female', '2004-12-03', '11th Grade', 'B', 106, 'Mark Davis', '+1234567806', 'Jennifer Davis', '+1234567807', 'Mark Davis', '+1234567806', 'Father', '987 Cedar Lane, Uptown, State 44444', '987 Cedar Lane, Uptown, State 44444', '2022-07-10', true, 'Teacher Garcia'),

('David Rodriguez', 'david.rodriguez@school.com', '+1234567808', 'Male', '2005-04-18', '10th Grade', 'B', 107, 'Carlos Rodriguez', '+1234567809', 'Maria Rodriguez', '+1234567810', 'Carlos Rodriguez', '+1234567809', 'Father', '147 Oak Ridge, Eastside, State 55555', '147 Oak Ridge, Eastside, State 55555', '2023-01-20', true, 'Counselor Martinez'),

('Ashley Taylor', 'ashley.taylor@school.com', '+1234567811', 'Female', '2006-06-30', '9th Grade', 'C', 108, 'Robert Taylor', '+1234567812', 'Nancy Taylor', '+1234567813', 'Helen Taylor', '+1234567814', 'Grandmother', '258 Birch Avenue, Westside, State 66666', '258 Birch Avenue, Westside, State 66666', '2023-02-28', true, 'Teacher Thompson'),

('Christopher Lee', 'christopher.lee@school.com', '+1234567815', 'Male', '2004-10-14', '11th Grade', 'C', 109, 'Andrew Lee', '+1234567816', 'Susan Lee', '+1234567817', 'Andrew Lee', '+1234567816', 'Father', '369 Willow Court, Northside, State 77777', '741 Spruce Way, Oldtown, State 88888', '2022-09-05', false, 'Dean Wilson'),

('Jessica Brown', 'jessica.brown@school.com', '+1234567818', 'Female', '2005-08-22', '10th Grade', 'A', 110, 'Kevin Brown', '+1234567819', 'Michelle Brown', '+1234567820', 'Kevin Brown', '+1234567819', 'Father', '852 Poplar Street, Southside, State 99999', '852 Poplar Street, Southside, State 99999', '2023-01-30', true, 'Teacher Anderson');

-- Create index on commonly queried fields
CREATE INDEX IF NOT EXISTS idx_students_email ON students(email);
CREATE INDEX IF NOT EXISTS idx_students_class_section ON students(class, section);
CREATE INDEX IF NOT EXISTS idx_students_roll ON students(roll);
CREATE INDEX IF NOT EXISTS idx_students_name ON students(name); 