# Contributing Guide

Welcome to velox! We're thrilled you're here. Whether you're fixing a typo, squashing a bug, or building a massive new feature, your contributions are what make the open-source community amazing.

These guidelines will help you get set up and smoothly navigate our contribution process.

---

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Ways to Contribute](#ways-to-contribute)
3. [Getting Started](#getting-started)
4. [Development Setup](#development-setup)
5. [Project Structure](#project-structure)
6. [Branching & Workflow](#branching--workflow)
7. [Coding Standards](#coding-standards)
8. [Commit Message Guidelines](#commit-message-guidelines)
9. [Testing](#testing)
10. [Submitting a Pull Request](#submitting-a-pull-request)
11. [Reporting Bugs](#reporting-bugs)
12. [Suggesting Features](#suggesting-features)
13. [Documentation Contributions](#documentation-contributions)
14. [Need Help?](#need-help)

---

## Code of Conduct

We are committed to maintaining a welcoming and inclusive environment. By participating, you agree to adhere to our Code of Conduct. Please read the `CODE_OF_CONDUCT.md` file (if available) to understand the expectations we have for all community members.

---

## Ways to Contribute

There are many ways to make an impact:
- **Write Code:** Tackle open bugs or implement requested features.
- **Improve Docs:** Fix typos, add examples, or rewrite confusing sections.
- **Report Issues:** Let us know when something is broken.
- **Suggest Features:** Share an idea for how we can improve.
- **Review PRs:** Help review and approve community contributions.

---

## Getting Started

Ensure you have the following installed before picking up a task:
- [Node.js](https://nodejs.org/) (v18+)
- [npm](https://www.npmjs.com/) or [Yarn](https://yarnpkg.com/)
- [Git](https://git-scm.com/)

> **Note:** Replace these with your actual stack if different.

---

## Development Setup

Ready to write some code? Follow these commands to set up your local environment:

1. **Fork** the repository and **clone** your fork:
   ```bash
   git clone https://github.com/[YOUR-USERNAME]/[PROJECT-NAME].git
   cd [PROJECT-NAME]
   ```

2. **Install dependencies**:
   ```bash
   npm install
   ```

3. **Start the development server**:
   ```bash
   npm run dev
   ```

4. **Build the project** (to verify production builds):
   ```bash
   npm run build
   ```

---

## Project Structure

Here is a quick overview of where things live:
```text
├── src/           # Application source code
│   ├── components/# Reusable UI elements
│   ├── pages/     # Routing and views
│   └── utils/     # Helpers and utility functions
├── public/        # Static assets (images, icons)
├── tests/         # Automated test suites
└── README.md      # Main project documentation
```

---

## Branching & Workflow

Stay organized by using feature branches:

1. Ensure your local `main` branch is up to date with the main repository.
2. Create a specific branch for your work:
   ```bash
   git checkout -b feature/cool-new-button
   git checkout -b bugfix/login-crash
   ```
3. Keep branches focused. Create separate branches for separate tasks.

---

## Coding Standards

Maintain consistent code quality with our linting tools. Before pushing, run:
```bash
npm run lint
npm run lint:fix
```

- **Naming Conventions:** Use `camelCase` for variables/functions, and `PascalCase` for React components/classes.
- **Clarity over Cleverness:** Ensure complex logic is well-commented.

---

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/). This keeps our git history clean and readable.

**Format:**
```text
type(scope): brief description
```

**Examples:**
- `feat(auth): add google oauth login`
- `fix(ui): resolve button alignment issue on mobile`
- `docs(readme): update installation instructions`

**Types include:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.

---

## Testing

We require tests to ensure everything remains stable.

- **Run the full test suite:**
  ```bash
  npm test
  ```
- **Run tests in watch mode** (while developing):
  ```bash
  npm run test:watch
  ```
- **Check test coverage:**
  ```bash
  npm run test:coverage
  ```

---

## Submitting a Pull Request

When your work is ready for review:

1. Push your branch up to your fork:
   ```bash
   git push origin feature/cool-new-button
   ```
2. Open a Pull Request against our `main` branch.
3. Fill out the PR template completely.

### Pull Request Checklist
- [ ] I've read the Contribution Guidelines.
- [ ] Code is styled and linted correctly.
- [ ] Tests have been added/updated and they pass.
- [ ] Documentation (like READMEs) updated if applicable.
- [ ] Commit messages follow the Conventional Commits guidelines.

---

## Reporting Bugs

Spotted a bug? Follow these steps:

1. **Search** the existing [Issues](https://github.com/[YOUR-USERNAME]/[PROJECT-NAME]/issues) to see if it's already reported.
2. If not, open a new Issue using our Bug Report format.
3. **Include details:** clear steps to reproduce, expected vs. actual behavior, and error logs/screenshots.

---

## Suggesting Features

We welcome your ideas!

1. Check existing issues to avoid duplicates.
2. Open a new issue outlining the feature.
3. Clearly describe the problem your feature solves and how it might look. Mockups are highly appreciated!

---

## Documentation Contributions

Good documentation is critical. If you spot a typo, confusing wording, or missing steps:

- For massive rewrites, open an Issue to discuss it first.
- For quick fixes or typos, feel free to jump straight to a PR!

---

## Need Help?

Stuck somewhere? We're here to help!
- **Open an Issue:** Tag it as a `question`.
- **Reach out:** Email us at `[your-email@example.com]` or join our `[Discord / Workspace link]`.

Thank you for contributing!