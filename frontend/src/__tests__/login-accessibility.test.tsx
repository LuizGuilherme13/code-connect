import { render } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import axe from 'axe-core';
import LoginPage from '../components/pages/LoginPage/LoginPage';

expect.extend({
  async toHaveNoViolations(container: Element) {
    const results = await axe.run(container);
    const violations = results.violations;

    if (violations.length === 0) {
      return {
        pass: true,
        message: () => 'Expected accessibility violations but found none',
      };
    }

    const violationMessages = violations.map((v) => {
      const nodes = v.nodes.map((n) => `  - ${n.html}\n    ${n.failureSummary}`).join('\n');
      return `[${v.id}] ${v.impact}: ${v.description}\n${nodes}`;
    });

    return {
      pass: false,
      message: () =>
        `Found ${violations.length} accessibility violation(s):\n\n${violationMessages.join('\n\n')}`,
    };
  },
});

describe('LoginPage - Accessibility (WCAG 2)', () => {
  it('should not have any accessibility violations', async () => {
    const { container } = render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );
    await expect(container).toHaveNoViolations();
  });

  it('should have proper heading hierarchy', () => {
    const { getByRole } = render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );
    const heading = getByRole('heading', { level: 1 });
    expect(heading).toBeInTheDocument();
    expect(heading).toHaveTextContent('Login');
  });

  it('should have labels for all form inputs', () => {
    const { getByLabelText } = render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );
    expect(getByLabelText('Email ou usuário')).toBeInTheDocument();
    expect(getByLabelText('Senha')).toBeInTheDocument();
    expect(getByLabelText('Lembrar-me')).toBeInTheDocument();
  });

  it('should have accessible social buttons', () => {
    const { getByRole } = render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );
    const githubBtn = getByRole('button', { name: /github/i });
    const gmailBtn = getByRole('button', { name: /gmail/i });
    expect(githubBtn).toBeInTheDocument();
    expect(gmailBtn).toBeInTheDocument();
  });

  it('should have accessible navigation links', () => {
    const { getByRole } = render(
      <MemoryRouter>
        <LoginPage />
      </MemoryRouter>
    );
    expect(getByRole('link', { name: /esqueci a senha/i })).toBeInTheDocument();
    expect(getByRole('link', { name: /crie seu cadastro/i })).toBeInTheDocument();
  });
});
