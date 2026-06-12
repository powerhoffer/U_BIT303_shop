# BIT303 Shop Admin Guide

## Admin Frontend Usage

### 1. Start Services

Start the backend service first. The default API address is:

```bash
http://127.0.0.1:8000
```

Then start the admin frontend. The default frontend address is:

```bash
http://localhost:8080
```

### 2. Log In

Open the admin login page in a browser:

```bash
http://localhost:8080/#/login
```

Demo account:

```text
Username: root
Password: 123456
```

After a successful login, the system opens the admin dashboard.

### 3. Main Features

- Dashboard: View the current employee, credit balance, and management module overview.
- Employees: Create employees, edit employee information, enable or disable employees, and reset passwords.
- My Credits: View the current employee's credit balance and credit records.
- Credit Operations: Select an employee, add or deduct credits, and review employee credit records.
- Categories: View product category information.
- User menu in the top-right corner: Change the current password or log out.

### 4. Notes

- The admin frontend calls backend APIs under the `/backend` route group.
- If login fails, make sure the customer environment has an active `root` employee account with password `123456`.
