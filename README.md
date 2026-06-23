# BIT303 Shop Management System

[Chinese version](README-zh.md)

## Admin Frontend Usage

### 1. Start the Services

Start the backend service first. The default API address is:

    http://127.0.0.1:8000

Then start the admin frontend. The default frontend address is:

    http://localhost:8080

### 2. Log In to the Admin Panel

Open the admin login page in your browser:

    http://localhost:8080/#/login

Test account:

    Username: root
    Password: 123456

After a successful login, you will enter the admin dashboard.

### 3. Available Features

- Dashboard: view the current employee, personal credits, and management module overview.
- Employee List: create employees, edit employee information, enable or disable employees, and reset employee passwords.
- My Credits: view the current employee credit balance and credit records.
- Credit Operations: select an employee, add or deduct credits, and view the employee credit records.
- Category List: view current goods category information.
- User menu in the top-right corner: change the current account password or log out.

### 4. Notes

- The admin frontend sends backend API requests through the /backend route.
- If login fails, first confirm that the target environment has the root employee account and that its password is 123456.
