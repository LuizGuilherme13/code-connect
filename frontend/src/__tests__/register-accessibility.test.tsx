import { render } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import axe from 'axe-core';
import RegisterPage from '../components/pages/RegisterPage/RegisterPage';

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

describe('RegisterPage - Accessibility (WCAG 2)', () => {
  it('should not have any accessibility violations', async () => {
    const { container } = render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>
    );
    await expect(container).toHaveNoViolations();
  });

  it('should have proper heading hierarchy', () => {
    const { getByRole } = render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>
    );
    const heading = getByRole('heading', { level: 1 });
    expect(heading).toBeInTheDocument();
    expect(heading).toHaveTextContent('Cadastro');
  });

  it('should have labels for all form inputs', () => {
    const { getByLabelText } = render(
      <MemoryRouter>
        <RegisterPage />
      </MemoryRouter>
    );
    expect(getByLabelText('Nome')).toBeInTheDocument();
    expect(getByLabelText('Email')).toBeInTheDocument();
    expect(getByLabelText('Senha')).toBeInTheDocument();
    expect(getByLabelText('Confirmar Senha')).toBeInTheDocument();
  });

  it('should have accessible social buttons', () => {
    const { getByRole } = render(
      <MemoryRouter>
        <RegisterPage />
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
        <RegisterPage />
      </MemoryRouter>
    );
    expect(getByRole('link', { name: /faça seu login/i })).toBeInTheDocument();
  });
});
